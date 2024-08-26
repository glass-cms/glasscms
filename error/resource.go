package item

// ResourceAlreadyExistsError represents an error when a resource already exists.
type ResourceAlreadyExistsError struct {
	UID      string
	Resource string
	err      error
}

func (e *ResourceAlreadyExistsError) Error() string {
	return e.err.Error()
}

// NewResourceAlreadyExistsError creates a new ResourceAlreadyExistsError.
func NewResourceAlreadyExistsError(uid, resource string, err error) *ResourceAlreadyExistsError {
	return &ResourceAlreadyExistsError{
		UID:      uid,
		Resource: resource,
		err:      err,
	}
}
