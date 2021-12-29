package logger

type dummyLog struct{

}

func NewDummyLog() *dummyLog {
	return &dummyLog{}
}

func CreateDummyLog (configurationName string) (interface{}, error) {
	return NewDummyLog(), nil
}



func (this *dummyLog) Warn(msg string) {
}

func (this *dummyLog) Info(msg string) {
}

func (this *dummyLog) Error(msg string, err error) {
}

func (this *dummyLog) Errorf(format string, err error, args ...interface{}) {

}

