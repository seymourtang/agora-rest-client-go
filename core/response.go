package core

import (
	"encoding/json"
	"net/http"
)

type ResponseInterface interface {
	IsSuccess() bool
}

type BaseResponse struct {
	*http.Response
	RawBody        []byte
	HttpStatusCode int
}

// UnmarshalToTarget unmarshal body into target var
// successful if err is nil
func (r *BaseResponse) UnmarshalToTarget(target interface{}) error {
	err := json.Unmarshal(r.RawBody, target)
	if err != nil {
		return err
	}
	return nil
}
