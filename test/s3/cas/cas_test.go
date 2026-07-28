// Package cas exercises the one property every conditional write rests on:
// a precondition and the write it guards must commit as a unit, so that N
// racers conditioned on the SAME version produce EXACTLY ONE winner.
//
// This is not a restatement of test/s3/versioning/s3_conditional_writes_test.go.
// That suite asks whether a precondition is EVALUATED correctly — stale ETag
// rejected, current ETag accepted — one request at a time, and it passed
// throughout the period the gateway was admitting two concurrent writers. The
// fault lived entirely in the gap between evaluating the precondition and
// committing the body, which only concurrency can open.
//
// Why it matters beyond S3 semantics: hanzoai/cloud fences per-tenant SQLite on
// this store, and its boot probe (org.ProbeCAS) fires four concurrent
// PutIfVersion calls and counts winners. Two winners means the fence can admit
// two writers for one org — split-brain — so cloud refuses to use the store at
// all and falls back to local-only. Production logged exactly that:
//
//	durability disabled — object-store conditional-PUT atomicity NOT confirmed —
//	update-race admitted 2 writers (want 1) — If-Match not atomic
//
// Four racers is the smallest race that can detect the fault, not a level that
// reliably provokes it. These tests run wider and repeat, because a single green
// round of a nondeterministic race is not evidence.
//
// # WHAT THESE TESTS DO NOT COVER, and why it matters
//
// Everything here signs its payload the ordinary way. That is not the only way an
// S3 client writes: over plain HTTP, minio-go — and so hanzoai/s3-go, which cloud's
// fence uses — signs with STREAMING-AWS4-HMAC-SHA256-PAYLOAD and sends an
// aws-chunked body. The two spellings take DIFFERENT paths through this gateway,
// and only one of them is atomic. Measured against both the deployed 4.34.6 and a
// build of this tree:
//
//	plain signature      15,700 conditional PUTs   0 over-admitted
//	streaming signature      ~300 races            ~30-40% over-admitted, up to 3 of 4
//
// So a green run here says nothing about the streaming path, and the fence runs on
// the streaming path. Covering it needs a client that emits aws-chunked bodies;
// aws-sdk-go-v2 will not do it for a seekable body over http. Until then, the
// reproducer is cloud's own TestProbeCASSoak_Staging, which drives the real client.
//
// Run against a live gateway:
//
//	s3 mini -dir=/tmp/d -s3.port=8333 -s3.config=docker/compose/s3.json
//	go test ./test/s3/cas/ -v
//
// Knobs: S3_ENDPOINT, S3_ACCESS_KEY, S3_SECRET_KEY, CAS_BUCKET, CAS_RACERS,
// CAS_ROUNDS, CAS_KEY=fresh, CAS_ETAG=raw.
package cas

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil && v > 0 {
		return v
	}
	return fallback
}

func client(t *testing.T) *s3.Client {
	t.Helper()
	endpoint := env("S3_ENDPOINT", "http://localhost:8333")
	c := s3.New(s3.Options{
		Region:       env("S3_REGION", "us-east-1"),
		BaseEndpoint: aws.String(endpoint),
		UsePathStyle: true,
		Credentials: credentials.NewStaticCredentialsProvider(
			env("S3_ACCESS_KEY", "some_access_key1"),
			env("S3_SECRET_KEY", "some_secret_key1"), ""),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := c.ListBuckets(ctx, &s3.ListBucketsInput{}); err != nil {
		t.Skipf("no S3 gateway at %s (%v) — start one with: s3 mini -dir=/tmp/d -s3.port=8333 -s3.config=docker/compose/s3.json", endpoint, err)
	}
	return c
}

// bucket returns somewhere to race in. By default it creates a throwaway one;
// CAS_BUCKET points the run at an existing bucket instead, which is how it runs
// against a live gateway where creating buckets is not ours to do. Keys are
// prefixed per run either way, so concurrent runs never collide.
func bucket(t *testing.T, c *s3.Client) string {
	t.Helper()
	if existing := os.Getenv("CAS_BUCKET"); existing != "" {
		return existing
	}
	name := fmt.Sprintf("cas-%d", time.Now().UnixNano())
	if _, err := c.CreateBucket(context.Background(), &s3.CreateBucketInput{Bucket: aws.String(name)}); err != nil {
		t.Fatalf("create bucket %s: %v", name, err)
	}
	t.Cleanup(func() {
		_, _ = c.DeleteBucket(context.Background(), &s3.DeleteBucketInput{Bucket: aws.String(name)})
	})
	return name
}

// scope namespaces a run's keys. It also mirrors what cloud's own probe writes
// (".probe/cas-"), so a run against the live gateway lands where that bucket's
// lifecycle rule already expects to reap probe objects rather than beside tenant
// data.
func scope(t *testing.T, c *s3.Client, bkt string) func(string) string {
	t.Helper()
	prefix := fmt.Sprintf(".probe/cas-%d/", time.Now().UnixNano())
	t.Cleanup(func() {
		out, err := c.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
			Bucket: aws.String(bkt), Prefix: aws.String(prefix),
		})
		if err != nil {
			return
		}
		for _, o := range out.Contents {
			_, _ = c.DeleteObject(context.Background(), &s3.DeleteObjectInput{Bucket: aws.String(bkt), Key: o.Key})
		}
	})
	return func(name string) string { return prefix + name }
}

// A precondition is the only thing that varies between the three races below, so
// it is the only thing parameterised — the concurrency, timing and accounting
// stay in one place where they cannot drift apart.
type precondition func(*s3.PutObjectInput)

func mustNotExist(in *s3.PutObjectInput) { in.IfNoneMatch = aws.String("*") }
func unconditional(*s3.PutObjectInput)   {}

// atVersion conditions a write on an ETag. RFC 9110 entity-tags are quoted, and
// that is how the AWS SDK hands them back — but S3 clients differ: minio-go (and
// so hanzoai/s3-go, which cloud's fence uses) strips the quotes before putting the
// value in the header. Both spellings name the same version, so a gateway must
// treat them identically. CAS_ETAG=raw sends the unquoted form to check that it
// does, because a gateway that fails to PARSE a precondition and then proceeds as
// if none was given turns a compare-and-swap into a blind overwrite.
func atVersion(etag string) precondition {
	if os.Getenv("CAS_ETAG") == "raw" {
		etag = strings.Trim(etag, `"`)
	}
	return func(in *s3.PutObjectInput) { in.IfMatch = aws.String(etag) }
}

// status maps an SDK error to the HTTP status the gateway returned, so a race
// result can distinguish the three outcomes that matter: 200 committed, 412
// correctly refused by the precondition, and anything else — the gateway failing
// rather than deciding.
func status(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var resp *awshttp.ResponseError
	if errors.As(err, &resp) {
		return resp.HTTPStatusCode()
	}
	var api smithy.APIError
	if errors.As(err, &api) && api.ErrorCode() == "PreconditionFailed" {
		return http.StatusPreconditionFailed
	}
	return -1
}

type round struct {
	winners  int
	winnerID string // ETag of whoever committed
	codes    map[int]int
	slowest  time.Duration
}

// race fires n concurrent PUTs at one key, every racer released at the same
// instant, and reports how many the gateway admitted.
func race(ctx context.Context, c *s3.Client, bkt, key string, cond precondition, n int) round {
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	r := round{codes: map[int]int{}}
	start := make(chan struct{})
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			body := []byte(fmt.Sprintf("racer-%d-%d", i, time.Now().UnixNano()))
			in := &s3.PutObjectInput{
				Bucket: aws.String(bkt),
				Key:    aws.String(key),
				Body:   bytes.NewReader(body),
			}
			cond(in)
			<-start // release every racer at once
			t0 := time.Now()
			out, err := c.PutObject(ctx, in)
			elapsed := time.Since(t0)

			mu.Lock()
			defer mu.Unlock()
			r.codes[status(err)]++
			if elapsed > r.slowest {
				r.slowest = elapsed
			}
			if err == nil {
				r.winners++
				if out.ETag != nil {
					r.winnerID = *out.ETag
				}
			}
		}(i)
	}
	close(start)
	wg.Wait()
	return r
}

func fmtCodes(codes map[int]int) string {
	keys := make([]int, 0, len(codes))
	for k := range codes {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		label := strconv.Itoa(k)
		if k == -1 {
			label = "transport-error"
		}
		parts = append(parts, fmt.Sprintf("%s×%d", label, codes[k]))
	}
	return strings.Join(parts, " ")
}

// checkCodes reports any racer that neither committed nor was refused by the
// precondition. A loser is supposed to get 412; a 500 means the gateway could not
// decide, which is safe (it did not write) but is still a failure to serve.
func checkCodes(t *testing.T, where string, codes map[int]int) int {
	t.Helper()
	hard := 0
	for code, n := range codes {
		if code != http.StatusOK && code != http.StatusPreconditionFailed {
			hard += n
			t.Errorf("%s: %d racer(s) got HTTP %d — a loser must be refused by the precondition (412), not fail", where, n, code)
		}
	}
	return hard
}

// TestCASUpdateRace is the load-bearing one: every round conditions all racers on
// the CURRENT version and requires exactly one to advance it. Chaining rounds on
// the previous winner's ETag means a round that admitted two writers also
// corrupts the next round's expectation, so the fault cannot hide in the count.
//
// CAS_KEY=fresh mints a new key per round and conditions on a version read back
// with a GET, which is what cloud's fence does (org.ProbeCAS: create-race →
// read-back → update-race). It is deliberately a separate mode, because whether
// the object under a conditional PUT was created seconds ago or has been rewritten
// a hundred times turned out to change the answer — see the package comment.
func TestCASUpdateRace(t *testing.T) {
	c := client(t)
	bkt := bucket(t, c)
	racers, rounds := envInt("CAS_RACERS", 16), envInt("CAS_ROUNDS", 25)
	name := scope(t, c, bkt)
	fresh := os.Getenv("CAS_KEY") == "fresh"
	key := name("update-race")
	ctx := context.Background()

	seed, err := c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bkt), Key: aws.String(key), Body: bytes.NewReader([]byte("seed")),
	})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	version := *seed.ETag

	var slowest time.Duration
	overAdmitted, hardErrors := 0, 0
	for i := 1; i <= rounds; i++ {
		if fresh {
			key = name(fmt.Sprintf("update-race-%d", i))
			created := race(ctx, c, bkt, key, mustNotExist, racers)
			if created.winners != 1 {
				overAdmitted++
				t.Errorf("round %d/%d: create-race admitted %d writers, want 1 [%s]", i, rounds, created.winners, fmtCodes(created.codes))
			}
			head, err := c.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(bkt), Key: aws.String(key)})
			if err != nil {
				t.Fatalf("round %d: read back %s: %v", i, key, err)
			}
			version = *head.ETag
		}
		r := race(ctx, c, bkt, key, atVersion(version), racers)
		if r.slowest > slowest {
			slowest = r.slowest
		}
		if r.winners != 1 {
			overAdmitted++
			t.Errorf("round %d/%d: %d racers conditioned on %s admitted %d writers, want exactly 1 — If-Match is not atomic [%s]",
				i, rounds, racers, version, r.winners, fmtCodes(r.codes))
		}
		hardErrors += checkCodes(t, fmt.Sprintf("round %d/%d", i, rounds), r.codes)
		if r.winnerID == "" {
			t.Fatalf("round %d: no winner, cannot continue the chain [%s]", i, fmtCodes(r.codes))
		}
		version = r.winnerID
	}
	t.Logf("update-race: %d rounds × %d racers = %d conditional PUTs; over-admitted rounds=%d hard errors=%d slowest racer=%v",
		rounds, racers, rounds*racers, overAdmitted, hardErrors, slowest)
}

// TestCASCreateRace covers the other precondition the fence uses: If-None-Match:*
// on a key that does not exist yet. Each round uses a fresh key, so every round is
// an independent create race.
func TestCASCreateRace(t *testing.T) {
	c := client(t)
	bkt := bucket(t, c)
	racers, rounds := envInt("CAS_RACERS", 16), envInt("CAS_ROUNDS", 25)
	key := scope(t, c, bkt)
	ctx := context.Background()

	overAdmitted := 0
	for i := 1; i <= rounds; i++ {
		where := fmt.Sprintf("round %d/%d", i, rounds)
		r := race(ctx, c, bkt, key(fmt.Sprintf("create-race-%d", i)), mustNotExist, racers)
		if r.winners != 1 {
			overAdmitted++
			t.Errorf("%s: %d racers creating one fresh key admitted %d writers, want exactly 1 — If-None-Match is not atomic [%s]",
				where, racers, r.winners, fmtCodes(r.codes))
		}
		checkCodes(t, where, r.codes)
	}
	t.Logf("create-race: %d rounds × %d racers; over-admitted rounds=%d", rounds, racers, overAdmitted)
}

// TestUnconditionalPutStaysConcurrent guards the other side of the fix. Making
// conditional writes safe would cost too much if it also serialized ordinary
// traffic: an unconditional PUT has no precondition to be atomic with,
// last-writer-wins is the defined S3 semantic, and taking a distributed lock per
// object write would be a large, silent throughput regression. Every racer must
// commit.
func TestUnconditionalPutStaysConcurrent(t *testing.T) {
	c := client(t)
	bkt := bucket(t, c)
	racers := envInt("CAS_RACERS", 16)
	r := race(context.Background(), c, bkt, scope(t, c, bkt)("unconditional"), unconditional, racers)
	if r.winners != racers {
		t.Errorf("%d unconditional PUTs admitted only %d writers [%s] — last-writer-wins must not be serialized behind a lock",
			racers, r.winners, fmtCodes(r.codes))
	}
	t.Logf("unconditional: %d/%d committed, slowest racer=%v", r.winners, racers, r.slowest)
}
