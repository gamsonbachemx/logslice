// Package tail provides a polling-based log tailer that periodically
// re-fetches lines from a remote source and forwards them to a channel.
//
// It is designed to work alongside the fetch package: callers supply a
// fetch function that returns newly available log lines, and Tailer
// handles the polling loop, back-off on errors, and context cancellation.
package tail
