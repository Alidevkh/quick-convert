package app

import (
	"errors"
	"path/filepath"

	"quickconvert/internal/engines"
)

// The two kinds of messages a converter can return instead of a result.

// UserError is a failure with a message meant for the user, shown in the
// notification as the reason, e.g. "This video has no audio".
type UserError struct{ Msg string }

func (e *UserError) Error() string { return e.Msg }

// Info is not a failure: the job decided not to create a file and wants to
// tell the user why, e.g. "Already optimized — no smaller file made".
type Info struct{ Msg string }

func (e *Info) Error() string { return e.Msg }

// ForFile adds a file's name to an error, for jobs with several files, so
// the user knows which one is the problem: "IMG_0042.heic: This file looks
// damaged or isn't a real HEIC file".
func ForFile(name string, err error) error {
	var userErr *UserError
	var engineErr *engines.EngineError
	switch {
	case errors.As(err, &userErr):
		return &UserError{Msg: name + ": " + userErr.Msg}
	case errors.As(err, &engineErr):
		return &UserError{Msg: name + ": " + reason(err, filepath.Ext(name))}
	default:
		return err // e.g. a missing engine: not about this file
	}
}
