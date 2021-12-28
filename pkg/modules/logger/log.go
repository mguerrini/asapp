package logger

var log Logger

type Logger interface {
	Info(msg string)
	Error(msg string, err error)
	Errorf(format string, err error, args ...interface{})
}

func init () {
	log = NewConsoleLog()
}

func Info(msg string ){
	log.Info(msg)
}

func Error(msg string, err error ){
	log.Error(msg, err)
}

func Errorf(format string, err error, args ...interface{}){
	log.Errorf(format, err, args)
}