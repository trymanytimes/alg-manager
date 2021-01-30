package firewalld

import (
	"context"
	"fmt"

	"github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	"github.com/trymanytimes/alg-manager/pkg/proto/mvccpb"
	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"
)

type FirewalldHandler struct {
	KVClient etcdserverpb.KVClient
}

func NewFirewalldHandler(kvClient etcdserverpb.KVClient) *FirewalldHandler {
	return &FirewalldHandler{KVClient: kvClient}
}

func (h *FirewalldHandler) DealEvent(event *mvccpb.Event) {
	rsp, err := h.KVClient.Range(context.Background(), &etcdserverpb.RangeRequest{Key: []byte("/"), RangeEnd: []byte("0")}, grpc.EmptyCallOption{})
	if err != nil {
		log.Errorf("Range error:%s", err.Error())
	}
	for _, v := range rsp.Kvs {
		fmt.Println("key:", string(v.Key))
		fmt.Println("value:", string(v.Value))
	}
}
