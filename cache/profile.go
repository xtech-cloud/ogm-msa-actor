package cache

import (
	"ogm-actor/model"
)

type Profile struct {
	Model *model.Profile
}

//TODO use redis/memory
// key is profile_uuid
var profileUUID_profile_map map[string]*Profile

// key is domain_uuid, value is profile_uuid
var domainUUID_profileUUIDS_map map[string][]string

type ProfileCAO struct {
}

func NewProfileCAO() *ProfileCAO {
	return &ProfileCAO{}
}

func (this *ProfileCAO) Filter(_domainUUID string) ([]string, error) {
	if _, ok := domainUUID_profileUUIDS_map[_domainUUID]; !ok {
		dao := model.NewProfileDAO(nil)
		profileAry, err := dao.FindByDomain(_domainUUID)
		if nil != err {
			return nil, err
		}
		ary := make([]string, len(profileAry))
		for i, v := range profileAry {
			ary[i] = v.UUID
		}
	}
	return domainUUID_profileUUIDS_map[_domainUUID], nil
}

func (this *ProfileCAO) Get(_uuid string) (*Profile, error) {
	if _, ok := profileUUID_profile_map[_uuid]; !ok {
		dao := model.NewProfileDAO(nil)
		profile, err := dao.Get(_uuid)
		if nil != err {
			return nil, err
		}
		profileUUID_profile_map[_uuid] = &Profile{
			Model: profile,
		}
	}
	return profileUUID_profile_map[_uuid], nil
}

func (this *ProfileCAO) Save(_profile *Profile) error {

	// 缓存不存在
	if profile, ok := profileUUID_profile_map[_profile.Model.UUID]; !ok {
		// 在数据库中更新或插入新值
		dao := model.NewProfileDAO(nil)
		err := dao.Upsert(_profile.Model)
		if nil != err {
			return err
		}
	} else {
		// 当值不一致时，更新数据库值
		changed := profile.Model.Access != _profile.Model.Access ||
			profile.Model.Alias != _profile.Model.Alias
		if changed {
			dao := model.NewProfileDAO(nil)
			err := dao.Update(_profile.Model.UUID, _profile.Model.Access, _profile.Model.Alias)
			if nil != err {
				return err
			}
		}
	}

	profileUUID_profile_map[_profile.Model.UUID] = _profile
	return nil
}

func (this *ProfileCAO) Delete(_uuid string) error {
	//delete(profileMap, _uuid)
	////TODO 删除数据库
	return nil
}
