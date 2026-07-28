package shell

import (
	"testing"

	"github.com/hanzoai/s3/s3/s3api/lifecycle_xml"
	"github.com/hanzoai/s3/s3/s3api/s3lifecycle"
)

// parseOneRule reads a shell-written policy back through the path the S3 API
// and the lifecycle worker use, so a format the shell alone understands fails.
func parseOneRule(t *testing.T, noncurrentDays int, expireDeleteMarker bool) *s3lifecycle.Rule {
	t.Helper()
	policy, err := lifecyclePolicy(noncurrentDays, expireDeleteMarker)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	rules, err := lifecycle_xml.ParseCanonical(policy)
	if err != nil {
		t.Fatalf("parse %s: %v", policy, err)
	}
	if len(rules) != 1 {
		t.Fatalf("got %d rules, want 1: %s", len(rules), policy)
	}
	if rules[0].Status != s3lifecycle.StatusEnabled {
		t.Fatalf("got status %q, want %q", rules[0].Status, s3lifecycle.StatusEnabled)
	}
	return rules[0]
}

func TestLifecyclePolicyNoncurrentDaysWithDeleteMarker(t *testing.T) {
	rule := parseOneRule(t, 7, true)

	if rule.NoncurrentVersionExpirationDays != 7 {
		t.Errorf("got NoncurrentVersionExpirationDays %d, want 7", rule.NoncurrentVersionExpirationDays)
	}
	if !rule.ExpiredObjectDeleteMarker {
		t.Errorf("ExpiredObjectDeleteMarker not set")
	}
	kinds := s3lifecycle.RuleActionKinds(rule)
	want := []s3lifecycle.ActionKind{s3lifecycle.ActionKindExpiredDeleteMarker, s3lifecycle.ActionKindNoncurrentDays}
	if len(kinds) != len(want) {
		t.Fatalf("got actions %v, want %v", kinds, want)
	}
	for i := range want {
		if kinds[i] != want[i] {
			t.Fatalf("got actions %v, want %v", kinds, want)
		}
	}
}

func TestLifecyclePolicyNoncurrentDaysOnly(t *testing.T) {
	rule := parseOneRule(t, 30, false)

	if rule.NoncurrentVersionExpirationDays != 30 {
		t.Errorf("got NoncurrentVersionExpirationDays %d, want 30", rule.NoncurrentVersionExpirationDays)
	}
	if rule.ExpiredObjectDeleteMarker {
		t.Errorf("ExpiredObjectDeleteMarker set without -expire-delete-marker")
	}
	if rule.ExpirationDays != 0 {
		t.Errorf("got ExpirationDays %d, want 0", rule.ExpirationDays)
	}
	kinds := s3lifecycle.RuleActionKinds(rule)
	if len(kinds) != 1 || kinds[0] != s3lifecycle.ActionKindNoncurrentDays {
		t.Errorf("got actions %v, want [noncurrent_days]", kinds)
	}
}
