package v1

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type Create struct {
	forwardedRegionPrefix core.ForwardedReginPrefix
	client                core.Client
	prefixPath            string
}

func (c *Create) buildPath() string {
	return string(c.forwardedRegionPrefix) + c.prefixPath + "/players"
}

type CreateReqBody struct {
	Player *CreateReqPlayerPayload `json:"player"`
}

type AudioOptions struct {
	Volume  int `json:"volume,omitempty"`
	Profile int `json:"profile,omitempty"`
}

type VideoOptions struct {
	Width               int    `json:"width,omitempty"`
	Height              int    `json:"height,omitempty"`
	WidthHeightAdaption bool   `json:"widthHeightAdaption"`
	FrameRate           int    `json:"frameRate,omitempty"`
	BitRate             int    `json:"bitRate,omitempty"`
	CodeC               string `json:"codec,omitempty"`
	FillMode            string `json:"fillMode,omitempty"`
}

type DataStreamOptions struct {
	Enable bool `json:"enable,omitempty"`
}

type CreateReqPlayerPayload struct {
	VideoOptions      *VideoOptions      `json:"videoOptions,omitempty"`
	AudioOptions      *AudioOptions      `json:"audioOptions,omitempty"`
	DataStreamOptions *DataStreamOptions `json:"dataStreamOptions,omitempty"`
	EncryptKey        string             `json:"encryptKey,omitempty"`
	EncryptKdfSalt    string             `json:"encryptKdfSalt,omitempty"`
	PlayTs            int                `json:"playTs,omitempty"`
	RepeatTime        int                `json:"repeatTime,omitempty"`
	SeekPosition      int                `json:"seekPosition,omitempty"`
	StreamUrl         string             `json:"streamUrl"`
	ChannelName       string             `json:"channelName"`
	Token             string             `json:"token"`
	UID               int                `json:"uid"`
	Account           string             `json:"account,omitempty"`
	IdleTimeout       int                `json:"idleTimeout,omitempty"`
	Name              string             `json:"name"`
}

type CreateSuccessResp struct {
	Player struct {
		UID      int    `json:"uid"`
		ID       string `json:"id"`
		CreateTs int    `json:"createTs"`
		Status   string `json:"status"`
	} `json:"player"`
	Fields string `json:"fields"`
}

type CreateResp struct {
	Response
	SuccessResponse CreateSuccessResp
}

func (c *Create) Do(ctx context.Context, area core.ForwardedReginPrefix, payload *CreateReqPlayerPayload) (*CreateResp, error) {
	c.forwardedRegionPrefix = area
	path := c.buildPath()

	responseData, err := c.client.DoREST(ctx, path, http.MethodPost, &CreateReqBody{
		Player: payload,
	})
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	resp := &CreateResp{}
	resp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse CreateSuccessResp
		if err = responseData.UnmarshalToTarget(&successResponse); err != nil {
			return nil, &ServiceErr{
				RawResponse: responseData.Response,
				RequestID:   resp.GetRequestID(),
				ResourceID:  resp.GetResourceID(),
				Err:         err,
			}
		}

		resp.SuccessResponse = successResponse
	} else {
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
