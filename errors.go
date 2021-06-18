package simpledb

import "errors"

// Custom Errors
var (
	ErrDataMustBeStructPointer = errors.New("provided data must be a pointer to a struct")
	ErrNotFound = errors.New("item not found")
)