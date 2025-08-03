package taggerr

import (
	"fmt"
	"slices"
	"strings"
)

type taggedError struct {
	tags map[any]int
	err  error
}

func (e taggedError) Error() string {
	if len(e.tags) == 0 {
		return e.err.Error()
	}

	slice := make([]any, 0, len(e.tags))
	for t := range e.tags {
		slice = append(slice, t)
	}

	slices.SortFunc(slice, func(a, b any) int {
		return e.tags[b] - e.tags[a] // latest added tag comes first
	})

	var sb strings.Builder
	sb.WriteByte('[')

	sb.WriteString(fmt.Sprint(slice[0]))

	for _, t := range slice[1:] {
		sb.WriteString(" | ")
		sb.WriteString(fmt.Sprint(t))
	}

	sb.WriteString("] @ ")
	sb.WriteString(e.err.Error())

	return sb.String()
}

func (e taggedError) Unwrap() error {
	return e.err
}
