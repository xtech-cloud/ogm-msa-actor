package handler

import (
	"context"

	proto "github.com/xtech-cloud/ogm-msp-actor/proto/actor"

	"github.com/asim/go-micro/v3/logger"
)

type Healthy struct{}

// Echo is a single request handler called via client.Call or the generated client code
func (this *Healthy) Echo(_ctx context.Context, _req *proto.EchoRequest, _rsp *proto.EchoResponse) error {
	logger.Infof("Received Healthy.Echo request: %v", _req)
	_rsp.Msg = _req.Msg
	return nil
}
