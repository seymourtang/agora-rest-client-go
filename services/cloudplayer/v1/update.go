package v1

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Update struct {
	forwardedRegionPrefix core.ForwardedReginPrefix
	client                core.Client
	prefixPath            string
}

func (u *Update) buildPath(playerID string) string {
	return string(u.forwardedRegionPrefix) + u.prefixPath + "/players" + "/" + playerID
}

type UpdateReqBody struct {
	Player *UpdateReqPlayerPayload `json:"player"`
}

type UpdateReqPlayerPayload struct {
}

type UpdateResp struct {
	Response
}

func (u *Update) Do(ctx context.Context, area core.ForwardedReginPrefix, playerID string, payload *UpdateReqPlayerPayload) (*UpdateResp, error) {
	u.forwardedRegionPrefix = area
	path := u.buildPath(playerID)

	responseData, err := u.client.DoREST(ctx, path, http.MethodPatch, &UpdateReqBody{
		Player: payload,
	})
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	resp := &UpdateResp{}
	resp.BaseResponse = responseData

	if responseData.HttpStatusCode != http.StatusOK {
		reasonResult := gjson.GetBytes(responseData.RawBody, "reason")
		if !reasonResult.Exists() {
			return nil, &ServiceErr{
				RawResponse: responseData.Response,
				RequestID:   resp.GetRequestID(),
				ResourceID:  resp.GetResourceID(),
				Err:         core.NewGatewayErr(responseData.HttpStatusCode, string(responseData.RawBody)),
			}
		}

		var errResponse ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return nil, &ServiceErr{
				RawResponse: responseData.Response,
				RequestID:   resp.GetRequestID(),
				ResourceID:  resp.GetResourceID(),
				Err:         err,
			}
		}
		resp.ErrResponse = errResponse
	}

	return resp, nil
}
