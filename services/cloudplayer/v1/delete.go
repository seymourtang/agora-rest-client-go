package v1

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Delete struct {
	forwardedRegionPrefix core.ForwardedReginPrefix
	client                core.Client
	prefixPath            string
}

func (d *Delete) buildPath(playerID string) string {
	return string(d.forwardedRegionPrefix) + d.prefixPath + "/players" + "/" + playerID
}

type DeleteResp struct {
	Response
}

func (d *Delete) Do(ctx context.Context, area core.ForwardedReginPrefix, playerID string) (*DeleteResp, error) {
	d.forwardedRegionPrefix = area
	path := d.buildPath(playerID)

	responseData, err := d.client.DoREST(ctx, path, http.MethodDelete, nil)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	resp := &DeleteResp{}
	resp.BaseResponse = responseData

	if responseData.HttpStatusCode != http.StatusOK {
		reasonResult := gjson.GetBytes(responseData.RawBody, "reason")
		if !reasonResult.Exists() {
			return resp, &ServiceErr{
				RawResponse: responseData.Response,
				RequestID:   resp.GetRequestID(),
				ResourceID:  resp.GetResourceID(),
				Err:         core.NewGatewayErr(responseData.HttpStatusCode, string(responseData.RawBody)),
			}
		}

		var errResponse ErrResponse
		if err = responseData.UnmarshalToTarget(&errResponse); err != nil {
			return resp, &ServiceErr{
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
