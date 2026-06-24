// Code generated from volume_server.proto; DO NOT EDIT.

package volume_serverwire

// Enum constants for volume_server.proto. Each proto enum maps to a uint32
// scalar on the wire; these consts carry the same numeric values the proto
// declares.

// ChecksumAlgorithm values (proto enum ChecksumAlgorithm): the checksum
// algorithm of an EC bitrot sidecar.
const (
	ChecksumNone   uint32 = 0
	ChecksumCRC32C uint32 = 1 // CRC32C (Castagnoli)
)

// VolumeScrubMode values (proto enum VolumeScrubMode): the depth/kind of a
// scrub pass.
const (
	VolumeScrubModeUnknown  uint32 = 0
	VolumeScrubModeIndex    uint32 = 1
	VolumeScrubModeFull     uint32 = 2
	VolumeScrubModeLocal    uint32 = 3
	VolumeScrubModeChecksum uint32 = 4 // EC only: verify shard bytes against the bitrot sidecar
)

// Oneof discriminator tags. A oneof maps to a uint32 "which" tag (0 = none set)
// plus the selected variant payload; the tag value equals the proto field
// number of the selected member.

// ReceiveFileRequest.data members (proto field numbers).
const (
	ReceiveFileDataNone        uint32 = 0
	ReceiveFileDataInfo        uint32 = 1
	ReceiveFileDataFileContent uint32 = 2
)
