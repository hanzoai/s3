// Drop-in module for github.com/tyler-smith/go-bip39, whose upstream GitHub
// repository was deleted (fetch fails under GOPROXY=direct). It re-exports the
// luxfi fork (github.com/luxfi/go-bip39, API-compatible) under the original
// module path so the lone transitive consumer (rclone's internxt backend, via
// github.com/internxt/rclone-adapter) keeps building without depending on the
// dead upstream. A plain `replace tyler-smith/go-bip39 => luxfi/go-bip39` is not
// possible: the luxfi fork renamed its module path and self-imports
// github.com/luxfi/go-bip39/wordlists, so one module version would have to serve
// two module paths ("used for two different module paths"). This thin shim is a
// distinct module under the original path, so `go mod tidy` resolves cleanly.
module github.com/tyler-smith/go-bip39

go 1.25

require github.com/luxfi/go-bip39 v1.1.2
