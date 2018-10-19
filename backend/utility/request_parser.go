package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ReadRequestData : reading and unmarshalling request body data
func ReadRequestData(req *http.Request, targetObject interface{}) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("error while reading request body: %s", err)
	}

	err = json.Unmarshal(body, &targetObject)
	if err != nil {
		return fmt.Errorf("error while parsing request body[%s]: %s", string(body), err)
	}
	return nil
}

// ToJSON marshals interface into json []byte
func ToJSON(results interface{}) []byte {
	b, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	return b
}
