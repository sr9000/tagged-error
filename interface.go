package taggerr

import "fmt"

type ETag interface {
	comparable
	fmt.Stringer
}
