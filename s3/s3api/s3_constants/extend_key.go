package s3_constants

const (
	ExtAmzOwnerKey              = "Hanzo-X-Amz-Owner"
	ExtAmzAclKey                = "Hanzo-X-Amz-Acl"
	ExtOwnershipKey             = "Hanzo-X-Amz-Ownership"
	ExtVersioningKey            = "Hanzo-X-Amz-Versioning"
	ExtVersionIdKey             = "Hanzo-X-Amz-Version-Id"
	ExtDeleteMarkerKey          = "Hanzo-X-Amz-Delete-Marker"
	ExtIsLatestKey              = "Hanzo-X-Amz-Is-Latest"
	ExtETagKey                  = "Hanzo-X-Amz-ETag"
	ExtLatestVersionIdKey       = "Hanzo-X-Amz-Latest-Version-Id"
	ExtLatestVersionFileNameKey = "Hanzo-X-Amz-Latest-Version-File-Name"
	ExtAllowEmptyFolders        = "Hanzo-X-Amz-Allow-Empty-Folders"
	// Cached list metadata in .versions directory for single-scan efficiency
	ExtLatestVersionSizeKey        = "Hanzo-X-Amz-Latest-Version-Size"
	ExtLatestVersionETagKey        = "Hanzo-X-Amz-Latest-Version-ETag"
	ExtLatestVersionMtimeKey       = "Hanzo-X-Amz-Latest-Version-Mtime"
	ExtLatestVersionOwnerKey       = "Hanzo-X-Amz-Latest-Version-Owner"
	ExtLatestVersionIsDeleteMarker = "Hanzo-X-Amz-Latest-Version-Is-Delete-Marker"
	ExtMultipartObjectKey          = "key"
	// Wall-clock nanoseconds (int64 as decimal string) captured at the
	// moment a versioned entry was demoted from current to noncurrent
	// by a later PUT or delete marker. Read by the s3 lifecycle engine
	// to compute NoncurrentDays due time; zero/missing falls back to
	// the entry's own mtime so legacy data still expires.
	ExtNoncurrentSinceNsKey = "Hanzo-X-Amz-Noncurrent-Since-Ns"

	// Per-bucket opt-in for the PutObject lifecycle TTL fast path ("true"
	// to enable). When on, an Expiration.Days rule is stamped as a volume
	// TTL at write time instead of being expired by the worker. Off by
	// default: a baked-in TTL can't honor a later policy change (rule
	// removed or lengthened) the way worker-driven expiration does.
	ExtLifecycleTtlFastPathKey = "Hanzo-X-Amz-Lifecycle-Ttl-Fast-Path"

	// S3 checksum storage keys (use x-hanzo- prefix to avoid leaking in generic header loop)
	ExtChecksumAlgorithm = "x-hanzo-checksum-algorithm"
	ExtChecksumValue     = "x-hanzo-checksum-value"

	// Bucket Policy
	ExtBucketPolicyKey = "Hanzo-X-Amz-Bucket-Policy"

	// Object Retention and Legal Hold
	ExtObjectLockModeKey     = "Hanzo-X-Amz-Object-Lock-Mode"
	ExtRetentionUntilDateKey = "Hanzo-X-Amz-Retention-Until-Date"
	ExtLegalHoldKey          = "Hanzo-X-Amz-Legal-Hold"
	ExtObjectLockEnabledKey  = "Hanzo-X-Amz-Object-Lock-Enabled"

	// Object Lock Bucket Configuration (individual components, not XML)
	ExtObjectLockDefaultModeKey  = "Lock-Default-Mode"
	ExtObjectLockDefaultDaysKey  = "Lock-Default-Days"
	ExtObjectLockDefaultYearsKey = "Lock-Default-Years"
)

// Object Lock and Retention Constants
const (
	// Retention modes
	RetentionModeGovernance = "GOVERNANCE"
	RetentionModeCompliance = "COMPLIANCE"

	// Legal hold status
	LegalHoldOn  = "ON"
	LegalHoldOff = "OFF"

	// Object lock enabled status
	ObjectLockEnabled = "Enabled"

	// Bucket versioning status
	VersioningEnabled   = "Enabled"
	VersioningSuspended = "Suspended"
)
