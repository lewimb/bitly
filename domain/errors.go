package domain

import "errors"

var ErrNotFound = errors.New("short url not found or expired")
