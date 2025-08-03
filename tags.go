package taggerr

import "errors"

func wrap(err error) taggedError {
	if te, ok := err.(taggedError); ok {
		return te
	}

	return taggedError{
		tags: make(map[any]int),
		err:  err,
	}
}

func WithTag[T ETag](err error, tag T) error {
	if err == nil {
		return nil
	}

	te := wrap(err)

	if _, ok := te.tags[tag]; !ok {
		te.tags[tag] = len(te.tags)
	}

	return te
}

func WithTag2[T1, T2 ETag](err error, tag1 T1, tag2 T2) error {
	if err == nil {
		return nil
	}

	te := wrap(err)

	if _, ok := te.tags[tag1]; !ok {
		te.tags[tag1] = len(te.tags)
	}

	if _, ok := te.tags[tag2]; !ok {
		te.tags[tag2] = len(te.tags)
	}

	return te
}

func WithTag3[T1, T2, T3 ETag](err error, tag1 T1, tag2 T2, tag3 T3) error {
	if err == nil {
		return nil
	}

	te := wrap(err)

	if _, ok := te.tags[tag1]; !ok {
		te.tags[tag1] = len(te.tags)
	}

	if _, ok := te.tags[tag2]; !ok {
		te.tags[tag2] = len(te.tags)
	}

	if _, ok := te.tags[tag3]; !ok {
		te.tags[tag3] = len(te.tags)
	}

	return te
}

func WithTags[T ~string](err error, tags ...T) error {
	if err == nil {
		return nil
	}

	if len(tags) == 0 {
		return err
	}

	te := wrap(err)
	for _, t := range tags {
		if _, ok := te.tags[t]; !ok {
			te.tags[t] = len(te.tags)
		}
	}

	return te
}

// HasTag checks if the top level error has the specified tag.
func HasTag[T comparable](err error, tag T) bool {
	if err == nil {
		return false
	}

	te, ok := err.(taggedError)
	if !ok {
		return false
	}

	_, exists := te.tags[tag]
	return exists
}

// DeepHasTag checks if the error or any of its unwrapped errors has the specified tag.
func DeepHasTag[T comparable](err error, tag T) bool {
	var te taggedError

	for errors.As(err, &te) {
		if _, ok := te.tags[tag]; ok {
			return true
		}

		err = te.err
	}

	return false
}
