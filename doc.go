// Package nera implements a positional, line-oriented micro-language.
//
// A file is a sequence of blank-line-delimited blocks. Each block is a
// mini-CSV: the first line is one or more comma-separated keys, and each
// following line is a comma-separated row of values matching the key
// count. Multiple value rows promote a block from a single entry into a
// collection. Values containing a literal comma must be wrapped in
// double quotes; the quotes are stripped during parsing. Bare commas
// outside quoted values are always treated as field separators.
//
// Document.Entries preserves file order; position is the API, not key
// names.
package nera
