package taggerr

import "fmt"

// ETag is an interface that represents a tag that can be used to annotate errors with non-string types.
type ETag interface {
	comparable
	fmt.Stringer
}
