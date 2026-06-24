// Code generated from filer.proto; DO NOT EDIT.

package filerwire

// Enum constants for filer.proto. Each proto enum maps to a uint32 scalar on the
// wire; these consts carry the same numeric values the proto declares.

// SSEType values (proto enum SSEType): server-side encryption type of a chunk.
const (
	SSETypeNone   uint32 = 0 // no server-side encryption
	SSETypeSSEC   uint32 = 1 // SSE with customer-provided keys
	SSETypeSSEKMS uint32 = 2 // SSE with KMS-managed keys
	SSETypeSSES3  uint32 = 3 // SSE with S3-managed keys
)

// FilerError values (proto enum FilerError): structured error codes for filer
// entry operations. Values are stable — do not reorder or reuse numbers.
const (
	FilerErrorOK                  uint32 = 0
	FilerErrorEntryNameTooLong    uint32 = 1 // name exceeds max_file_name_length
	FilerErrorParentIsFile        uint32 = 2 // parent path component is a file
	FilerErrorExistingIsDirectory uint32 = 3 // cannot overwrite directory with file
	FilerErrorExistingIsFile      uint32 = 4 // cannot overwrite file with directory
	FilerErrorEntryAlreadyExists  uint32 = 5 // O_EXCL and entry already exists
	FilerErrorPreconditionFailed  uint32 = 6 // WriteCondition not satisfied
)

// PosixLockOp values (proto enum PosixLockOp): advisory lock operation.
const (
	PosixLockOpTryLock           uint32 = 0 // grant lock or report conflict (non-blocking)
	PosixLockOpUnlock            uint32 = 1 // release owner's locks over its range
	PosixLockOpGetLk             uint32 = 2 // report a conflicting lock, if any
	PosixLockOpReleasePosixOwner uint32 = 3 // drop owner's fcntl locks (flush-time)
	PosixLockOpReleaseFlockOwner uint32 = 4 // drop owner's flock locks (release-time)
	PosixLockOpKeepAlive         uint32 = 5 // renew the session's lease on this owner
)

// WriteConditionKind values (proto enum WriteCondition.Kind): the precondition a
// Clause expresses.
const (
	WriteConditionKindNone                  uint32 = 0 // unconditional
	WriteConditionKindIfNotExists           uint32 = 1 // fail if the entry exists
	WriteConditionKindIfExists              uint32 = 2 // fail if the entry is absent
	WriteConditionKindIfETagMatch           uint32 = 3 // fail if absent or etag matches none of the set
	WriteConditionKindIfETagNotMatch        uint32 = 4 // fail if present and etag matches any of the set
	WriteConditionKindIfUnmodifiedSince     uint32 = 5 // fail if present and mtime > unix_time
	WriteConditionKindIfModifiedSince       uint32 = 6 // fail if present and mtime <= unix_time
	WriteConditionKindIfExtendedNotEqual    uint32 = 7 // fail if present and extended[ext_key] == ext_value
	WriteConditionKindIfExtendedTimeElapsed uint32 = 8 // fail if present and extended[ext_key] is in the future
)

// ObjectMutationType values (proto enum ObjectMutation.Type): the kind of
// entry-level change a mutation applies.
const (
	ObjectMutationTypePut             uint32 = 0 // create or replace the entry
	ObjectMutationTypeDelete          uint32 = 1 // delete the entry at directory/name
	ObjectMutationTypePatchExtended   uint32 = 2 // merge set_extended / remove delete_extended
	ObjectMutationTypeRecomputeLatest uint32 = 3 // scan a directory and re-point a parent entry
)

// Oneof discriminator tags. A oneof maps to a uint32 "which" tag (0 = none set)
// plus the selected variant object; the tag value equals the proto field number
// of the selected member.

// StreamMutateEntryRequest.request members (proto field numbers).
const (
	StreamMutateRequestNone   uint32 = 0
	StreamMutateRequestCreate uint32 = 2
	StreamMutateRequestUpdate uint32 = 3
	StreamMutateRequestDelete uint32 = 4
	StreamMutateRequestRename uint32 = 5
)

// StreamMutateEntryResponse.response members (proto field numbers).
const (
	StreamMutateResponseNone   uint32 = 0
	StreamMutateResponseCreate uint32 = 3
	StreamMutateResponseUpdate uint32 = 4
	StreamMutateResponseDelete uint32 = 5
	StreamMutateResponseRename uint32 = 6
)
