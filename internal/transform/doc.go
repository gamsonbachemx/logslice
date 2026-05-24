// Package transform provides field projection and renaming for structured
// JSON log entries.
//
// A FieldSelector is constructed from a comma-separated specification such as
//
//	"level,msg,ts"
//
// Fields may optionally be aliased:
//
//	"level:severity,msg:message"
//
// When applied to a parsed log entry, only the selected fields are retained
// and any aliases are applied. Entries with no selector pass through unchanged.
package transform
