package queue

import "errors"

var ErrEmptyQueue = errors.New("queue is empty")
var ErrOutOfCapacity = errors.New("queue is out of capacity")
