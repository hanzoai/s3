// Code generated from worker.proto; DO NOT EDIT.

package workerwire

// S3LifecycleParams_Subtype enum values. Mapped to uint32 scalars on the wire.
// Zero is an UNSPECIFIED sentinel: a TaskParams payload whose subtype is unset
// must not silently route into a READ task.
const (
	S3LifecycleParamsSubtypeUnspecified uint32 = 0
	S3LifecycleParamsSubtypeRead        uint32 = 1
	S3LifecycleParamsSubtypeBootstrap   uint32 = 2
	S3LifecycleParamsSubtypeDrain       uint32 = 3
)
