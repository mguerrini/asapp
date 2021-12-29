package logger

var logDefault Logger

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string, err error)
	Errorf(format string, err error, args ...interface{})
}

func init () {
	logDefault = NewConsoleLog()
}

func Info(msg string ){
	logDefault.Info(msg)
}

func Error(msg string, err error ){
	logDefault.Error(msg, err)
}

func Errorf(format string, err error, args ...interface{}){
	logDefault.Errorf(format, err, args)
}