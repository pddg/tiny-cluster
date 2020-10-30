package errors

import "golang.org/x/xerrors"

const (
	// CodeErrOther is the error code for ErrOther.
	CodeErrOther = iota
	// CodeErrNotFound is the error code for ErrNotFound.
	CodeErrNotFound
	// CodeErrTimedOut is the error code for ErrTimedOut.
	CodeErrTimedOut
	// CodeErrAlreadyExists is the error code for ErrAlreadyExists.
	CodeErrAlreadyExists
	// CodeErrAuthFailed is the error code for ErrAuthFailed.
	CodeErrAuthFailed
	// CodeErrPermissionDenied is the error code for ErrErrPermissionDenied.
	CodeErrPermissionDenied
	// CodeErrContextCanceled is the error code for ErrContextCanceled.
	CodeErrContextCanceled
)

var (
	// ErrOther indicates that the error is unknown.
	ErrOther = newError(CodeErrOther, "something error")
	// ErrNotFound indicates that the item was not found.
	ErrNotFound = newError(CodeErrNotFound, "the item does not exist")
	// ErrTimedOut indicates that the communication is timed out.
	ErrTimedOut = newError(CodeErrTimedOut, "the operation timed out")
	// ErrAlreadyExists will be occured if the duplicated item is created.
	ErrAlreadyExists = newError(CodeErrAlreadyExists, "the item has already existed")

	// Authentication and Authorization
	// ErrAuthFailed indicates that the authentication was failed.
	// Username or password (or other authentication factor) is wrong.
	ErrAuthFailed = newError(CodeErrAuthFailed, "authentication failed")
	// ErrPermissionDenied indicates that the operation was denied because the user to exec it did not have the enough permission.
	ErrPermissionDenied = newError(CodeErrPermissionDenied, "permission denied")

	// ErrContextCanceled will be occured when the context was canceled.
	ErrContextCanceled = newError(CodeErrContextCanceled, "Context canceled")
)

func newError(code int, message string) error {
	return &TinyClusterError{
		code:    code,
		message: message,
	}
}

// TinyClusterError is a base error type for TinyCluster.
type TinyClusterError struct {
	code    int
	message string
}

func (e TinyClusterError) Is(err error) bool {
	var tmpErr *TinyClusterError
	return xerrors.As(err, tmpErr) && tmpErr.code == e.code
}

func (e TinyClusterError) Error() string {
	return e.message
}
