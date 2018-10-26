package log

func ErrCh(l *Log) chan error {
	errCh := make(chan error)

	go func() {
		for v := range errCh {
			l.Errorf(v.Error())
		}
	}()

	return errCh
}

func WarnCh(l *Log) chan error {
	errCh := make(chan error)

	go func() {
		for v := range errCh {
			l.Warnf(v.Error())
		}
	}()

	return errCh
}
