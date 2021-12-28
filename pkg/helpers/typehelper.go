package helpers

import "reflect"


func  GetTypeName(obj interface{}) string {
	if IsNil(obj) {
		return "nil"
	}
	tObj := reflect.TypeOf(obj)

	return tObj.String()
}


func IsNil(i interface{}) bool {
	return i == nil || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}

