package helpers

import (
	"encoding/json"
	"net/http"
)

// RespondJSON translates an interface to json for response
func RespondJSON(w http.ResponseWriter, resp interface{}) {
	retJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(retJSON)
}

// BindJSON deserialize the body
func BindJSON(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&obj)
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