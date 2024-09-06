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

func NewAlreadyExistsError(name, resource string, err error) *AlreadyExistsError {
	return &AlreadyExistsError{
		Name:     name,
		Resource: resource,
		err:      err,
	}
}

// NotFoundError represents an error when a resource is not found.
type NotFoundError struct {
	Name     string
	Resource string
	err      error
}

func (e *NotFoundError) Error() string {
	return e.err.Error()
}

func NewNotFoundError(name, resource string, err error) *NotFoundError {
	return &NotFoundError{
		Name:     name,
		Resource: resource,
		err:      err,
	}
}
