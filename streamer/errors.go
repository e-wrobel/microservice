package streamer

import "errors"

var ErrMarshaling = errors.New("unable to marshal data")
var ErrUnmarshalling = errors.New("unable to unmarshal data")
var ErrNotPrepared = errors.New("unable to prepare request")
var ErrNotSent = errors.New("unable to send request")
var ErrBodyCorrupted = errors.New("unable to ready response body")
