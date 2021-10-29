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

	dao := model.NewDomainDAO(nil)
	if exists := dao.Exists(_req.Name); exists {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "domain already exists"
		return nil
	}

	domain := &model.Domain{
		UUID: model.NewUUID(),
		Name: _req.Name,
	}
	err := dao.Insert(domain)
	if nil != err {
		return err
	}

	return nil
}

func (this *Domain) Delete(_ctx context.Context, _req *proto.DomainDeleteRequest, _rsp *proto.BlankResponse) error {
	logger.Infof("Received Domain.Delete request: %v", _req)
	_rsp.Status = &proto.Status{}

	dao := model.NewDomainDAO(nil)
	err := dao.Delete(_req.Uuid)
	if nil != err {
		return err
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
		return err
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
