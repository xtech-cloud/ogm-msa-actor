package cache

import (
	"ogm-actor/model"
)

type Guard struct {
	Model *model.Guard
}

//TODO use redis/memory
// key is guard_uuid
var guardUUID_guard_map map[string]*Guard

// key is domain_uuid, value is guard_uuid
var domainUUID_guardUUIDS_map map[string]map[string]string

type GuardCAO struct {
}

func NewGuardCAO() *GuardCAO {
	return &GuardCAO{}
}

func (this *GuardCAO) Filter(_domainUUID string) (map[string]string, error) {
	if _, ok := domainUUID_guardUUIDS_map[_domainUUID]; !ok {
		domainUUID_guardUUIDS_map[_domainUUID] = make(map[string]string)
		dao := model.NewGuardDAO(nil)
		guardAry, err := dao.FindByDomain(_domainUUID)
		if nil != err {
			return nil, err
		}
		for _, v := range guardAry {
			domainUUID_guardUUIDS_map[_domainUUID][v.UUID] = v.DeviceUUID
		}
	}
	return domainUUID_guardUUIDS_map[_domainUUID], nil
}

func (this *GuardCAO) Get(_uuid string) (*Guard, error) {
	if _, ok := guardUUID_guard_map[_uuid]; !ok {
		dao := model.NewGuardDAO(nil)
		guard, err := dao.Get(_uuid)
		if nil != err {
			return nil, err
		}
		guardUUID_guard_map[_uuid] = &Guard{
			Model: guard,
		}
	}
	return guardUUID_guard_map[_uuid], nil
}

func (this *GuardCAO) Save(_guard *Guard) (*Guard, error) {

	if _, ok := domainUUID_guardUUIDS_map[_guard.Model.DomainUUID]; !ok {
		domainUUID_guardUUIDS_map[_guard.Model.DomainUUID] = make(map[string]string)
	}
	domainUUID_guardUUIDS_map[_guard.Model.DomainUUID][_guard.Model.UUID] = _guard.Model.DeviceUUID

	// guard的数据仅需要在缓存不存在更新
	if _, ok := guardUUID_guard_map[_guard.Model.UUID]; !ok {
		dao := model.NewGuardDAO(nil)
		if !dao.Exists(_guard.Model.UUID) {
			// 在数据库中插入新值
			err := dao.Upsert(_guard.Model)
			if nil != err {
				return nil, err
			}
		} else {
			// 在数据库中取值
			guardInDB, err := dao.Get(_guard.Model.UUID)
			if nil != err {
				return nil, err
			}
			_guard.Model.Alias = guardInDB.Alias
			_guard.Model.Access = guardInDB.Access
		}
		guardUUID_guard_map[_guard.Model.UUID] = _guard
	}

	return guardUUID_guard_map[_guard.Model.UUID], nil
}

func (this *GuardCAO) Load(_uuid string) error {
	// 从数据库中取出
	dao := model.NewGuardDAO(nil)
	guardInDB, err := dao.Get(_uuid)
	if nil != err {
		return err
	}

	//写入到缓存
	guardUUID_guard_map[_uuid] = &Guard{
		Model: guardInDB,
	}

	if _, ok := domainUUID_guardUUIDS_map[guardInDB.DomainUUID]; !ok {
		domainUUID_guardUUIDS_map[guardInDB.DomainUUID] = make(map[string]string)
	}
	domainUUID_guardUUIDS_map[guardInDB.DomainUUID][guardInDB.UUID] = guardInDB.DeviceUUID
	return nil
}

func (this *GuardCAO) Delete(_uuid string) error {
	//delete(profileMap, _uuid)
	////TODO 删除数据库
	return nil
}
