package handler

import (
	"context"
	"ogm-actor/model"

	proto "github.com/xtech-cloud/ogm-msp-actor/proto/actor"

	"github.com/asim/go-micro/v3/logger"
)

type Application struct{}

func (this *Application) List(_ctx context.Context, _req *proto.ApplicationListRequest, _rsp *proto.ApplicationListResponse) error {
	logger.Infof("Received Application.List request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Domain {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "domain is required"
		return nil
	}

	dao := model.NewApplicationDAO(nil)
	total, application, err := dao.FindByDomain(_req.Domain)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Application = make([]*proto.ApplicationEntity, len(application))
	for i := 0; i < len(application); i++ {
		_rsp.Application[i] = &proto.ApplicationEntity{
			Uuid:     application[i].UUID,
			Name:     application[i].Name,
			Version:  application[i].Version,
			Program:  application[i].Program,
			Location: application[i].Location,
			Url:      application[i].Url,
			Upgrade:  application[i].Upgrade,
		}
	}
	_rsp.Total = total
	return nil
}

func (this *Application) Update(_ctx context.Context, _req *proto.ApplicationUpdateRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Application.Update request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	if "" == _req.Version {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "version is required"
		return nil
	}

	if "" == _req.Url {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "url is required"
		return nil
	}

	if "" == _req.Program {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "program is required"
		return nil
	}

	if "" == _req.Location {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "location is required"
		return nil
	}

	dao := model.NewApplicationDAO(nil)
	application := &model.Application{
		UUID:     _req.Uuid,
		Name:     _req.Name,
		Version:  _req.Version,
		Program:  _req.Program,
		Location: _req.Location,
		Url:      _req.Url,
		Upgrade:  _req.Upgrade,
	}
	err := dao.Update(application)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
	}
	return nil
}

func (this *Application) Add(_ctx context.Context, _req *proto.ApplicationAddRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Application.Add request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Domain {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "domain is required"
		return nil
	}

	if "" == _req.Name {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "name is required"
		return nil
	}

	if "" == _req.Version {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "version is required"
		return nil
	}

	if "" == _req.Url {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "url is required"
		return nil
	}

	if "" == _req.Program {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "program is required"
		return nil
	}

	if "" == _req.Location {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "location is required"
		return nil
	}

	dao := model.NewApplicationDAO(nil)
	applicationUUID := model.ToApplicationUUID(_req.Domain, _req.Name)
	application := &model.Application{
		UUID:       applicationUUID,
		DomainUUID: _req.Domain,
		Name:       _req.Name,
		Version:    _req.Version,
		Program:    _req.Program,
		Location:   _req.Location,
		Url:        _req.Url,
		Upgrade:    _req.Upgrade,
	}
	err := dao.Upsert(application)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
	}
	return nil
}

func (this *Application) Remove(_ctx context.Context, _req *proto.ApplicationRemoveRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Applicaton.Remove request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}
	dao := model.NewApplicationDAO(nil)
	err := dao.Delete(_req.Uuid)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
	}
	return nil
}
