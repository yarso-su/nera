// errors.go
package nera

import "errors"

// Sentinel errors returned by Parse and the Document accessor methods.
// Use errors.Is to check for these across any wrapping.
var (
	ErrEmptyValueRow   = errors.New("entry has a key row but no value row")
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrTypeMismatch    = errors.New("entry type does not match requested accessor")
)
