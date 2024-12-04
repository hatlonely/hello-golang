package kvstore

import "github.com/pkg/errors"

var ErrNotFound = errors.New("not found")

func IsNotFound(err error) bool {
	return errors.Cause(err) == ErrNotFound
}
