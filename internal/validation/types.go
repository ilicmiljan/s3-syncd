package validation

type TaskFieldError struct {
	TaskID  string
	Field   string
	Message string
}

type BucketFieldError struct {
	BucketID string
	Field    string
	Message  string
}

type TaskValidationError struct {
	Errors []TaskFieldError
}

type BucketValidationError struct {
	Errors []BucketFieldError
}
