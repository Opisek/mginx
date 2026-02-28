package parsing

import (
	"encoding/json"
	"errors"
	"fmt"
	util "mginx/protocol/internal"
	"mginx/protocol/payloads"
)

func ParseStatusResponse(buffer []byte) (payloads.StatusResponse, error) {
	jsonString, buffer, err := util.ParseString(buffer)
	if err != nil {
		return payloads.StatusResponse{},
			errors.Join(errors.New("could not parse json string"), err)
	}

	if len(buffer) != 0 {
		return payloads.StatusResponse{},
			fmt.Errorf("remaining bytes at the end of payload: %v", len(buffer))
	}

	var payload payloads.StatusResponse

	err = json.Unmarshal([]byte(jsonString), &payload)
	if err != nil {
		return payloads.StatusResponse{},
			errors.Join(errors.New("could not unmarshal json object"), err)
	}

	return payload, nil
}
