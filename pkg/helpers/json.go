package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type JsonHelper struct {

}

func NewJsonHelper() *JsonHelper {
	return &JsonHelper{}
}

func  (this *JsonHelper) GetStringFromJson(target string, path string)  (*string, error) {
	aPath := this.getPath(path)

	val, err := this.getValueFromJson(target, aPath)

	if err != nil {
		return nil, err
	}

	if IsNil(val) {
		return nil, nil
	}

	output := fmt.Sprintf("%v", val)
	return &output, nil
}

func  (this *JsonHelper) GetIntFromJson(target string, path string)  (*int, error) {
	aPath := this.getPath(path)

	val, err := this.getValueFromJson(target, aPath)

	return this.convertToInt(val, err)
}


func (this *JsonHelper) getPath(path string)  []string {
	path = strings.ReplaceAll(path, "[", ".")
	path = strings.ReplaceAll(path, "]", ".")

	pathArray := strings.Split(path, ".")

	output := make([]string, 0)
	//elimino todos lo indices vacios
	for _, p := range pathArray {
		if p != "" {
			output = append(output, p)
		}
	}
	return output
}

func  (this *JsonHelper) getValueFromJson(targetJson string, path []string)  (interface{}, error) {
	if targetJson == "" {
		return nil, nil
	}

	if targetJson[0] == '{' {
		var targetMap map[string]interface{}

		err2 := json.Unmarshal([]byte(targetJson), &targetMap);
		if err2 != nil {
			return nil, err2
		}


		return this.getValueFromInterface(targetMap, path)
	} else {
		var targetMap []interface{}

		err2 := json.Unmarshal([]byte(targetJson), &targetMap);
		if err2 != nil {
			return nil, err2
		}

		return this.getValueFromInterface(targetMap, path)
	}
}

func  (this *JsonHelper) getValue(target interface{}, path []string)  (interface{}, error) {
	if IsNil(target) {
		return nil, nil
	}
	var targetMap map[string]interface{}

	targetJson, err1 := json.Marshal(target)

	if err1 != nil {
		return nil, err1
	}

	err2 := json.Unmarshal([]byte(targetJson), &targetMap);
	if err2 != nil {
		return nil, err2
	}

	return this.getValueFromInterface(targetMap, path)
}

func  (this *JsonHelper) getValueFromInterface(from interface{}, path []string)  (interface{}, error) {
	if IsNil(from) {
		return nil, nil
	}

	if path == nil || len(path) == 0 {
		return from, nil
	}

	val := reflect.ValueOf(from)

	switch val.Type().Kind() {
	case reflect.Array:
		return this.getValueFromArray(val, path)
	case reflect.Slice:
		return this.getValueFromArray(val, path)
	case reflect.Map:
		valMap, _ := from.(map[string]interface{})
		return this.getValueFromMap(valMap, path)
	default:
		return nil, errors.New("Invalid path")
	}
}

func  (this *JsonHelper) getValueFromMap(from map[string]interface{}, path []string)  (interface{}, error) {
	val, found := from[path[0]]

	if !found {
		return nil, nil
	}

	return this.getValueFromInterface(val, path[1:])
}

func  (this *JsonHelper) getValueFromArray(fromArray reflect.Value, path []string)  (interface{}, error) {
	var indexStr string = path[0]

	if indexStr[0] == '[' {
		indexStr = indexStr[1:]
	}

	if len(indexStr) == 0 || (len(indexStr) == 1 &&  indexStr[len(indexStr)-1] == ']') {
		return nil, errors.New("Invalid path")
	}

	if  indexStr[len(indexStr)-1] == ']' {
		indexStr = indexStr[:len(indexStr)-1]
	}

	index, err := strconv.Atoi(indexStr)

	if err != nil {
		return nil, err
	}

	if 0 <= index && index < fromArray.Len(){
		item := fromArray.Index(index).Interface()
		return this.getValueFromInterface(item, path[1:])
	} else {
		return nil, nil
	}
}

func  (this *JsonHelper) convertToInt(val interface{}, err error)  (*int, error) {
	if err != nil {
		return nil, err
	}

	if IsNil(val){
		return nil, err
	}

	reflectVal := reflect.ValueOf(val)

	switch reflectVal.Type().Kind() {
	case reflect.String:
		strVal := val.(string)
		iVal, err := strconv.Atoi(strVal)

		if err != nil {
			return nil, err
		}

		return &iVal, nil
	case reflect.Float64:
		f64, ok  := val.(float64)
		if ok {
			pInt := int(f64)
			return &pInt, nil
		}

		return nil, errors.New("Invalid return type")
	case reflect.Float32:
		t32, ok  := val.(float32)
		if ok {
			pInt := int(t32)
			return &pInt, nil
		}

		return nil, errors.New("Invalid return type")
	default:
		//trata de convertirlo
		t, ok  := val.(int)
		if ok {
			return &t, nil
		}
		t2, ok2  := val.(*int)
		if ok2 {
			return t2, nil
		}

		return nil, errors.New("Invalid return type")
	}
}


// Could be something more intelligent with reflection
func Copy (source interface{}, destiny interface{} ) {
	if IsNil(source) || IsNil(destiny) {
		return
	}

	sourceJson, _ := json.Marshal(source)
	str := string(sourceJson)
	if len(str) == 0{

	}
	json.Unmarshal(sourceJson, destiny)
}