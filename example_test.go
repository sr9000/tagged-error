package taggerr_test

import (
	"errors"
	"fmt"
	taggerr "github.com/sr9000/tagged-error"
	"github.com/stretchr/testify/assert"
	"math"
	"slices"
	"strings"
	"testing"
)

func ExamplePrint() {
	// This example demonstrates how to use the Print function from the taggerr package.
	// It prints a tagged error with a custom tag and an underlying error message.

	// Got an error
	err := errors.New("foobar error")

	// Add some tags to the error
	err = taggerr.WithTags(err, "buzz tag", "another tag")

	// Print the tagged error
	// Remember, last added tag comes first
	fmt.Println(err)

	// Output:
	// [another tag | buzz tag] @ foobar error
}

func ExampleWithTags() {
	// This example demonstrates how to use the Enum (string) as a tag in the taggerr package.
	// It creates a tagged error with a custom tag and prints it.

	type Enum string

	const (
		Tag1 Enum = "tag1"
		Tag2 Enum = "tag2"
	)

	// Got an error
	err := errors.New("foobar error")

	// Add some tags to the error
	err = taggerr.WithTags(err, Tag1, Tag2)

	// Print the tagged error
	// Remember, last added tag comes first
	fmt.Println(err)

	// Output:
	// [tag2 | tag1] @ foobar error
}

type specialTag struct {
	a, b int
}

func (t specialTag) String() string {
	return fmt.Sprintf("special-tag(a=%d, b=%d)", t.a, t.b)
}

// This test checks if the error has the expected tags.
// It creates an error, adds tags to it (string and specialTag)
// And then verifies that the tags are present (also checks if unpresent tags are missing).
func TestHasTags(t *testing.T) {
	t.Parallel()

	err := errors.New("test error")
	err = taggerr.WithTags(err, "tag1", "tag2")
	err = taggerr.WithTag(err, specialTag{a: 1, b: 2})

	assert.True(t, taggerr.HasTag(err, "tag1"))
	assert.True(t, taggerr.HasTag(err, "tag2"))
	assert.True(t, taggerr.HasTag(err, specialTag{a: 1, b: 2}))

	assert.False(t, taggerr.HasTag(err, "tag3"))
	assert.False(t, taggerr.HasTag(err, specialTag{a: 3, b: 4}))
	assert.False(t, taggerr.HasTag(err, 12345)) // Check for a non-existent tag type

	s := err.Error()
	assert.Equal(t, "[special-tag(a=1, b=2) | tag2 | tag1] @ test error", s)
}

func TestWrapTaggedError(t *testing.T) {
	t.Parallel()

	// Create an error with tags
	err := errors.New("original error")
	err = taggerr.WithTags(err, "tag1", "tag2")

	// Wrap the error
	werr := fmt.Errorf("wrapping error: %w", err)

	// Check if the wrapped error has no tags
	assert.False(t, taggerr.HasTag(werr, "tag1"))
	assert.False(t, taggerr.HasTag(werr, "tag2"))
	assert.False(t, taggerr.HasTag(werr, "wubba-lubba-dub-dub")) // Check for a non-existent tag

	// Check if the original error can be unwrapped
	terr := taggerr.WithTags(werr, "foo", "bar")

	// New tagged error has its tags
	assert.True(t, taggerr.HasTag(terr, "foo"))
	assert.True(t, taggerr.HasTag(terr, "bar"))

	assert.False(t, taggerr.HasTag(terr, "wubba-lubba-dub-dub")) // Check for a non-existent tag

	// But still missing tags from the initially tagged error
	assert.False(t, taggerr.HasTag(terr, "tag1"))
	assert.False(t, taggerr.HasTag(terr, "tag2"))

	// DeepHasTag founds all the tags
	assert.True(t, taggerr.DeepHasTag(terr, "foo"))
	assert.True(t, taggerr.DeepHasTag(terr, "bar"))
	assert.True(t, taggerr.DeepHasTag(terr, "tag1"))
	assert.True(t, taggerr.DeepHasTag(terr, "tag2"))

	assert.False(t, taggerr.DeepHasTag(terr, "wubba-lubba-dub-dub")) // Check for a non-existent tag

	s := terr.Error()
	assert.Equal(t, "[bar | foo] @ wrapping error: [tag2 | tag1] @ original error", s)
}

// TestRepeatedTagsDoNotChangeTheOrder checks that adding the same tag multiple times does not change the order of tags.
func TestRepeatedTagsDoNotChangeTheOrder(t *testing.T) {
	t.Parallel()

	slice := make([]string, 100)
	for i := range slice {
		slice[i] = fmt.Sprint("tag", i+1)
	}

	// Create an error with tags
	err := errors.New("original error")
	for i := range slice {
		err = taggerr.WithTags(err, slice[i])
	}

	// Add the same tags again
	for i := range 1_000_000 {
		j := int(math.Sqrt(2.0*float64(i)*float64(i))) % len(slice)
		err = taggerr.WithTags(err, slice[j])
	}

	// Check if the tags are still in the same order
	slices.Reverse(slice)
	expected := fmt.Sprintf("[%s] @ original error", strings.Join(slice, " | "))
	got := err.Error()

	assert.Equal(t, expected, got)
}
