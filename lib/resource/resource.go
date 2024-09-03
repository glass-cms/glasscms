package resource

// AlreadyExistsError represents an error when a resource already exists.
type AlreadyExistsError struct {
	UID      string
	Resource string
	err      error
}

func (e *AlreadyExistsError) Error() string {
	return e.err.Error()
}

// NewAlreadyExistsError creates a new ResourceAlreadyExistsError.
func NewAlreadyExistsError(uid, resource string, err error) *AlreadyExistsError {
	return &AlreadyExistsError{
		UID:      uid,
		Resource: resource,
		err:      err,
	}
}
