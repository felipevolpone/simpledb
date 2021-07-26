package simpledb

import "errors"

// Custom Errors
var (
	ErrDataMustBeSlicePointer = errors.New("provided data must be a pointer to a slice")
	ErrDataMustBeStructPointer = errors.New("provided data must be a pointer to a struct")
	ErrNotFound = errors.New("item not found")
)