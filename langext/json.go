package langext

import (
	"bytes"
	"encoding/json"
	"github.com/joomcode/errorx"
)

type H map[string]any

type A []any

func TryPrettyPrintJson(str string) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return str
	}
	return prettyJSON.String()
}

func PrettyPrintJson(str string) (string, bool) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return str, false
	}
	return prettyJSON.String(), true
}

func PatchJson[JV string | []byte](rawjson JV, key string, value any) (JV, error) {
	var err error

	var jsonpayload map[string]any
	err = json.Unmarshal([]byte(rawjson), &jsonpayload)
	if err != nil {
		return *new(JV), errorx.Decorate(err, "failed to unmarshal payload")
	}

	jsonpayload[key] = value

	newjson, err := json.Marshal(jsonpayload)
	if err != nil {
		return *new(JV), errorx.Decorate(err, "failed to re-marshal payload")
	}

	return JV(newjson), nil
}

func PatchRemJson[JV string | []byte](rawjson JV, key string) (JV, error) {
	var err error

	var jsonpayload map[string]any
	err = json.Unmarshal([]byte(rawjson), &jsonpayload)
	if err != nil {
		return *new(JV), errorx.Decorate(err, "failed to unmarshal payload")
	}

	delete(jsonpayload, key)

	newjson, err := json.Marshal(jsonpayload)
	if err != nil {
		return *new(JV), errorx.Decorate(err, "failed to re-marshal payload")
	}

	return JV(newjson), nil
}
