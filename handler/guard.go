package handler

import (
	"context"
	"ogm-actor/model"

	proto "github.com/xtech-cloud/ogm-msp-actor/proto/actor"

	"github.com/asim/go-micro/v3/logger"
)

type Guard struct{}

func (this *Guard) Fetch(_ctx context.Context, _req *proto.GuardFetchRequest, _rsp *proto.GuardFetchResponse) error {
	logger.Infof("Received Guard.Fetch request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Domain {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "domain is required"
		return nil
	}

	daoJoin := model.NewJoinDAO(nil)
	device, err := daoJoin.ListDeviceByDomain(_req.Domain)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Device = make([]*proto.DeviceEntity, len(device))
	for i := 0; i < len(device); i++ {
		_rsp.Device[i] = &proto.DeviceEntity{
			Uuid:            device[i].UUID,
			SerialNumber:    device[i].SerialNumber,
			Name:            device[i].Name,
			OperatingSystem: device[i].OperatingSystem,
			SystemVersion:   device[i].SystemVersion,
			Shape:           device[i].Shape,
		}
	}

	daoGuard := model.NewGuardDAO(nil)
	guard, err := daoGuard.FindByDomain(_req.Domain)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Access = make(map[string]int32)
	for i := 0; i < len(guard); i++ {
		_rsp.Access[guard[i].DeviceUUID] = guard[i].Access
	}

	_rsp.Alias = make(map[string]string)
	for i := 0; i < len(guard); i++ {
		_rsp.Alias[guard[i].DeviceUUID] = guard[i].Alias
	}

	return nil
}

func (this *Guard) Edit(_ctx context.Context, _req *proto.GuardEditRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Guard.Edit request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Domain {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "domain is required"
		return nil
	}

	if "" == _req.Device {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "device is required"
		return nil
	}

	dao := model.NewGuardDAO(nil)
	guardUUID := model.ToGuardUUID(_req.Domain, _req.Device)
	err := dao.Update(guardUUID, _req.Access, _req.Alias)
	if "" == _req.Device {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	return nil
}

func (this *Guard) Delete(_ctx context.Context, _req *proto.GuardDeleteRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Guard.Delete request: %v", _req)
	_rsp.Status = &proto.Status{}
	return nil
}
