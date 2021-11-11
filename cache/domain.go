package cache

import (
	"ogm-actor/model"
)

type Domain struct {
	Model    *model.Domain
	Property map[string]string
	// map[device.sn]map[command]parameter
	Task map[string]map[string]string
}

//TODO use redis/memory
// key is domain_uuid
var domainUUID_domain_map map[string]*Domain

type DomainCAO struct {
}

func NewDomainCAO() *DomainCAO {
	return &DomainCAO{}
}

func (this *DomainCAO) Get(_uuid string) (*Domain, error) {
	if _, ok := domainUUID_domain_map[_uuid]; !ok {
		dao := model.NewDomainDAO(nil)
		domain, err := dao.Get(_uuid)
		if nil != err {
			return nil, err
		}
		domainUUID_domain_map[_uuid] = &Domain{
			Model:    domain,
			Property: make(map[string]string),
			Task:     make(map[string]map[string]string),
		}
	}
	return domainUUID_domain_map[_uuid], nil
}

func (this *DomainCAO) Save(_domain *Domain) error {
	// 缓存不存在
	if domain, ok := domainUUID_domain_map[_domain.Model.UUID]; !ok {
		dao := model.NewDomainDAO(nil)
		// 在数据库中更新或插入设备实体
		err := dao.Upsert(_domain.Model)
		if nil != err {
			return err
		}
	} else {
		// 当值不一致时，更新数据库值
		changed := domain.Model.Name != _domain.Model.Name
		if changed {
			dao := model.NewDomainDAO(nil)
			err := dao.Update(_domain.Model.UUID, _domain.Model.Name)
			if nil != err {
				return err
			}
		}
	}

	domainUUID_domain_map[_domain.Model.UUID] = _domain
	return nil
}

func (this *DomainCAO) Delete(_uuid string) error {
	//delete(deviceUUID_device_map, _uuid)
	return nil
}
