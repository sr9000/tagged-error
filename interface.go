package taggerr

import "fmt"

// Tag is an interface that represents a tag that can be used to annotate errors with non-string types.
type Tag interface {
	comparable
	fmt.Stringer
}
