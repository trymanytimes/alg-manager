package updatedip

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/trymanytimes/alg-manager/config"
	"github.com/trymanytimes/alg-manager/pkg/proto/etcdserverpb"
	"github.com/trymanytimes/alg-manager/pkg/proto/mvccpb"
	"github.com/trymanytimes/alg-manager/pkg/util"
	"github.com/zdnscloud/cement/log"
	"google.golang.org/grpc"
)

type UpdatedIPHandler struct {
	KVClient etcdserverpb.KVClient
}

func NewUpdatedIPHandler(kvClient etcdserverpb.KVClient, conf *config.ControllerConfig) (*UpdatedIPHandler, error) {
	instance := &UpdatedIPHandler{KVClient: kvClient}
	rsp, err := instance.KVClient.Range(context.Background(), &etcdserverpb.RangeRequest{
		Key:      []byte("/project_nm/version1_0_0/cluster/001/config_vips/used_VIP6/"),
		RangeEnd: []byte("/project_nm/version1_0_0/cluster/001/config_vips/used_VIP60"),
	}, grpc.EmptyCallOption{})
	if err != nil {
		return nil, fmt.Errorf("Range error:%s", err.Error())
	}
	for _, v := range rsp.Kvs {
		key := string(v.Key)[strings.LastIndex(string(v.Key), "/")+1:]
		command := "ip add add " + key + "/96 dev " + conf.Interface
		cmd := exec.Command("/bin/bash", "-c", command)
		if _, err := cmd.Output(); err != nil {
			log.Errorf("exec %s error: %s", command, err.Error())
		}
	}
	// /project_nm/version1_0_0/ate/website/domain_id/mail.test01.com/protocol_map/protocol_map_id/WH2I/dst_ip_port/***
	allIPPorts, err := util.GetIPPorts(instance.KVClient)
	if err != nil {
		return nil, fmt.Errorf("GetIPPorts error:%s", err.Error())
	}
	allNodes, err := util.GetNodeIPs(instance.KVClient)
	for ipPort := range allIPPorts {
		command := "ipvsadm -A -t " + ipPort + " -s sh "
		cmd := exec.Command("/bin/bash", "-c", command)
		if _, err := cmd.Output(); err != nil {
			log.Errorf("exec %s error: %s", command, err.Error())
		}
		parts := strings.Split(ipPort, "]:")
		if err := util.AddOrRemoveFirewallRule(parts[0][1:], parts[1], 1); err != nil {
			log.Errorf("AddOrRemoveFirewallRule err:%s", err.Error())
		}
		for _, nodeip := range allNodes {
			command := "ipvsadm -a -t " + ipPort + " -r " + nodeip + " -g "
			cmd := exec.Command("/bin/bash", "-c", command)
			if _, err := cmd.Output(); err != nil {
				log.Errorf("exec %s error: %s", command, err.Error())
			}
		}
	}
	return instance, nil
}

func (h *UpdatedIPHandler) DealEvent(event *mvccpb.Event) {
	if strings.Index(string(event.Kv.Key), "dst_ip_port") > 1 {
		log.Infof("UpdatedIPHandler start")
		nodeIPs, err := util.GetNodeIPs(h.KVClient)
		if err != nil {
			log.Errorf("GetNodeIPs error:%s", err.Error())
			return
		}
		changedIPPort := string(event.Kv.Key)[strings.LastIndex(string(event.Kv.Key), "/")+1:]
		allipports, err := util.GetIPPorts(h.KVClient)
		if err != nil {
			log.Errorf("GetIPPorts error:%s", err.Error())
			return
		}
		found := false
		ipPortCount := 0
		for ipport, count := range allipports {
			if ipport == changedIPPort {
				found = true
				ipPortCount = count
				break
			}
		}
		parts := strings.Split(changedIPPort, "]:")
		if event.GetType().String() == "PUT" {
			if found && ipPortCount == 1 {
				command := "ipvsadm -A -t " + changedIPPort + " -s sh "
				cmd := exec.Command("/bin/bash", "-c", command)
				if _, err := cmd.Output(); err != nil {
					log.Errorf("exec %s error: %s", command, err.Error())
				}
				for _, nodeip := range nodeIPs {
					command := "ipvsadm -a -t " + changedIPPort + " -r " + nodeip + " -g "
					cmd := exec.Command("/bin/bash", "-c", command)
					if _, err := cmd.Output(); err != nil {
						log.Errorf("exec %s error: %s", command, err.Error())
					}
				}
				if err := util.AddOrRemoveFirewallRule(parts[0][1:], parts[1], 1); err != nil {
					log.Errorf("AddOrRemoveFirewallRule err:%s", err.Error())
				}
			}
		} else if event.GetType().String() == "DELETE" {
			if !found {
				command := "ipvsadm -D -t " + changedIPPort
				cmd := exec.Command("/bin/bash", "-c", command)
				if _, err := cmd.Output(); err != nil {
					log.Errorf("exec %s error: %s", command, err.Error())
				}
				if err := util.AddOrRemoveFirewallRule(parts[0][1:], parts[1], 2); err != nil {
					log.Errorf("AddOrRemoveFirewallRule err:%s", err.Error())
				}
			}
		}
	}
}
