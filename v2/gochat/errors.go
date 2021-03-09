package gochat

import "errors"

var (
	ENOENT    = errors.New("record not found.")
	ECONFLICT = errors.New("entry conflicts with a saved record.")
)
