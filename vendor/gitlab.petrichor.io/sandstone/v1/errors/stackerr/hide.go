package stackerr

type hide struct {
	error
}

func Hide(err error) error {
	return &hide{
		error: err,
	}
}
