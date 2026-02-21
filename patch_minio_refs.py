#!/usr/bin/env python3
"""
Patch compiled JS bundles to replace MinIO/min.io references with Hanzo Space equivalents.

Usage:
    python3 patch_minio_refs.py [CONSOLE_ROOT]

    CONSOLE_ROOT defaults to the current working directory.
    The script globs for all *.js files (excluding *.map) in:
        CONSOLE_ROOT/web-app/build/static/js/

Ordering matters: longer/more-specific strings are replaced first to prevent
double-replacement. Safe patterns (Go API paths, SDK names, JS variable names)
are preserved.
"""

import glob
import os
import re
import sys

# Relative path from console root to the JS build output
JS_DIR = os.path.join("web-app", "build", "static", "js")

# --- replacement pairs, applied in order ---
# Each entry: (label, old, new, is_regex)
REPLACEMENTS = [
    # -- User-visible compound phrases (longest first) --
    ("MinIO Object Store",    "MinIO Object Store",    "Hanzo Space",        False),
    ("MinIO Subscription",    "MinIO Subscription",    "Hanzo Subscription", False),
    ("MinIO SUBNET",          "MinIO SUBNET",          "Hanzo Support",      False),
    ("MinIO Console",         "MinIO Console",         "Hanzo Space",        False),
    ("MinIO Browser",         "MinIO Browser",         "Hanzo Space",        False),
    ("MinIO Operator",        "MinIO Operator",        "Hanzo Operator",     False),
    ("MinIO, Inc.",           "MinIO, Inc.",            "Hanzo AI",           False),
    ("MinIO Inc.",            "MinIO Inc.",             "Hanzo AI",           False),

    # -- SUBNET (standalone, not inside longer words) --
    ("SUBNET (standalone)",   r'(?<![a-zA-Z])SUBNET(?![a-zA-Z_])', "Support", True),

    # -- URLs: specific subdomains first --
    ("docs.min.io",           "docs.min.io",           "docs.hanzo.ai/storage", False),
    ("blog.min.io",           "blog.min.io",           "hanzo.ai/blog",         False),
    ("subnet.min.io",         "subnet.min.io",         "hanzo.ai/support",      False),
    ("dl.min.io",             "dl.min.io",             "hanzo.space",           False),
    ("play.min.io",           "play.min.io",           "hanzo.space",           False),

    # -- URLs: specific paths on min.io --
    ("min.io/product/subnet", "min.io/product/subnet", "hanzo.ai/support",      False),
    ("min.io/downloads",      "min.io/downloads",      "hanzo.ai/downloads",    False),
    ("min.io/signup",         "min.io/signup",          "hanzo.ai/signup",       False),
    ("min.io/pricing",        "min.io/pricing",         "hanzo.ai/pricing",      False),
    ("min.io/community",      "min.io/community",       "hanzo.ai/community",   False),

    # -- URL catch-all: remaining min.io refs --
    ("min.io (catch-all)",    "min.io",                 "hanzo.ai",             False),

    # -- Brand catch-all: remaining MinIO --
    ("MinIO (catch-all)",     "MinIO",                  "Hanzo Space",          False),

    # -- Lowercase selective replacements --
    ("minio console",         "minio console",          "hanzo space",          False),
    ("minio server",          "minio server",           "hanzo server",         False),
    ("minio object",          "minio object",           "hanzo object",         False),
]


def patch_file(path):
    """Apply all replacements to a single JS file. Returns total replacement count."""
    print(f"\n--- Patching: {path}")
    with open(path, "r", encoding="utf-8") as f:
        data = f.read()

    original_len = len(data)
    total = 0

    for label, pattern, replacement, is_regex in REPLACEMENTS:
        if is_regex:
            new_data, n = re.subn(pattern, replacement, data)
        else:
            n = data.count(pattern)
            new_data = data.replace(pattern, replacement)
        total += n
        data = new_data
        if n > 0:
            print(f"  {label:30s} : {n:5d} replacements")

    if total == 0:
        print("  (no replacements needed)")
        return 0

    print(f"  Total: {total} replacements, delta: {len(data) - original_len:+,} bytes")

    with open(path, "w", encoding="utf-8") as f:
        f.write(data)

    # Verification
    remaining_minio = len(re.findall(r'MinIO', data))
    remaining_min_io = len(re.findall(r'min\.io', data))
    if remaining_minio > 0:
        print(f"  WARNING: {remaining_minio} 'MinIO' references remain")
    if remaining_min_io > 0:
        print(f"  WARNING: {remaining_min_io} 'min.io' references remain")

    return total


def main():
    # Determine console root directory
    if len(sys.argv) > 1:
        root = sys.argv[1]
    else:
        root = os.getcwd()

    js_dir = os.path.join(root, JS_DIR)
    if not os.path.isdir(js_dir):
        print(f"ERROR: JS build directory not found: {js_dir}", file=sys.stderr)
        print(f"Run from the console repo root, or pass it as an argument.", file=sys.stderr)
        return 1

    # Glob for all .js files, skip .map files
    pattern = os.path.join(js_dir, "*.js")
    js_files = sorted(f for f in glob.glob(pattern) if not f.endswith(".map"))

    if not js_files:
        print(f"ERROR: No .js files found in {js_dir}", file=sys.stderr)
        return 1

    print(f"Console root: {root}")
    print(f"JS directory: {js_dir}")
    print(f"Found {len(js_files)} JS file(s):")
    for f in js_files:
        print(f"  {os.path.basename(f)}")

    grand_total = 0
    patched_count = 0
    for path in js_files:
        n = patch_file(path)
        grand_total += n
        if n > 0:
            patched_count += 1

    print(f"\n{'='*60}")
    print(f"Done. {grand_total} total replacements across {patched_count} file(s).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
