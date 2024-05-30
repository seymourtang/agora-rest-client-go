package v1

import (
	"fmt"
	"net/http"
)

type ServiceErr struct {
	RawResponse *http.Response
	RequestID   string
	ResourceID  string
	Err         error
}

func (e *ServiceErr) UnWrap() error {
	return e.Err
}

func (e *ServiceErr) Error() string {
	return fmt.Sprintf("requestID:%s,resourceID:%s,err:%s", e.RequestID, e.ResourceID, e.Err.Error())
}
