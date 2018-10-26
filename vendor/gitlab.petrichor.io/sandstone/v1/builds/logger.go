package builds

type logger interface {
	Normalf(msg string, args ...interface{}) (err error)
	Noticef(msg string, args ...interface{}) (err error)
	Successf(msg string, args ...interface{}) (err error)
	Debugf(msg string, args ...interface{}) (err error)
	Warnf(msg string, args ...interface{}) (err error)
	Errorf(msg string, args ...interface{}) (err error)
	Write(in []byte) (written int, err error)
}
