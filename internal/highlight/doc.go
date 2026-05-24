// Package highlight provides ANSI terminal color highlighting for logslice output.
//
// It wraps regex-matched substrings in ANSI escape sequences so that
// patterns of interest stand out when streaming logs to a terminal.
// When output is redirected or color is disabled the Highlighter acts as
// a transparent pass-through, leaving strings unmodified.
package highlight
