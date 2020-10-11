package usecase

import "golang.org/x/xerrors"

var (
	ErrContextCanceled = xerrors.New("Context canceled")
)
