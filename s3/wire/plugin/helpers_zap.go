// Code generated from plugin.proto; DO NOT EDIT.

package pluginwire

import (
	"math"

	zap "github.com/zap-proto/go"
)

// float64bits returns the IEEE-754 bit pattern of f, for packing a double into
// a flat uint64 list element.
func float64bits(f float64) uint64 { return math.Float64bits(f) }

// setChild stores a pre-built child message frame as a singular nested message
// at the given bytes-pointer field offset of ob. A singular nested message and
// an out-of-line list element use the same representation — a complete ZAP frame
// stored out of line — so both are read back by re-parsing (see childObject /
// List.ObjectAt). A zero-length child writes a null pointer.
func setChild(ob *zap.ObjectBuilder, fieldOff int, child []byte) {
	ob.SetBytes(fieldOff, child)
}

// emptyChild is a valid, field-less ZAP frame used as the backing buffer for an
// absent nested message. Reads against its root land past the (zero-length) data
// segment and the runtime returns zero values — so accessing a field of an
// absent singular message yields the zero value instead of dereferencing a nil
// buffer.
var emptyChild = func() []byte {
	b := zap.NewBuilder(zap.HeaderSize)
	b.StartObject(0).FinishAsRoot()
	return b.Finish()
}()

// childObject re-parses a nested message frame into its root Object. A singular
// nested message and an out-of-line list element share one representation — a
// complete ZAP frame stored out of line — so both are read back by re-parsing.
// An empty or malformed frame falls back to the emptyChild sentinel, whose
// scalar/text accessors return zero values, so an absent singular message reads
// as zero and never panics.
func childObject(frame []byte) zap.Object {
	m, err := zap.Parse(frame)
	if err != nil {
		m, _ = zap.Parse(emptyChild)
	}
	return m.Root()
}
