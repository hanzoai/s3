// Package polarstore implements a clean-room adaptive object compression engine
// inspired by PolarStore (Hu et al., USENIX FAST '26), adapted for pre-encryption
// object storage paths in Hanzo S3 with hard encryption-boundary guards.
package polarstore

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/zstd"
)

// Codec identifies the compression codec applied to a stored object payload.
type Codec uint8

const (
	// CodecNone indicates no compression was applied (identity bypass).
	CodecNone Codec = 0
	// CodecS2 selects the fast-tier S2 codec (high throughput, low CPU).
	CodecS2 Codec = 1
	// CodecZstd selects the balanced-tier Zstandard codec (high ratio, fast decode).
	CodecZstd Codec = 2
	// CodecGzip selects legacy gzip compression (primarily for compatibility and matched baselines).
	CodecGzip Codec = 3
)

// String returns the canonical lowercase codec identifier.
func (c Codec) String() string {
	switch c {
	case CodecNone:
		return "none"
	case CodecS2:
		return "s2"
	case CodecZstd:
		return "zstd"
	case CodecGzip:
		return "gzip"
	default:
		return fmt.Sprintf("unknown(%d)", uint8(c))
	}
}

// ParseCodec parses a canonical codec name.
func ParseCodec(s string) (Codec, error) {
	switch s {
	case "", "none", "identity":
		return CodecNone, nil
	case "s2":
		return CodecS2, nil
	case "zstd":
		return CodecZstd, nil
	case "gzip":
		return CodecGzip, nil
	default:
		return CodecNone, fmt.Errorf("polarstore: unsupported codec %q", s)
	}
}

var (
	// ErrUnsupportedCodec indicates an unrecognized codec identifier.
	ErrUnsupportedCodec = errors.New("polarstore: unsupported codec")
	// ErrCorruptEnvelope indicates a malformed or truncated framed payload.
	ErrCorruptEnvelope = errors.New("polarstore: corrupt compression envelope")
)

var (
	envelopeMagic = [4]byte{'P', 'L', 'S', 'T'}
	envelopeVer   = byte(1)
)

const envelopeHeaderSize = 4 + 1 + 1 + 4 // magic(4) + version(1) + codec(1) + origLen(4)

var (
	zstdEncoderPool = sync.Pool{
		New: func() any {
			enc, err := zstd.NewWriter(
				nil,
				zstd.WithEncoderLevel(zstd.SpeedDefault),
				zstd.WithEncoderConcurrency(1),
			)
			if err != nil {
				panic(fmt.Sprintf("polarstore: init zstd encoder: %v", err))
			}
			return enc
		},
	}

	zstdDecoderPool = sync.Pool{
		New: func() any {
			dec, err := zstd.NewReader(
				nil,
				zstd.WithDecoderConcurrency(1),
			)
			if err != nil {
				panic(fmt.Sprintf("polarstore: init zstd decoder: %v", err))
			}
			return dec
		},
	}

	gzipWriterPool = sync.Pool{
		New: func() any {
			w, err := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
			if err != nil {
				panic(fmt.Sprintf("polarstore: init gzip writer: %v", err))
			}
			return w
		},
	}

	gzipFastWriterPool = sync.Pool{
		New: func() any {
			w, err := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
			if err != nil {
				panic(fmt.Sprintf("polarstore: init gzip fast writer: %v", err))
			}
			return w
		},
	}

	gzipReaderPool = sync.Pool{
		New: func() any {
			return new(gzip.Reader)
		},
	}
)

// CompressWithCodec encodes src using the requested codec.
// For CodecNone, src is returned unchanged without copying.
func CompressWithCodec(codec Codec, src []byte) ([]byte, error) {
	switch codec {
	case CodecNone:
		return src, nil
	case CodecS2:
		return s2.Encode(nil, src), nil
	case CodecZstd:
		enc := zstdEncoderPool.Get().(*zstd.Encoder)
		defer zstdEncoderPool.Put(enc)
		return enc.EncodeAll(src, make([]byte, 0, len(src)/2+64)), nil
	case CodecGzip:
		return compressGzipPooled(&gzipWriterPool, src)
	default:
		return nil, ErrUnsupportedCodec
	}
}

// CompressGzipFast encodes src using gzip.BestSpeed (matching legacy Hanzo S3 util.GzipData).
func CompressGzipFast(src []byte) ([]byte, error) {
	return compressGzipPooled(&gzipFastWriterPool, src)
}

func compressGzipPooled(pool *sync.Pool, src []byte) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(len(src)/2 + 64)
	w := pool.Get().(*gzip.Writer)
	w.Reset(&buf)
	defer func() {
		w.Reset(io.Discard)
		pool.Put(w)
	}()
	if _, err := w.Write(src); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecompressWithCodec decodes src using the specified codec.
// For CodecNone, src is returned unchanged without copying.
func DecompressWithCodec(codec Codec, src []byte) ([]byte, error) {
	switch codec {
	case CodecNone:
		return src, nil
	case CodecS2:
		return s2.Decode(nil, src)
	case CodecZstd:
		dec := zstdDecoderPool.Get().(*zstd.Decoder)
		defer zstdDecoderPool.Put(dec)
		return dec.DecodeAll(src, nil)
	case CodecGzip:
		gr := gzipReaderPool.Get().(*gzip.Reader)
		if err := gr.Reset(bytes.NewReader(src)); err != nil {
			gzipReaderPool.Put(gr)
			return nil, err
		}
		var out bytes.Buffer
		_, err := io.Copy(&out, gr)
		closeErr := gr.Close()
		gzipReaderPool.Put(gr)
		if err != nil {
			return nil, err
		}
		if closeErr != nil {
			return nil, closeErr
		}
		return out.Bytes(), nil
	default:
		return nil, ErrUnsupportedCodec
	}
}

// MarshalEnvelope wraps a compressed or bypassed payload in a self-describing 10-byte frame.
func MarshalEnvelope(codec Codec, originalSize uint32, payload []byte) []byte {
	frame := make([]byte, envelopeHeaderSize+len(payload))
	copy(frame[0:4], envelopeMagic[:])
	frame[4] = envelopeVer
	frame[5] = byte(codec)
	binary.BigEndian.PutUint32(frame[6:10], originalSize)
	copy(frame[envelopeHeaderSize:], payload)
	return frame
}

// UnmarshalEnvelope parses and validates a self-describing envelope frame.
func UnmarshalEnvelope(frame []byte) (Codec, uint32, []byte, error) {
	if len(frame) < envelopeHeaderSize {
		return CodecNone, 0, nil, ErrCorruptEnvelope
	}
	if frame[0] != envelopeMagic[0] || frame[1] != envelopeMagic[1] ||
		frame[2] != envelopeMagic[2] || frame[3] != envelopeMagic[3] {
		return CodecNone, 0, nil, ErrCorruptEnvelope
	}
	if frame[4] != envelopeVer {
		return CodecNone, 0, nil, fmt.Errorf("%w: unsupported version %d", ErrCorruptEnvelope, frame[4])
	}
	codec := Codec(frame[5])
	switch codec {
	case CodecNone, CodecS2, CodecZstd, CodecGzip:
	default:
		return CodecNone, 0, nil, ErrUnsupportedCodec
	}
	origSize := binary.BigEndian.Uint32(frame[6:10])
	return codec, origSize, frame[envelopeHeaderSize:], nil
}
