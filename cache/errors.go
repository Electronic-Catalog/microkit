package cache

import "errors"

type Error error

var (
	NotFoundError Error = errors.New("key not-found")
)
