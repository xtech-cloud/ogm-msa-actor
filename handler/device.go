package handler

import (
	"context"
	"ogm-actor/model"

	proto "github.com/xtech-cloud/ogm-msp-actor/proto/actor"

	"github.com/asim/go-micro/v3/logger"
)

type Device struct{}

func (this *Device) List(_ctx context.Context, _req *proto.ListRequest, _rsp *proto.DeviceListResponse) error {
	logger.Infof("Received Device.List request: %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	if _req.Offset > 0 {
		offset = _req.Offset
	}
	count := int64(0)
	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewDeviceDAO(nil)
	device, err := dao.List(offset, count)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Total = dao.Count()

	_rsp.Device = make([]*proto.DeviceEntity, len(device))
	for i := 0; i < len(device); i++ {
		_rsp.Device[i] = &proto.DeviceEntity{
            Uuid: device[i].UUID,
			SerialNumber: device[i].SerialNumber,
			Name: device[i].Name,
			OperatingSystem: device[i].OperatingSystem,
			SystemVersion: device[i].SystemVersion,
			Shape: device[i].Shape,
		}
	}

	return nil
}

