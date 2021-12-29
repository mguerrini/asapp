package logger

import (
	"fmt"
	"os"
)

type consoleLog struct{

}

func NewConsoleLog() *consoleLog {
	return &consoleLog{}
}

func CreateConsoleLog (configurationName string) (interface{}, error) {
	return NewConsoleLog(), nil
}

func (this *consoleLog) Warn(msg string) {
	fmt.Fprint(os.Stdout, "[WARN] " + msg + "\n")
}

func (this *consoleLog) Info(msg string) {
	fmt.Fprint(os.Stdout, "[INFO] " + msg + "\n")
}

func (this *consoleLog) Error(msg string, err error) {
	if err != nil {
		fmt.Fprint(os.Stdout, "[ERROR] " + msg + " - "+err.Error()+"\n")
	} else {
		fmt.Fprint(os.Stdout, "[ERROR] " + msg + "\n")
	}
}

func (this *consoleLog) Errorf(format string, err error, args ...interface{}) {
	this.Error(fmt.Sprintf(format, args...), err)
}


