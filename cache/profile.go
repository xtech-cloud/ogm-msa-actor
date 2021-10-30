package cache

import (
)

type Profile struct {
	Access int32
	Alias  string
}

//TODO use redis/memory
var profileMap map[string]*Profile

type ProfileCAO struct {
}

func NewProfileCAO() *ProfileCAO {
	return &ProfileCAO{}
}

func (this *ProfileCAO) Find(_uuid string) (*Profile, error) {
	profile, _ := profileMap[_uuid]
	return profile, nil
}

func (this *ProfileCAO) Save(_uuid string, _profile *Profile) error {
	profileMap[_uuid] = _profile
	return nil
}

func (this *ProfileCAO) Delete(_uuid string) error {
	delete(profileMap, _uuid)
	return nil
}
