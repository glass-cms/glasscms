package resource

// AlreadyExistsError represents an error when a resource already exists.
type AlreadyExistsError struct {
	Name     string
	Resource string
	err      error
}

func (e *AlreadyExistsError) Error() string {
	return e.err.Error()
}

// NewAlreadyExistsError creates a new ResourceAlreadyExistsError.
func NewAlreadyExistsError(name, resource string, err error) *AlreadyExistsError {
	return &AlreadyExistsError{
		Name:     name,
		Resource: resource,
		err:      err,
	}
}
