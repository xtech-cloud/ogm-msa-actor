package handler

import (
	"context"
	"ogm-actor/cache"
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

func (this *Domain) Find(_ctx context.Context, _req *proto.DomainFindRequest, _rsp *proto.DomainFindResponse) error {
	logger.Infof("Received Domain.Find request: %v", _req)
	_rsp.Status = &proto.Status{}

	dao := model.NewDomainDAO(nil)
	domain, err := dao.FindByName(_req.Name)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	if nil == domain {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "domain not found"
		return nil
	}

	_rsp.Domain = &proto.DomainEntity{
		Uuid: domain.UUID,
		Name: domain.Name,
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

func (this *Domain) Search(_ctx context.Context, _req *proto.DomainSearchRequest, _rsp *proto.DomainSearchResponse) error {
	logger.Infof("Received Domain.Search request: %v", _req)
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
	total, domain, err := dao.Search(offset, count, _req.Name)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Total = total

	_rsp.Domain = make([]*proto.DomainEntity, len(domain))
	for i := 0; i < len(domain); i++ {
		_rsp.Domain[i] = &proto.DomainEntity{
			Uuid: domain[i].UUID,
			Name: domain[i].Name,
		}
	}

	return nil
}

func (this *Domain) Update(_ctx context.Context, _req *proto.DomainUpdateRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Domain.Update request: %v", _req)
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

	dao := model.NewDomainDAO(nil)
	err := dao.Update(_req.Uuid, _req.Name)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	return nil
}

func (this *Domain) Execute(_ctx context.Context, _req *proto.DomainExecuteRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Domain.Execute request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Uuid {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "uuid is required"
		return nil
	}

	if "" == _req.Command {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "command is required"
		return nil
	}
	caoDomain := cache.NewDomainCAO()
	domain, err := caoDomain.Get(_req.Uuid)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	for _, sn := range _req.Device {
		// 赋值需要执行的任务
		if _, ok := domain.Task[sn]; !ok {
			domain.Task[sn] = make(map[string]string)
		}
		domain.Task[sn][_req.Command] = _req.Parameter
	}
	return nil
}
