// Package checkpoint provides persistent stream-position tracking for logslice.
//
// A Store records the last successfully processed byte offset and timestamp
// for each remote URL. State is written atomically to a JSON file so that
// logslice can resume streaming from the correct position after a crash or
// planned restart, avoiding both duplicate and missed log lines.
//
// Usage:
//
//	store, err := checkpoint.New("/var/lib/logslice/checkpoint.json")
//	if err != nil { /* handle */ }
//
//	if st, ok := store.Get(url); ok {
//		// resume from st.Offset
//	}
//
//	store.Set(url, newOffset, lastLineTime)
package checkpoint
