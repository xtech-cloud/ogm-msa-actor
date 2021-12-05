package cache

import (
	"ogm-actor/model"
)

type Device struct {
	Model            *model.Device
	Battery          int32             // 电量
	Volume           int32             // 音量
	Brightness       int32             // 亮度
	Storage          string            // 存储类型
	StorageBlocks    int64             // 存储总容量
	StorageAvailable int64             // 存储可用容量
	Network          string            // 网络类型
	NetworkStrength  int32             // 网络强度
	Program          map[string]string // 程序信息<程序名，程序版本>
	Healthy          int32             // 健康值
}

//TODO use redis/memory
// key is device_uuid
var deviceUUID_device_map map[string]*Device

type DeviceCAO struct {
}

func NewDeviceCAO() *DeviceCAO {
	return &DeviceCAO{}
}

func (this *DeviceCAO) Get(_uuid string) (*Device, error) {
	if _, ok := deviceUUID_device_map[_uuid]; !ok {
		daoDevice := model.NewDeviceDAO(nil)
		device, err := daoDevice.Get(_uuid)
		if nil != err {
			return nil, err
		}
		deviceUUID_device_map[_uuid] = &Device{
			Model:   device,
			Program: make(map[string]string),
		}
	}
	return deviceUUID_device_map[_uuid], nil
}

func (this *DeviceCAO) Save(_device *Device) error {
	// 缓存不存在
	if device, ok := deviceUUID_device_map[_device.Model.UUID]; !ok {
		dao := model.NewDeviceDAO(nil)
		// 在数据库中更新或插入设备实体
		err := dao.Upsert(_device.Model)
		if nil != err {
			return err
		}
	} else {
		// 当值不一致时，更新数据库值
		changed := device.Model.Name != _device.Model.Name ||
			device.Model.OperatingSystem != _device.Model.OperatingSystem ||
			device.Model.SystemVersion != _device.Model.SystemVersion ||
			device.Model.Shape != _device.Model.Shape
		if changed {
			dao := model.NewDeviceDAO(nil)
			err := dao.Update(device.Model)
			if nil != err {
				return err
			}
		}
	}

	deviceUUID_device_map[_device.Model.UUID] = _device
	return nil
}

func (this *DeviceCAO) Delete(_uuid string) error {
	delete(deviceUUID_device_map, _uuid)
	return nil
}
