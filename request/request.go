package request

import "encoding/json"

func MarshalRawRequest(rawRequest interface{}) string {
	if rawRequest == nil {
		return ""
	}

	body, _ := json.Marshal(rawRequest)

	return string(body)
}
