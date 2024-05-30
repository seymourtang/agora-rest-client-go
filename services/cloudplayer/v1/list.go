package v1

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type List struct {
	forwardedRegionPrefix core.ForwardedReginPrefix
	client                core.Client
	prefixPath            string
}

func (l *List) buildPath(param ListQueryParams) string {
	query := url.Values{}
	for k, v := range param {
		query.Add(k, v)
	}

	queryString := query.Encode()
	if queryString != "" {
		queryString = "?" + queryString
	}

	return string(l.forwardedRegionPrefix) + l.prefixPath + "/players" + queryString
}

type ListSuccessResp struct {
	TotalSize     int    `json:"totalSize"`
	Fields        string `json:"fields"`
	NextPageToken string `json:"nextPageToken"`
	Players       []struct {
		Name        string `json:"name"`
		StreamURL   string `json:"streamUrl"`
		ChannelName string `json:"channelName"`
		UID         int    `json:"uid"`
		ID          int    `json:"id"`
		CreateTs    int    `json:"createTs"`
		Status      string `json:"status"`
	} `json:"players"`
}

type ListResp struct {
	Response
	SuccessResponse ListSuccessResp
}

type ListQueryParams map[string]string

type ListQueryOpt interface {
	Apply(ListQueryParams)
}

type ListFilterOption string

func (f ListFilterOption) Apply(opt ListQueryParams) {
	filter := string(f)
	if filter != "" {
		opt["filter"] = fmt.Sprintf("channelName eq %s", filter)
	}
}

type ListPageSizeOption int

func (p ListPageSizeOption) Apply(opt ListQueryParams) {
	pageSize := int(p)
	if pageSize > 0 {
		opt["pageSize"] = strconv.Itoa(pageSize)
	}
}

type ListPageTokenOption string

func (t ListPageTokenOption) Apply(opt ListQueryParams) {
	pageToken := string(t)
	if pageToken != "" {
		opt["pageToken"] = pageToken
	}
}

func (l *List) Do(ctx context.Context, area core.ForwardedReginPrefix, queryOpts ...ListQueryOpt) (*ListResp, error) {
	param := make(map[string]string)

	for _, opt := range queryOpts {
		opt.Apply(param)
	}

	l.forwardedRegionPrefix = area
	path := l.buildPath(param)

	responseData, err := l.client.DoREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		var internalErr *core.InternalErr
		if !errors.As(err, &internalErr) {
			return nil, err
		}
	}

	resp := &ListResp{}
	resp.BaseResponse = responseData

	if responseData.HttpStatusCode == http.StatusOK {
		var successResponse ListSuccessResp
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
