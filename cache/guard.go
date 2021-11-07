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
var domainUUID_guardUUIDS_map map[string][]string

type GuardCAO struct {
}

func NewGuardCAO() *GuardCAO {
	return &GuardCAO{}
}

func (this *GuardCAO) Filter(_domainUUID string) ([]string, error) {
	if _, ok := domainUUID_guardUUIDS_map[_domainUUID]; !ok {
		dao := model.NewGuardDAO(nil)
		guardAry, err := dao.FindByDomain(_domainUUID)
		if nil != err {
			return nil, err
		}
		ary := make([]string, len(guardAry))
		for i, v := range guardAry {
			ary[i] = v.UUID
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

func (this *GuardCAO) Save(_guard *Guard) error {

	// 缓存不存在
	if guard, ok := guardUUID_guard_map[_guard.Model.UUID]; !ok {
		// 在数据库中更新或插入新值
		dao := model.NewGuardDAO(nil)
		err := dao.Upsert(_guard.Model)
		if nil != err {
			return err
		}
	} else {
		// 当值不一致时，更新数据库值
		changed := guard.Model.Access != _guard.Model.Access ||
			guard.Model.Alias != _guard.Model.Alias
		if changed {
			dao := model.NewGuardDAO(nil)
			err := dao.Update(_guard.Model.UUID, _guard.Model.Access, _guard.Model.Alias)
			if nil != err {
				return err
			}
		}
	}

	guardUUID_guard_map[_guard.Model.UUID] = _guard
	return nil
}

func (this *GuardCAO) Delete(_uuid string) error {
	//delete(profileMap, _uuid)
	////TODO 删除数据库
	return nil
}
