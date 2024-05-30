package v1

import (
	"github.com/AgoraIO-Community/agora-rest-client-go/core"
)

type BaseCollection struct {
	prefixPath string
	client     core.Client
}

func NewCollection(prefixPath string, client core.Client) *BaseCollection {
	return &BaseCollection{
		prefixPath: "/v1" + prefixPath,
		client:     client,
	}
}

func (c *BaseCollection) Create() *Create {
	return &Create{
		forwardedRegionPrefix: core.DefaultForwardedReginPrefix,
		client:                c.client,
		prefixPath:            c.prefixPath,
	}
}

func (c *BaseCollection) Delete() *Delete {
	return &Delete{
		forwardedRegionPrefix: core.DefaultForwardedReginPrefix,
		client:                c.client,
		prefixPath:            c.prefixPath,
	}
}

func (c *BaseCollection) List() *List {
	return &List{
		forwardedRegionPrefix: core.DefaultForwardedReginPrefix,
		client:                c.client,
		prefixPath:            c.prefixPath,
	}
}

func (c *BaseCollection) Update() *Update {
	return &Update{
		forwardedRegionPrefix: core.DefaultForwardedReginPrefix,
		client:                c.client,
		prefixPath:            c.prefixPath,
	}
}
