package storage

import "errors"

var ErrInitDatabase = errors.New("unable to initialize database")
var ErrCreatingTables = errors.New("unable to create tables")
