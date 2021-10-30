package cache

import (
	"ogm-actor/model"
)

type Device struct {
	Model            model.Device
	Battery          int32             // 电量
	Volume           int32             // 音量
	Brightness       int32             // 亮度
	Storage          string            // 存储类型
	StorageBlocks    int64             // 存储总容量
	StorageAvailable int64             // 存储可用容量
	Network          string            // 网络类型
	NetworkStrength  int32             // 网络强度
	Program          map[string]string // 程序信息<程序名，程序版本>
}

//TODO use redis/memory
var deviceMap map[string]*Device

type DeviceCAO struct {
}

func NewDeviceCAO() *DeviceCAO {
	return &DeviceCAO{}
}

func (this *DeviceCAO) Find(_uuid string) (*Device, error) {
	device, _ := deviceMap[_uuid]
	return device, nil
}

func (this *DeviceCAO) Save(_uuid string, _entity *Device) error {
	deviceMap[_uuid] = _entity
	return nil
}

func (this *DeviceCAO) Delete(_uuid string) error {
	delete(deviceMap, _uuid)
	return nil
}
