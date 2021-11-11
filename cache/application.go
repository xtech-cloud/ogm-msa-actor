package cache

import (
	"encoding/json"
	"ogm-actor/model"

	"github.com/asim/go-micro/v3/logger"
)

//TODO use redis/memory
// key is domain_uuid
var domainUUID_applicationManifest_map map[string]string
var domainUUID_applicationMD5_map map[string]string

type Application struct {
	Model *model.Application
}

type ApplicationCAO struct {
}

func NewApplicationCAO() *ApplicationCAO {
	return &ApplicationCAO{}
}

func (this *ApplicationCAO) Reload(_domainUUID string) {
	dao := model.NewApplicationDAO(nil)
	_, applicationAry, err := dao.FindByDomain(_domainUUID)
	if err != nil {
		logger.Error(err)
		domainUUID_applicationManifest_map[_domainUUID] = ""
		domainUUID_applicationMD5_map[_domainUUID] = ""
		return
	}
	// 格式化为json
	bytes, err := json.Marshal(applicationAry)
	if err != nil {
		logger.Error(err)
		domainUUID_applicationManifest_map[_domainUUID] = ""
		domainUUID_applicationMD5_map[_domainUUID] = ""
		return
	}
	// 编码为base64
	base64_str := model.ToBase64(bytes)
	domainUUID_applicationManifest_map[_domainUUID] = base64_str
	domainUUID_applicationMD5_map[_domainUUID] = model.ToUUID(base64_str)
}

func (this *ApplicationCAO) GetManifest(_domainUUID string) string {
	if _, ok := domainUUID_applicationManifest_map[_domainUUID]; !ok {
		this.Reload(_domainUUID)
	}

    return domainUUID_applicationManifest_map[_domainUUID]
}

func (this *ApplicationCAO) GetMD5(_domainUUID string) string {
	if _, ok := domainUUID_applicationMD5_map[_domainUUID]; !ok {
		this.Reload(_domainUUID)
	}

    return domainUUID_applicationMD5_map[_domainUUID]
}
