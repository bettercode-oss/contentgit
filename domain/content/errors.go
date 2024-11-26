package content

import "github.com/pkg/errors"

var (
	ErrFieldNotFound        = errors.New("not found field")
	ErrFieldUpdateConflict  = errors.New("field update conflict")
	ErrContentAlreadyExists = errors.New("content with given id already exists")
	ErrUnknownEventType     = errors.New("unknown event type")
)
