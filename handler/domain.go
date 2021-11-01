package handler

import (
	"context"
	"ogm-actor/model"

	proto "github.com/xtech-cloud/ogm-msp-actor/proto/actor"

	"github.com/asim/go-micro/v3/logger"
)

type Domain struct{}

func (this *Domain) Create(_ctx context.Context, _req *proto.DomainCreateRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Domain.Create request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	dao := model.NewDomainDAO(nil)
	if exists := dao.Exists(_req.Name); exists {
		_rsp.Status.Code = 2
		_rsp.Status.Message = "domain already exists"
		return nil
	}

	domain := &model.Domain{
		UUID: model.NewUUID(),
		Name: _req.Name,
	}
	err := dao.Insert(domain)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	return nil
}

func (this *Domain) Delete(_ctx context.Context, _req *proto.DomainDeleteRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Domain.Delete request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewDomainDAO(nil)
	err := dao.Delete(_req.Uuid)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	return nil
}

func (this *Domain) List(_ctx context.Context, _req *proto.ListRequest, _rsp *proto.DomainListResponse) error {
	logger.Infof("Received Domain.List request: %v", _req)
	_rsp.Status = &proto.Status{}

	offset := int64(0)
	if _req.Offset > 0 {
		offset = _req.Offset
	}
	count := int64(0)
	if _req.Count > 0 {
		count = _req.Count
	}

	dao := model.NewDomainDAO(nil)
	domain, err := dao.List(offset, count)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Total = dao.Count()

	_rsp.Domain = make([]*proto.DomainEntity, len(domain))
	for i := 0; i < len(domain); i++ {
		_rsp.Domain[i] = &proto.DomainEntity{
			Uuid: domain[i].UUID,
			Name: domain[i].Name,
		}
	}

	return nil
}

func (this *Domain) Execute(_ctx context.Context, _req *proto.DomainExecuteRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Domain.Execute request: %v", _req)
	_rsp.Status = &proto.Status{}
	return nil
}

func (this *Domain) FetchDevice(_ctx context.Context, _req *proto.DomainFetchDeviceRequest, _rsp *proto.DomainFetchDeviceResponse) error {
	logger.Infof("Received Domain.FetchDevice request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid{
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	dao := model.NewJoinDAO(nil)
	device, err := dao.ListDeviceByDomain(_req.Uuid)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Device = make([]*proto.DeviceEntity, len(device))
	for i := 0; i < len(device); i++ {
		_rsp.Device[i] = &proto.DeviceEntity{
			SerialNumber:    device[i].SerialNumber,
			Name:            device[i].Name,
			OperatingSystem: device[i].OperatingSystem,
			SystemVersion:   device[i].SystemVersion,
			Shape:           device[i].Shape,
		}
	}

    _rsp.Access = make(map[string]int32)
	for i := 0; i < len(device); i++ {
        _rsp.Access[device[i].SerialNumber] = 0
    }

    _rsp.Alias = make(map[string]string)
	for i := 0; i < len(device); i++ {
        _rsp.Alias[device[i].SerialNumber] = ""
    }

	return nil
}

func (this *Domain) AcceptDevice(_ctx context.Context, _req *proto.DomainAcceptDeviceRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Domain.AcceptDevice request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}
	if "" == _req.Device {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "device is required"
		return nil
	}

    dao := model.NewProfileDAO(nil)
	profileUUID := model.ToProfileUUID(_req.Uuid, _req.Device)
    err := dao.UpdateAccess(profileUUID, 1)
	if "" == _req.Device {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	return nil
}

func (this *Domain) RejectDevice(_ctx context.Context, _req *proto.DomainRejectDeviceRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Domain.RejectDevice request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}
	if "" == _req.Device {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "device is required"
		return nil
	}

    dao := model.NewProfileDAO(nil)
	profileUUID := model.ToProfileUUID(_req.Uuid, _req.Device)
    err := dao.UpdateAccess(profileUUID, 2)
	if "" == _req.Device {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	return nil
}
