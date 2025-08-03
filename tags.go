package taggerr

import "errors"

func wrap(err error) taggedError {
	if te, ok := err.(taggedError); ok { //nolint
		return te
	}

	return taggedError{
		tags: make(map[any]int),
		err:  err,
	}
}

// WithTag is a convenience function to add a single custom tag to an error.
func WithTag[T Tag](err error, tag T) error {
	if err == nil {
		return nil
	}

	te := wrap(err)

	if _, ok := te.tags[tag]; !ok {
		te.tags[tag] = len(te.tags)
	}

	return te
}

// WithTag2 is a convenience function to add two custom tags to an error.
func WithTag2[T1, T2 Tag](err error, tag1 T1, tag2 T2) error {
	if err == nil {
		return nil
	}

	terr := wrap(err)

	if _, ok := terr.tags[tag1]; !ok {
		terr.tags[tag1] = len(terr.tags)
	}

	if _, ok := terr.tags[tag2]; !ok {
		terr.tags[tag2] = len(terr.tags)
	}

	return terr
}

// WithTag3 is a convenience function to add three custom tags to an error.
func WithTag3[T1, T2, T3 Tag](err error, tag1 T1, tag2 T2, tag3 T3) error {
	if err == nil {
		return nil
	}

	terr := wrap(err)

	if _, ok := terr.tags[tag1]; !ok {
		terr.tags[tag1] = len(terr.tags)
	}

	if _, ok := terr.tags[tag2]; !ok {
		terr.tags[tag2] = len(terr.tags)
	}

	if _, ok := terr.tags[tag3]; !ok {
		terr.tags[tag3] = len(terr.tags)
	}

	return terr
}

// WithTags is a convenience function to add multiple string-like tags to an error.
func WithTags[T ~string](err error, tags ...T) error {
	if err == nil {
		return nil
	}

	if len(tags) == 0 {
		return err
	}

	terr := wrap(err)
	for _, t := range tags {
		if _, ok := terr.tags[t]; !ok {
			terr.tags[t] = len(terr.tags)
		}
	}

	return terr
}

// HasTag checks if the top level error has the specified tag.
func HasTag[T comparable](err error, tag T) bool {
	if err == nil {
		return false
	}

	terr, ok := err.(taggedError) //nolint
	if !ok {
		return false
	}

	_, exists := terr.tags[tag]

	return exists
}

// DeepHasTag checks if the error or any of its unwrapped errors has the specified tag.
func DeepHasTag[T comparable](err error, tag T) bool {
	var terr taggedError

	for errors.As(err, &terr) {
		if _, ok := terr.tags[tag]; ok {
			return true
		}

		err = terr.err
	}

	return false
}
