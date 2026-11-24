package app

// UserError is a failure with a message meant for the user, shown in the
// notification as the reason, e.g. "This video has no audio".
type UserError struct{ Msg string }

func (e *UserError) Error() string { return e.Msg }

// Info is not a failure: the job decided not to create a file and wants to
// tell the user why, e.g. "Already optimized — no smaller file made".
type Info struct{ Msg string }

func (e *Info) Error() string { return e.Msg }
