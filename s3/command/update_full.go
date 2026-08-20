//go:build elastic && sqlite

package command

// set true if gtags are set
func init() {
	isFullVersion = true
}
