#!/usr/bin/env python3
"""Patch the compiled React bundle to update Login page text.

Replaces marketing copy in the minified JS bundle since the MDS npm
package is unavailable and we cannot rebuild from source.
"""

import sys

BUNDLE = "/Users/z/work/hanzo/storage-console-clean/web-app/build/static/js/main.4bcf28ca.js"

# Each tuple: (old_string, new_string)
REPLACEMENTS = [
    # --- Hero section ---
    (
        "Unified Object Storage for AI Infrastructure",
        "scale, speed, and security",
    ),
    (
        "S3-compatible object storage with enterprise SSO, multi-tenant isolation, and global CDN delivery. Built for teams that ship.",
        "Store and retrieve any amount of data, anytime, from anywhere. S3-compatible, enterprise-grade, with SSO and multi-tenant isolation built in.",
    ),
    (
        "Sign In with Hanzo ID",
        "Get started",
    ),

    # --- Stats (value/label pairs in minified object literals) ---
    (
        'value:"Multi",label:"Tenant"',
        'value:"11\\u00d79s",label:"Durability"',
    ),
    (
        'value:"99.9%",label:"Uptime"',
        'value:"<50ms",label:"Latency"',
    ),
    (
        'label:"Scale"',
        'label:"Scalable"',
    ),

    # --- Features heading ---
    (
        "Everything you need",
        "Everything you need to store data at scale",
    ),

    # --- Feature cards ---
    (
        "SSO via Hanzo ID",
        "Enterprise SSO",
    ),
    (
        "Enterprise single sign-on through hanzo.id with OIDC. No separate credentials needed.",
        "Authenticate with Hanzo ID via OIDC. Centralized identity, fine-grained access policies, and audit logging out of the box.",
    ),
    (
        "Multi-Org Buckets",
        "Multi-Tenant Isolation",
    ),
    (
        "Isolated storage per organization with fine-grained access policies and quotas.",
        "Per-organization buckets with independent access controls, quotas, and lifecycle policies. Built for teams and platforms.",
    ),
    (
        "Global CDN Ready",
        "AI-Optimized Workloads",
    ),
    (
        "Edge-cached delivery through Cloudflare for static assets and media files.",
        "Purpose-built for model weights, training datasets, embeddings, and inference artifacts. High-throughput parallel uploads.",
    ),
    # Note: "Server-Side Encryption" appears twice in the bundle.
    # Target the feature card occurrence by replacing title+desc together.
    (
        'title:"Server-Side Encryption",desc:"AES-256 encryption at rest with KMS integration for key management."',
        'title:"Encryption & Compliance",desc:"AES-256 server-side encryption at rest, TLS in transit, and KMS-managed keys. SOC 2 and GDPR ready."',
    ),
    (
        "Object Versioning",
        "Edge Delivery",
    ),
    (
        "Full version history for every object. Roll back changes and recover deleted files.",
        "Serve assets globally through Cloudflare\u2019s edge network. Automatic caching, custom domains, and signed URLs.",
    ),
    (
        "Drop-in replacement for Amazon S3. Works with every S3 client, SDK, and tool.",
        "Drop-in replacement for Amazon S3. Use the same SDKs, CLI tools, and integrations you already rely on \u2014 zero migration friction.",
    ),

    # --- CTA section ---
    (
        "Ready to get started?",
        "Start storing data in seconds",
    ),
    (
        "Sign in with your Hanzo ID to access your storage.",
        "Sign in with your Hanzo ID. No credit card, no setup wizard.",
    ),
]


def main():
    with open(BUNDLE, "r", encoding="utf-8") as f:
        data = f.read()

    original_len = len(data)
    total_replaced = 0
    failed = []

    for old, new in REPLACEMENTS:
        count = data.count(old)
        if count == 0:
            failed.append(old[:70])
            continue
        data = data.replace(old, new)
        total_replaced += count
        print(f"  [{count}x] {old[:60]!r}")
        print(f"     -> {new[:60]!r}")

    if failed:
        print(f"\nWARNING: {len(failed)} replacement(s) NOT found:")
        for f_str in failed:
            print(f"  MISS: {f_str!r}")

    with open(BUNDLE, "w", encoding="utf-8") as f:
        f.write(data)

    delta = len(data) - original_len
    print(f"\nDone. {total_replaced} replacements made across {len(REPLACEMENTS)} patterns.")
    print(f"File size delta: {delta:+d} bytes ({original_len} -> {len(data)})")

    if failed:
        sys.exit(1)


if __name__ == "__main__":
    main()
