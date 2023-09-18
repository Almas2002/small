package httpErrors

import "github.com/pkg/errors"

var AlreadyExists = errors.New("already exists")

var InvalidArguments = errors.New("invalid arguments")

var NotFound = errors.New("not found")
