package s3api

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"path"
	"strings"

	"github.com/hanzoai/s3/s3/glog"
	"github.com/hanzoai/s3/s3/pb/filer_pb"
	"github.com/hanzoai/s3/s3/s3api/s3_constants"
	"github.com/hanzoai/s3/s3/s3api/s3lifecycle"
	stats_collect "github.com/hanzoai/s3/s3/stats"
	s3_lifecyclewire "github.com/hanzoai/s3/s3/wire/s3_lifecycle"
)

// Compile-time check: S3ApiServer is the live LifecycleDeleter backing the
// native ZAP HanzoS3LifecycleInternal service.
var _ s3_lifecyclewire.LifecycleDeleter = (*S3ApiServer)(nil)

// entryIdentity is the in-process value form of the wire EntryIdentity:
// computeEntryIdentity builds it from the live filer entry and
// identityMatches compares it against the request's CAS witness.
type entryIdentity struct {
	MtimeNs      int64
	Size         int64
	HeadFid      string
	ExtendedHash []byte
}

// LifecycleDelete executes one (rule, action) verdict: re-fetch, identity
// CAS, object-lock check, dispatch by kind. Errors surface as outcomes;
// reader cursors and pending state are the worker's concern. req is the
// zero-copy ZAP view over the wire LifecycleDeleteRequest; the returned
// result carries the outcome enum and reason for the wire response.
func (s3a *S3ApiServer) LifecycleDelete(req s3_lifecyclewire.LifecycleDeleteRequest) (s3_lifecyclewire.LifecycleDeleteResult, error) {
	bucket, objectPath, versionID := req.Bucket(), req.ObjectPath(), req.VersionID()
	if bucket == "" || objectPath == "" {
		return blocked("FATAL_EVENT_ERROR: empty bucket or object_path"), nil
	}

	// MPU init lives at .uploads/<id>/; not handled by getObjectEntry.
	if req.ActionKind() == s3_lifecyclewire.ActionKindAbortMpu {
		return s3a.lifecycleAbortMPU(req)
	}

	entry, err := s3a.getObjectEntry(bucket, objectPath, versionID)
	if err != nil {
		if errors.Is(err, filer_pb.ErrNotFound) || errors.Is(err, ErrObjectNotFound) || errors.Is(err, ErrVersionNotFound) || errors.Is(err, ErrLatestVersionNotFound) {
			return noopResolved("NOT_FOUND"), nil
		}
		glog.V(1).Infof("lifecycle: live fetch %s/%s@%s: %v", bucket, objectPath, versionID, err)
		return retryLater("TRANSPORT_ERROR: " + err.Error()), nil
	}

	if !identityMatches(computeEntryIdentity(entry), req) {
		return noopResolved("STALE_IDENTITY"), nil
	}

	// Lifecycle never bypasses governance/compliance; the http.Request is
	// only read when bypass is allowed, so nil is safe here.
	if err := s3a.enforceObjectLockProtections(nil, bucket, objectPath, versionID, false); err != nil {
		glog.V(2).Infof("lifecycle: SKIPPED_OBJECT_LOCK %s/%s@%s: %v", bucket, objectPath, versionID, err)
		return s3_lifecyclewire.LifecycleDeleteResult{
			Outcome: s3_lifecyclewire.LifecycleDeleteOutcomeSkippedObjectLock,
			Reason:  err.Error(),
		}, nil
	}

	return s3a.lifecycleDispatch(context.Background(), req, entry)
}

func (s3a *S3ApiServer) lifecycleDispatch(ctx context.Context, req s3_lifecyclewire.LifecycleDeleteRequest, entry *filer_pb.Entry) (s3_lifecyclewire.LifecycleDeleteResult, error) {
	bucket, objectPath, versionID := req.Bucket(), req.ObjectPath(), req.VersionID()
	metadataOnly := entryUsesMetadataOnlyDelete(entry)
	switch req.ActionKind() {
	case s3_lifecyclewire.ActionKindExpirationDays, s3_lifecyclewire.ActionKindExpirationDate:
		// Current-version expiration: Enabled -> delete marker; Suspended
		// -> delete null + new marker; Off -> remove. Filer errors classify
		// as RETRY_LATER; the worker's budget promotes to BLOCKED.
		state, vErr := s3a.getVersioningState(bucket)
		if vErr != nil {
			if errors.Is(vErr, filer_pb.ErrNotFound) {
				return noopResolved("BUCKET_NOT_FOUND"), nil
			}
			return retryLater("TRANSPORT_ERROR: versioning lookup: " + vErr.Error()), nil
		}
		switch state {
		case s3_constants.VersioningEnabled:
			if _, err := s3a.createDeleteMarker(bucket, objectPath); err != nil {
				return retryLater("TRANSPORT_ERROR: createDeleteMarker: " + err.Error()), nil
			}
			return done(), nil
		case s3_constants.VersioningSuspended:
			// Best-effort null delete; NotFound is benign.
			if err := s3a.deleteSpecificObjectVersion(ctx, bucket, objectPath, "null", metadataOnly); err != nil {
				if !errors.Is(err, filer_pb.ErrNotFound) && !errors.Is(err, ErrVersionNotFound) {
					return retryLater("TRANSPORT_ERROR: deleteNullVersion: " + err.Error()), nil
				}
			}
			if _, err := s3a.createDeleteMarker(bucket, objectPath); err != nil {
				return retryLater("TRANSPORT_ERROR: createDeleteMarker: " + err.Error()), nil
			}
			return done(), nil
		default:
			err := s3a.WithFilerClient(false, func(c filer_pb.HanzoFilerClient) error {
				return s3a.deleteUnversionedObjectWithClient(c, bucket, objectPath, metadataOnly)
			})
			if err != nil {
				if errors.Is(err, filer_pb.ErrNotFound) || errors.Is(err, ErrObjectNotFound) {
					return noopResolved("NOT_FOUND_AT_DELETE"), nil
				}
				return retryLater("TRANSPORT_ERROR: deleteUnversioned: " + err.Error()), nil
			}
			recordMetadataOnlyIf(metadataOnly, req)
			return done(), nil
		}

	case s3_lifecyclewire.ActionKindNoncurrentDays,
		s3_lifecyclewire.ActionKindNewerNoncurrent,
		s3_lifecyclewire.ActionKindExpiredDeleteMarker:
		// EXPIRED_DELETE_MARKER targets the marker version itself.
		if versionID == "" {
			return blocked("FATAL_EVENT_ERROR: version_id required for noncurrent / delete-marker delete"), nil
		}
		// Latest-pointer guard for noncurrent kinds: refuse to delete
		// the version that the .versions/ directory currently points
		// to. The router can't always tell current from noncurrent
		// without sibling state, so the server checks here.
		if req.ActionKind() == s3_lifecyclewire.ActionKindNoncurrentDays ||
			req.ActionKind() == s3_lifecyclewire.ActionKindNewerNoncurrent {
			isLatest, lookupErr := s3a.isCurrentLatestVersion(bucket, objectPath, versionID)
			if lookupErr != nil {
				if errors.Is(lookupErr, filer_pb.ErrNotFound) || errors.Is(lookupErr, ErrObjectNotFound) {
					return noopResolved("NOT_FOUND"), nil
				}
				return retryLater("TRANSPORT_ERROR: latest-pointer lookup: " + lookupErr.Error()), nil
			}
			if isLatest {
				return noopResolved("VERSION_IS_LATEST"), nil
			}
		}
		// Re-check sole-survivor: a fresh PUT can land between schedule
		// and dispatch. Identity-CAS upstream covers the marker bytes;
		// this covers the directory shape.
		if req.ActionKind() == s3_lifecyclewire.ActionKindExpiredDeleteMarker {
			outcome, ok, err := s3a.checkSoleSurvivorMarker(ctx, bucket, objectPath, versionID)
			if ok || err != nil {
				return outcome, err
			}
		}
		if err := s3a.deleteSpecificObjectVersion(ctx, bucket, objectPath, versionID, metadataOnly); err != nil {
			if errors.Is(err, filer_pb.ErrNotFound) || errors.Is(err, ErrVersionNotFound) || errors.Is(err, ErrObjectNotFound) {
				return noopResolved("NOT_FOUND_AT_DELETE"), nil
			}
			return retryLater("TRANSPORT_ERROR: deleteSpecificVersion: " + err.Error()), nil
		}
		recordMetadataOnlyIf(metadataOnly, req)
		return done(), nil

	case s3_lifecyclewire.ActionKindAbortMpu:
		return blocked("FATAL_EVENT_ERROR: ABORT_MPU dispatched after fetch"), nil

	default:
		return blocked("FATAL_EVENT_ERROR: unknown action_kind " + actionKindLabel(req.ActionKind())), nil
	}
}

func (s3a *S3ApiServer) lifecycleAbortMPU(req s3_lifecyclewire.LifecycleDeleteRequest) (s3_lifecyclewire.LifecycleDeleteResult, error) {
	bucket, objectPath := req.Bucket(), req.ObjectPath()
	// objectPath is `.uploads/<upload_id>` (set by the router from the
	// init directory's bucket-relative path); reject anything that isn't
	// exactly that shape so a malformed event can't escalate to a wider rm.
	const uploadsPrefix = s3_constants.MultipartUploadsFolder + "/"
	if !strings.HasPrefix(objectPath, uploadsPrefix) {
		return blocked("FATAL_EVENT_ERROR: ABORT_MPU object_path missing .uploads/ prefix"), nil
	}
	uploadID := objectPath[len(uploadsPrefix):]
	// Reject "." and ".." explicitly: util.JoinPath in the filer cleans
	// path components, so .uploads/.. would resolve to the bucket root.
	if uploadID == "" || uploadID == "." || uploadID == ".." || strings.ContainsRune(uploadID, '/') {
		return blocked("FATAL_EVENT_ERROR: ABORT_MPU object_path malformed: " + objectPath), nil
	}

	uploadsFolder := s3a.genUploadsFolder(bucket)
	// Pre-check existence: filer.DeleteEntry suppresses ErrNotFound and
	// returns success, so without this check an already-aborted upload
	// would report DONE instead of the correct NOOP_RESOLVED.
	exists, err := s3a.exists(uploadsFolder, uploadID, true)
	if err != nil {
		if errors.Is(err, filer_pb.ErrNotFound) {
			return noopResolved("NOT_FOUND"), nil
		}
		return retryLater("TRANSPORT_ERROR: exists: " + err.Error()), nil
	}
	if !exists {
		return noopResolved("NOT_FOUND"), nil
	}
	if err := s3a.rm(uploadsFolder, uploadID, true, true); err != nil {
		if errors.Is(err, filer_pb.ErrNotFound) {
			return noopResolved("NOT_FOUND_AT_DELETE"), nil
		}
		glog.V(1).Infof("lifecycle abort_mpu %s/%s: %v", bucket, objectPath, err)
		return retryLater("TRANSPORT_ERROR: rm: " + err.Error()), nil
	}
	return done(), nil
}

// checkSoleSurvivorMarker returns (terminal-result, true, nil) when state
// has drifted: count != 1, the surviving entry is a different version, the
// .versions/ directory's latest pointer doesn't name versionId, or a bare
// null-version exists outside .versions/. It returns (_, false, nil) to
// proceed with the delete. Pointer missing while a marker is present is
// treated as retry-later — the create races with the directory metadata
// update.
func (s3a *S3ApiServer) checkSoleSurvivorMarker(ctx context.Context, bucket, object, versionId string) (s3_lifecyclewire.LifecycleDeleteResult, bool, error) {
	bucketDir := s3a.bucketDir(bucket)
	versionsDir := bucketDir + "/" + object + s3_constants.VersionsFolder
	count := 0
	var firstName string
	err := s3a.WithFilerClient(false, func(client filer_pb.HanzoFilerClient) error {
		return filer_pb.HanzoList(ctx, client, versionsDir, "", func(entry *filer_pb.Entry, _ bool) error {
			count++
			if count == 1 && entry != nil {
				firstName = entry.Name
			}
			return nil
		}, "", false, 2)
	})
	if err != nil {
		if errors.Is(err, filer_pb.ErrNotFound) {
			return noopResolved("NOT_FOUND"), true, nil
		}
		return retryLater("TRANSPORT_ERROR: sole-survivor list: " + err.Error()), true, nil
	}
	if count == 0 {
		return noopResolved("NOT_FOUND"), true, nil
	}
	if count > 1 {
		return noopResolved("NOT_SOLE_SURVIVOR"), true, nil
	}
	// HanzoList delivered a single callback but with a nil entry; we
	// can't compare names so retry rather than silently bypass the
	// marker-replaced check.
	if firstName == "" {
		return retryLater("PENDING_SURVIVOR_ENTRY"), true, nil
	}
	if versionId != "" && firstName != s3a.getVersionFileName(versionId) {
		return noopResolved("MARKER_REPLACED"), true, nil
	}
	// Latest-pointer check: createDeleteMarker writes the marker file
	// and then updates the parent directory's Extended map. Reading
	// before the second step lands would see count==1 but no pointer;
	// retry-later rather than mistakenly delete.
	parent, name := path.Split(versionsDir)
	parent = strings.TrimRight(parent, "/")
	if parent == "" {
		parent = "/"
	}
	versionsEntry, err := s3a.getEntry(parent, name)
	if err != nil {
		if errors.Is(err, filer_pb.ErrNotFound) {
			return noopResolved("NOT_FOUND"), true, nil
		}
		return retryLater("TRANSPORT_ERROR: latest-pointer lookup: " + err.Error()), true, nil
	}
	if versionsEntry == nil {
		return noopResolved("NOT_FOUND"), true, nil
	}
	latest, hasPointer := versionsEntry.Extended[s3_constants.ExtLatestVersionIdKey]
	if !hasPointer || len(latest) == 0 {
		return retryLater("PENDING_LATEST_POINTER"), true, nil
	}
	if string(latest) != versionId {
		return noopResolved("MARKER_NOT_LATEST"), true, nil
	}
	// Null-version check: pre-versioning objects survive as the bare
	// <bucket>/<key>. Both regular files and explicit S3 directory-key
	// markers (object names ending in /) qualify; the listing path
	// (s3api_object_versioning.go processExplicitDirectory) treats both
	// as the null version. getEntry uses NewFullPath so a trailing slash
	// in object splits the same as a regular key.
	bareEntry, err := s3a.getEntry(bucketDir, object)
	if err != nil {
		if errors.Is(err, filer_pb.ErrNotFound) {
			return s3_lifecyclewire.LifecycleDeleteResult{}, false, nil
		}
		return retryLater("TRANSPORT_ERROR: null-version lookup: " + err.Error()), true, nil
	}
	if bareEntry != nil && (!bareEntry.IsDirectory || bareEntry.IsDirectoryKeyObject()) {
		return noopResolved("NULL_VERSION_PRESENT"), true, nil
	}
	return s3_lifecyclewire.LifecycleDeleteResult{}, false, nil
}

// isCurrentLatestVersion reports whether versionId is the version the
// .versions/ directory currently points to. Hanzo records the latest
// version on the parent directory's Extended map; without consulting it,
// a noncurrent-kind dispatch can't safely distinguish current from
// noncurrent and would risk deleting the live version. Returns
// (false, nil) when the directory has no latest pointer (e.g., the
// bucket isn't versioned in this object's history).
func (s3a *S3ApiServer) isCurrentLatestVersion(bucket, object, versionId string) (bool, error) {
	versionsDir := s3a.bucketDir(bucket) + "/" + object + s3_constants.VersionsFolder
	parent, name := path.Split(versionsDir)
	parent = strings.TrimRight(parent, "/")
	if parent == "" {
		parent = "/"
	}
	entry, err := s3a.getEntry(parent, name)
	if err != nil {
		return false, err
	}
	if entry == nil || len(entry.Extended) == 0 {
		return false, nil
	}
	latest, ok := entry.Extended[s3_constants.ExtLatestVersionIdKey]
	if !ok {
		return false, nil
	}
	return string(latest) == versionId, nil
}

// computeEntryIdentity captures (mtime, size, head fid, sorted-Extended hash):
// an overwrite changes mtime/size/fid; a metadata edit changes Extended; a
// snapshot-restore that preserves mtime+size still differs in head_fid.
func computeEntryIdentity(entry *filer_pb.Entry) *entryIdentity {
	if entry == nil {
		return nil
	}
	id := &entryIdentity{}
	if entry.Attributes != nil {
		// FuseAttributes splits the timestamp across Mtime (seconds) and
		// MtimeNs (nanosecond component); EntryIdentity.MtimeNs is the
		// combined nanoseconds-since-epoch value.
		id.MtimeNs = entry.Attributes.Mtime*int64(1e9) + int64(entry.Attributes.MtimeNs)
		id.Size = int64(entry.Attributes.FileSize)
	}
	if len(entry.GetChunks()) > 0 {
		id.HeadFid = entry.GetChunks()[0].GetFileIdString()
	}
	id.ExtendedHash = s3lifecycle.HashExtended(entry.Extended)
	return id
}

// identityMatches compares the live entry identity against the request's
// CAS witness (req.ExpectedIdentity). A nil/empty witness (early bootstrap)
// always matches.
func identityMatches(live *entryIdentity, req s3_lifecyclewire.LifecycleDeleteRequest) bool {
	// No CAS witness (early bootstrap): null nested object. Must gate before
	// touching any accessor, which would panic on the nil backing message.
	if !req.HasExpectedIdentity() {
		return true
	}
	want := req.ExpectedIdentity()
	if live == nil {
		return false
	}
	if live.MtimeNs != want.MtimeNs() || live.Size != want.Size() {
		return false
	}
	if live.HeadFid != want.HeadFid() {
		return false
	}
	return bytes.Equal(live.ExtendedHash, want.ExtendedHash())
}

// entryUsesMetadataOnlyDelete reports whether the lifecycle delete path
// can skip per-chunk DeleteFile RPCs and rely on the volume's TTL to
// reclaim chunks. Per-write TTL stamping (PR 9377) sets Attributes.TtlSec
// on every entry whose lifecycle rule fits within volume TTL — observing
// a non-zero TtlSec on the live entry is the authoritative signal.
// Defensive nil-checks because the caller may be racing a concurrent
// rewrite that nil-ed Attributes briefly during meta-log replay.
func entryUsesMetadataOnlyDelete(entry *filer_pb.Entry) bool {
	return entry != nil && entry.Attributes != nil && entry.Attributes.TtlSec > 0
}

// recordMetadataOnlyIf bumps the metadata-only counter when on=true.
// Skipped when off so callers don't need a guard at every call site.
// rule_hash is hex-encoded so operators can group by rule when
// debugging; nil rule_hash collapses to the empty string.
func recordMetadataOnlyIf(on bool, req s3_lifecyclewire.LifecycleDeleteRequest) {
	if !on {
		return
	}
	stats_collect.S3LifecycleMetadataOnlyCounter.WithLabelValues(req.Bucket(), hex.EncodeToString(req.RuleHash())).Inc()
}

// actionKindLabel renders an ActionKind* constant as its proto enum name
// for error messages, matching the prior pb stringer output.
func actionKindLabel(k uint32) string {
	switch k {
	case s3_lifecyclewire.ActionKindExpirationDays:
		return "EXPIRATION_DAYS"
	case s3_lifecyclewire.ActionKindExpirationDate:
		return "EXPIRATION_DATE"
	case s3_lifecyclewire.ActionKindNoncurrentDays:
		return "NONCURRENT_DAYS"
	case s3_lifecyclewire.ActionKindNewerNoncurrent:
		return "NEWER_NONCURRENT"
	case s3_lifecyclewire.ActionKindAbortMpu:
		return "ABORT_MPU"
	case s3_lifecyclewire.ActionKindExpiredDeleteMarker:
		return "EXPIRED_DELETE_MARKER"
	default:
		return "ACTION_KIND_UNSPECIFIED"
	}
}

func done() s3_lifecyclewire.LifecycleDeleteResult {
	return s3_lifecyclewire.LifecycleDeleteResult{Outcome: s3_lifecyclewire.LifecycleDeleteOutcomeDone}
}
func noopResolved(reason string) s3_lifecyclewire.LifecycleDeleteResult {
	return s3_lifecyclewire.LifecycleDeleteResult{Outcome: s3_lifecyclewire.LifecycleDeleteOutcomeNoopResolved, Reason: reason}
}
func blocked(reason string) s3_lifecyclewire.LifecycleDeleteResult {
	return s3_lifecyclewire.LifecycleDeleteResult{Outcome: s3_lifecyclewire.LifecycleDeleteOutcomeBlocked, Reason: reason}
}
func retryLater(reason string) s3_lifecyclewire.LifecycleDeleteResult {
	return s3_lifecyclewire.LifecycleDeleteResult{Outcome: s3_lifecyclewire.LifecycleDeleteOutcomeRetryLater, Reason: reason}
}
