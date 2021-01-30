package etcdwatch

import (
	"context"
	"fmt"

	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"

	"github.com/trymanytimes/alg-manager/config"
	rpcpb "github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	"github.com/trymanytimes/alg-manager/pkg/proto/mvccpb"
)

type ETCDWatchClient struct {
	connect           *grpc.ClientConn
	WatchRequest      rpcpb.WatchRequest
	Watch_watchclient rpcpb.Watch_WatchClient
	WatchClient       rpcpb.WatchClient
}

func NewETCDWatchClient(conn *grpc.ClientConn, conf *config.ControllerConfig) (*ETCDWatchClient, error) {
	instance := &ETCDWatchClient{connect: conn, WatchClient: rpcpb.NewWatchClient(conn)}
	var err error
	instance.Watch_watchclient, err = instance.WatchClient.Watch(context.TODO(), grpc.EmptyCallOption{})
	if err != nil {
		return nil, fmt.Errorf("Watch error,%s", err.Error())
	}
	return instance, nil
}

type WatchHandler interface {
	DealEvent(event *mvccpb.Event)
}

func (c *ETCDWatchClient) Watch(key, end string, watchHandler WatchHandler) error {
	c.WatchRequest = rpcpb.WatchRequest{
		RequestUnion: &rpcpb.WatchRequest_CreateRequest{
			CreateRequest: &rpcpb.WatchCreateRequest{
				Key:      []byte(key),
				RangeEnd: []byte(end),
			},
		},
	}
	if err := c.Watch_watchclient.Send(&c.WatchRequest); err != nil {
		log.Fatalf("Send error,%s", err.Error())
	}
	return nil
}
