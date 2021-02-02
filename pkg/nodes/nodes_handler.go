package nodes

import (
	"context"
	"os/exec"
	"strings"

	"github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	"github.com/trymanytimes/alg-manager/pkg/proto/mvccpb"
	"github.com/trymanytimes/alg-manager/pkg/util"
	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"
)

type NodesHandler struct {
	KVClient etcdserverpb.KVClient
}

func NewNodesHandler(kvClient etcdserverpb.KVClient) *NodesHandler {
	instance := &NodesHandler{KVClient: kvClient}
	return instance
}

func (h *NodesHandler) DealEvent(event *mvccpb.Event) {
	log.Infof("NodesHandler start")
	//get the node in cluster:/project_nm/version1_0_0/nodes/{nodeid}
	//get the node ip: /project_nm/version1_0_0/host/{hostid}/hostinfo/ipv6_mgr_addr
	rsp, err := h.KVClient.Range(context.Background(), &etcdserverpb.RangeRequest{
		Key:      []byte("/project_nm/version1_0_0/host/" + string(event.Kv.Key)[strings.LastIndex(string(event.Kv.Key), "/")+1:] + "/hostinfo/ipv6_mgr_addr"),
		RangeEnd: []byte(""),
	}, grpc.EmptyCallOption{})
	if err != nil {
		log.Errorf("Range error:%s", err.Error())
		return
	}
	if len(rsp.Kvs) != 1 {
		log.Errorf("no match host for id:%s exists", string(event.Kv.Key)[strings.LastIndex(string(event.Kv.Key), "/")+1:])
		return
	}
	allIPPorts, err := util.GetIPPorts(h.KVClient)
	if err != nil {
		log.Errorf("GetIPPorts error:%s", err.Error())
		return
	}
	if event.GetType().String() == "PUT" {
		for ipPort := range allIPPorts {
			command := "ipvsadm -a -t " + ipPort + " -r " + string(rsp.Kvs[0].Value) + " -g "
			cmd := exec.Command("/bin/bash", "-c", command)
			if _, err := cmd.Output(); err != nil {
				log.Errorf("exec %s error: %s", command, err.Error())
			}
		}
	} else if event.GetType().String() == "DELETE" {
		for ipPort := range allIPPorts {
			command := "ipvsadm -d -t " + ipPort + " -r " + string(rsp.Kvs[0].Value)
			cmd := exec.Command("/bin/bash", "-c", command)
			if _, err := cmd.Output(); err != nil {
				log.Errorf("exec %s error: %s", command, err.Error())
			}
		}
	}
}
