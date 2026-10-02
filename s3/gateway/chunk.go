package gateway

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"hash/crc32"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hanzos3/crc64nvme"
)

// trailerPayload is the payload hash of an aws-chunked body whose chunks carry
// no signatures and whose last chunk carries a checksum of the whole.
const trailerPayload = "STREAMING-UNSIGNED-PAYLOAD-TRAILER"

// chunkSize is the size of every chunk but the last in a trailed body.
const chunkSize = 64 << 10

// checksum is an algorithm a trailed body's last chunk carries.
type checksum struct {
	trailer string // the trailer header, "x-amz-checksum-crc32"
	name    string // the x-amz-sdk-checksum-algorithm value, "CRC32"
	size    int    // bytes in the sum
	new     func() hash.Hash
}

// checksums are the algorithms the upstream verifies, by trailer header.
var checksums = map[string]checksum{
	"x-amz-checksum-crc32":     {"x-amz-checksum-crc32", "CRC32", 4, func() hash.Hash { return crc32.NewIEEE() }},
	"x-amz-checksum-crc32c":    {"x-amz-checksum-crc32c", "CRC32C", 4, func() hash.Hash { return crc32.New(crc32.MakeTable(crc32.Castagnoli)) }},
	"x-amz-checksum-crc64nvme": {"x-amz-checksum-crc64nvme", "CRC64NVME", 8, func() hash.Hash { return crc64nvme.New() }},
	"x-amz-checksum-sha1":      {"x-amz-checksum-sha1", "SHA1", sha1.Size, sha1.New},
	"x-amz-checksum-sha256":    {"x-amz-checksum-sha256", "SHA256", sha256.Size, sha256.New},
}

// checksumFor returns the algorithm to trail a body with: the one the client
// trailed its own body with, since a multipart upload created for that
// algorithm refuses a part carrying another, and CRC32 otherwise.
func checksumFor(clientTrailer string) checksum {
	if c, ok := checksums[strings.ToLower(strings.TrimSpace(clientTrailer))]; ok {
		return c
	}
	return checksums["x-amz-checksum-crc32"]
}

// hasChecksum reports whether h already carries a checksum of the body:
// Content-MD5 or an x-amz-checksum-<algorithm> value.
func hasChecksum(h http.Header) bool {
	if h.Get("Content-Md5") != "" {
		return true
	}
	for name := range h {
		canon := http.CanonicalHeaderKey(name)
		if strings.HasPrefix(canon, "X-Amz-Checksum-") &&
			canon != "X-Amz-Checksum-Algorithm" && canon != "X-Amz-Checksum-Type" {
			return true
		}
	}
	return false
}

// trailedLength is the length on the wire of an n-byte body sent trailed
// with c.
func trailedLength(n int64, c checksum) int64 {
	full, rest := n/chunkSize, n%chunkSize
	total := full * (int64(len(strconv.FormatInt(chunkSize, 16))) + 2 + chunkSize + 2)
	if rest > 0 {
		total += int64(len(strconv.FormatInt(rest, 16))) + 2 + rest + 2
	}
	// "0\r\n" + "<trailer>:" + base64 sum + "\r\n" + "\r\n"
	return total + 3 + int64(len(c.trailer)) + 1 + int64(base64.StdEncoding.EncodedLen(c.size)) + 2 + 2
}

// errPayloadHash is returned when a body does not hash to the SHA-256 its
// sender declared. The upstream request fails with it, so nothing is stored.
var errPayloadHash = errors.New("payload does not match its declared x-amz-content-sha256")

// errShortBody is returned when a body ends before its declared length.
var errShortBody = errors.New("payload shorter than its declared length")

// trailed encodes an n-byte body as aws-chunked with a checksum trailer, so
// the upstream checks every byte it stores against a sum computed here. When
// want is a SHA-256 the sender declared, the body is hashed as it passes and
// the last chunk is never written unless it matches.
type trailed struct {
	src  io.Reader
	left int64
	algo checksum
	crc  hash.Hash
	sum  hash.Hash
	want string
	buf  []byte
	out  bytes.Buffer
	done bool
}

func newTrailed(src io.Reader, n int64, want string, algo checksum) *trailed {
	t := &trailed{src: src, left: n, algo: algo, crc: algo.new(), want: want, buf: make([]byte, chunkSize)}
	if want != "" {
		t.sum = sha256.New()
	}
	return t
}

func (t *trailed) Read(p []byte) (int, error) {
	for t.out.Len() == 0 {
		if t.done {
			return 0, io.EOF
		}
		if err := t.fill(); err != nil {
			return 0, err
		}
	}
	return t.out.Read(p)
}

// fill encodes the next chunk, or the final chunk and its trailer.
func (t *trailed) fill() error {
	if t.left == 0 {
		if t.sum != nil && !strings.EqualFold(hex.EncodeToString(t.sum.Sum(nil)), t.want) {
			return errPayloadHash
		}
		t.out.WriteString("0\r\n" + t.algo.trailer + ":" + base64.StdEncoding.EncodeToString(t.crc.Sum(nil)) + "\r\n\r\n")
		t.done = true
		return nil
	}
	size := int64(chunkSize)
	if t.left < size {
		size = t.left
	}
	if _, err := io.ReadFull(t.src, t.buf[:size]); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return errShortBody
		}
		return err
	}
	t.left -= size
	t.crc.Write(t.buf[:size])
	if t.sum != nil {
		t.sum.Write(t.buf[:size])
	}
	t.out.WriteString(strconv.FormatInt(size, 16) + "\r\n")
	t.out.Write(t.buf[:size])
	t.out.WriteString("\r\n")
	return nil
}

// trail sets the headers that announce a body of n decoded bytes trailed
// with c.
func trail(h http.Header, n int64, c checksum) {
	encoding := "aws-chunked"
	if rest := h.Get("Content-Encoding"); rest != "" {
		encoding += "," + rest
	}
	h.Set("Content-Encoding", encoding)
	h.Set("X-Amz-Decoded-Content-Length", strconv.FormatInt(n, 10))
	h.Set("X-Amz-Trailer", c.trailer)
	h.Set("X-Amz-Sdk-Checksum-Algorithm", c.name)
}

// contentMD5 is the Content-MD5 value of b.
func contentMD5(b []byte) string {
	sum := md5.Sum(b)
	return base64.StdEncoding.EncodeToString(sum[:])
}

// sha256Hex is the hex SHA-256 of b.
func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// isHash reports whether v is a hex SHA-256, as opposed to one of the
// UNSIGNED-PAYLOAD or STREAMING-* markers.
func isHash(v string) bool {
	if len(v) != 64 {
		return false
	}
	_, err := hex.DecodeString(v)
	return err == nil
}
