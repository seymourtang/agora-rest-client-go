package individualrecording

import (
	"context"

	baseV1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudrecording/v1"
)

type Acquire struct {
	Base *baseV1.Acquire
}

var _ baseV1.AcquireIndividualRecording = (*Acquire)(nil)

func (a *Acquire) Do(ctx context.Context, cname string, uid string, enablePostponeTranscodingMix bool, clientRequest *baseV1.AcquireIndividualRecodingClientRequest) (*baseV1.AcquireResp, error) {
	var startParameter *baseV1.StartClientRequest
	if clientRequest.StartParameter != nil {
		startParameter = &baseV1.StartClientRequest{
			Token:               clientRequest.StartParameter.Token,
			StorageConfig:       clientRequest.StartParameter.StorageConfig,
			RecordingConfig:     clientRequest.StartParameter.RecordingConfig,
			RecordingFileConfig: clientRequest.StartParameter.RecordingFileConfig,
			SnapshotConfig:      clientRequest.StartParameter.SnapshotConfig,
			AppsCollection:      clientRequest.StartParameter.AppsCollection,
			TranscodeOptions:    clientRequest.StartParameter.TranscodeOptions,
		}
	}

	scene := 0
	if enablePostponeTranscodingMix {
		scene = 2
	}
	return a.Base.Do(ctx, &baseV1.AcquireReqBody{
		Cname: cname,
		Uid:   uid,
		ClientRequest: &baseV1.AcquireClientRequest{
			Scene:               scene,
			ResourceExpiredHour: clientRequest.ResourceExpiredHour,
			ExcludeResourceIds:  clientRequest.ExcludeResourceIds,
			RegionAffinity:      clientRequest.RegionAffinity,
			StartParameter:      startParameter,
		},
	})
}
