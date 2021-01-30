package updatedip

import (
	"context"
	"fmt"

	"github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	"github.com/trymanytimes/alg-manager/pkg/proto/mvccpb"
	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"
)

type UpdatedIPHandler struct {
	KVClient etcdserverpb.KVClient
}

func NewUpdatedIPHandler(kvClient etcdserverpb.KVClient) *UpdatedIPHandler {
	return &UpdatedIPHandler{KVClient: kvClient}
}

func (h *UpdatedIPHandler) DealEvent(event *mvccpb.Event) {
	rsp, err := h.KVClient.Range(context.Background(), &etcdserverpb.RangeRequest{Key: []byte("/"), RangeEnd: []byte("0")}, grpc.EmptyCallOption{})
	if err != nil {
		log.Errorf("Range error:%s", err.Error())
	}
	for _, v := range rsp.Kvs {
		fmt.Println("key:", string(v.Key))
		fmt.Println("value:", string(v.Value))
	}
}
