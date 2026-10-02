package gateway

import "strings"

// Names maps the bucket a client names to the upstream bucket that holds it.
//
// Upstream bucket names are one namespace shared by every account in the world,
// so a client's "org-db" cannot be the upstream's "org-db": each local name is
// wrapped in a prefix and suffix the account owns, "org-db" ->
// "hanzo-s3-org-db-532217001883". The wrapping is the whole mapping, so it is
// reversible and needs no table.
type Names struct {
	Prefix string
	Suffix string
}

// Upstream returns the upstream bucket for a local bucket name, or false when
// the local name is not a valid bucket or the wrapped name exceeds 63 bytes.
func (n Names) Upstream(bucket string) (string, bool) {
	if !validBucket(bucket) {
		return "", false
	}
	name := n.Prefix + bucket + n.Suffix
	if len(name) > 63 {
		return "", false
	}
	return name, true
}

// Local returns the local name of an upstream bucket, or false when the
// upstream bucket is not one this mapping produces.
func (n Names) Local(upstream string) (string, bool) {
	if len(upstream) <= len(n.Prefix)+len(n.Suffix) ||
		!strings.HasPrefix(upstream, n.Prefix) || !strings.HasSuffix(upstream, n.Suffix) {
		return "", false
	}
	bucket := upstream[len(n.Prefix) : len(upstream)-len(n.Suffix)]
	if !validBucket(bucket) {
		return "", false
	}
	return bucket, true
}

// validBucket accepts 3 to 63 lowercase letters, digits and hyphens that start
// and end with a letter or digit. Dots are refused: the upstream is addressed
// as <bucket>.<host> over TLS, and a dot in the label breaks the certificate
// match.
func validBucket(b string) bool {
	if len(b) < 3 || len(b) > 63 {
		return false
	}
	for i := 0; i < len(b); i++ {
		c := b[i]
		switch {
		case 'a' <= c && c <= 'z', '0' <= c && c <= '9':
		case c == '-' && i > 0 && i < len(b)-1:
		default:
			return false
		}
	}
	return true
}
