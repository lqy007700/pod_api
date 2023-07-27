package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zxnlx/common"
	"github.com/zxnlx/pod/proto/pod"
	"github.com/zxnlx/pod_api/plugin/form"
	"github.com/zxnlx/pod_api/proto/pod_api"
	"strconv"
)

type PodApi struct {
	PodService pod.PodService
}

func (p *PodApi) FindPodById(ctx context.Context, req *pod_api.Request, resp *pod_api.Response) error {
	common.Info("pod-api 接收到 findPodById 请求")
	if _, ok := req.Get["pod_id"]; !ok {
		common.Error("pod_id 不存在")
		resp.Code = 500
		return errors.New("podid 不存在")
	}

	podIdString := req.Get["pod_id"].Values[0]
	podId, err := strconv.ParseInt(podIdString, 10, 64)
	if err != nil {
		return err
	}

	podInfo, err := p.PodService.FindPodById(ctx, &pod.PodId{Id: podId})
	if err != nil {
		return err
	}

	resp.Code = 0
	marshal, err := json.Marshal(podInfo)
	if err != nil {
		return err
	}
	resp.Body = string(marshal)
	return nil
}

func (p *PodApi) AddPod(ctx context.Context, req *pod_api.Request, resp *pod_api.Response) error {
	common.Info("pod-api 接收到 AddPod 请求")
	addPodInfo := &pod.PodInfo{}

	dataSlice, ok := req.Post["pod_port"]

	if ok {
		//特殊处理
		var podSlice []*pod.PodPort
		for _, v := range dataSlice.Values {
			i, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				common.Error(err)
			}
			port := &pod.PodPort{
				ContainerPort: int32(i),
				Protocol:      "TCP",
			}
			podSlice = append(podSlice, port)
		}
		addPodInfo.PodPort = podSlice
	}

	//form类型转化到结构体中
	form.FromToPodStruct(req.Post, addPodInfo)

	response, err := p.PodService.AddPod(ctx, addPodInfo)
	if err != nil {
		common.Error(err)
		return err
	}

	resp.Code = 200
	b, _ := json.Marshal(response)
	resp.Body = string(b)
	return nil
}

func (p *PodApi) DeletePodById(ctx context.Context, req *pod_api.Request, resp *pod_api.Response) error {
	fmt.Println("接受到 podApi.DeletePodById 的请求")
	if _, ok := req.Get["pod_id"]; !ok {
		return errors.New("参数异常")
	}
	//获取要删除的ID
	podIdString := req.Get["pod_id"].Values[0]
	podId, err := strconv.ParseInt(podIdString, 10, 64)
	if err != nil {
		common.Error(err)
		return err
	}
	//删除指定服务
	response, err := p.PodService.DelPod(ctx, &pod.PodId{
		Id: podId,
	})
	if err != nil {
		common.Error(err)
		return err
	}

	resp.Code = 200
	b, _ := json.Marshal(response)
	resp.Body = string(b)
	return nil
}

func (p *PodApi) UpdatePod(ctx context.Context, request *pod_api.Request, response *pod_api.Response) error {
	//TODO implement me
	panic("implement me")
}

func (p *PodApi) Call(ctx context.Context, request *pod_api.Request, response *pod_api.Response) error {
	//TODO implement me
	panic("implement me")
}
