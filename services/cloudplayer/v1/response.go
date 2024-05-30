package v1

import (
	"net/http"

	client "github.com/AgoraIO-Community/agora-rest-client-go/core"
)

const (
	ResourceIDHeaderKey = "X-Resource-Id"
	RequestIDHeaderKey  = "X-Request-Id"
)

type ErrResponse struct {
	Reason string `json:"reason"`
}

type Response struct {
	*client.BaseResponse
	ErrResponse ErrResponse
}

func (b *Response) IsSuccess() bool {
	if b.BaseResponse != nil {
		return b.HttpStatusCode == http.StatusOK
	} else {
		return false
	}
}

func (b *Response) GetResourceID() string {
	if b.BaseResponse != nil && b.Response != nil {
		return b.Header.Get(ResourceIDHeaderKey)
	}
	return ""
}

func (b *Response) GetRequestID() string {
	if b.BaseResponse != nil && b.Response != nil {
		return b.Header.Get(RequestIDHeaderKey)
	}
	return ""
}
