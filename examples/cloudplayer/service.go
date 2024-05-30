package main

import (
	"context"
	"log"
	"time"

	"github.com/AgoraIO-Community/agora-rest-client-go/core"
	"github.com/AgoraIO-Community/agora-rest-client-go/services/cloudplayer"
	v1 "github.com/AgoraIO-Community/agora-rest-client-go/services/cloudplayer/v1"
)

type Service struct {
	region     core.RegionArea
	appId      string
	cname      string
	uid        string
	credential core.Credential
}

func NewService(region core.RegionArea, appId, cname, uid string) *Service {
	return &Service{
		region:     region,
		appId:      appId,
		cname:      cname,
		uid:        uid,
		credential: nil,
	}
}

func (s *Service) SetCredential(username, password string) {
	s.credential = core.NewBasicAuthCredential(username, password)
}

func (s *Service) Run(token string, streamUrl string) {
	ctx := context.Background()
	c := core.NewClient(&core.Config{
		AppID:      s.appId,
		Credential: s.credential,
		RegionCode: s.region,
		Logger:     core.NewDefaultLogger(core.LogDebug),
	})

	implV1 := cloudplayer.NewAPI(c).V1()

	// create
	createResp, err := implV1.Create().Do(ctx, core.CNForwardedReginPrefix, &v1.CreateReqPlayerPayload{
		StreamUrl:   streamUrl,
		ChannelName: s.cname,
		Token:       token,
		UID:         0,
		Account:     "",
	})
	if err != nil {
		log.Println(err)
		return
	}

	if createResp.IsSuccess() {
		log.Printf("create cloud player success:%+v", createResp)
	} else {
		log.Printf("create cloud player failed:%+v", createResp)
		return
	}

	resourceID := createResp.GetResourceID()
	if resourceID == "" {
		log.Printf("create cloud player failed, resourceID is empty")
		return
	}

	defer func() {
		// stop
	}()

	// list
	for i := 0; i < 3; i++ {
		listResp, err := implV1.List().Do(ctx, core.CNForwardedReginPrefix, v1.ListPageSizeOption(10))
		if err != nil {
			log.Println(err)
			return
		}
		if listResp.IsSuccess() {
			log.Printf("list cloud player success:%+v", listResp)
		} else {
			log.Printf("list cloud player failed:%+v", listResp)
			return
		}

		time.Sleep(time.Second * 10)
	}

	// update
	updateResp, err := implV1.Update().Do(ctx, core.CNForwardedReginPrefix, resourceID, &v1.UpdateReqPlayerPayload{})
	if err != nil {
		log.Println(err)
		return
	}

	if updateResp.IsSuccess() {
		log.Printf("update cloud player success:%+v", updateResp)
	} else {
		log.Printf("update cloud player failed:%+v", updateResp)
		return
	}
}
