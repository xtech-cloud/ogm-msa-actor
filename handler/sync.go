package handler

import (
	"context"
	"ogm-actor/cache"
	"ogm-actor/model"

	proto "github.com/xtech-cloud/ogm-msp-actor/proto/actor"

	"github.com/asim/go-micro/v3/logger"
)

type Sync struct{}

// 推送的频率会非常高，需要使用缓存，尽量减少数据库的直接访问
func (this *Sync) Push(_ctx context.Context, _req *proto.SyncPushRequest, _rsp *proto.SyncPushResponse) error {
	logger.Infof("Received Sync.Push request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Domain {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "domain is required"
		return nil
	}

	if nil == _req.Device {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "device is required"
		return nil
	}

	if "" == _req.Device.SerialNumber {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "device.serialnumber is required"
		return nil
	}

	program := make(map[string]string)
	for k, v := range _req.Device.Program {
		program[k] = v
	}

	deviceUUID := model.ToUUID(_req.Device.SerialNumber)
	device := &cache.Device{
		Model: model.Device{
			UUID:            deviceUUID,
			SerialNumber:    _req.Device.SerialNumber,
			Name:            _req.Device.Name,
			OperatingSystem: _req.Device.OperatingSystem,
			SystemVersion:   _req.Device.SystemVersion,
			Shape:           _req.Device.Shape,
		},
		Battery:          _req.Device.Battery,
		Volume:           _req.Device.Volume,
		Brightness:       _req.Device.Brightness,
		Storage:          _req.Device.Storage,
		StorageBlocks:    _req.Device.StorageBlocks,
		StorageAvailable: _req.Device.StorageAvailable,
		Network:          _req.Device.Network,
		NetworkStrength:  _req.Device.NetworkStrength,
		Program:          program,
	}
	profileUUID := model.ToUUID(_req.Domain + deviceUUID)

	// 在缓存中查找
	caoDevice := cache.NewDeviceCAO()
	deviceInCache, err := caoDevice.Find(deviceUUID)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	// 缓存和数据库中都没有时，在数据库中插入新值
	if nil == deviceInCache {
		daoDevice := model.NewDeviceDAO(nil)
		if !daoDevice.Exists(deviceUUID) {
			// 插入设备实体
			err = daoDevice.Insert(&device.Model)
			if nil != err {
				_rsp.Status.Code = -1
				_rsp.Status.Message = err.Error()
				return nil
			}
			// 插入简介实体
			daoProfile := model.NewProfileDAO(nil)
			if !daoProfile.Exists(profileUUID) {
				profile := &model.Profile{
					UUID:   profileUUID,
					Domain: _req.Domain,
					Device: deviceUUID,
					Access: 0,
					Alias:  "",
				}
				err = daoProfile.Insert(profile)
				if nil != err {
					_rsp.Status.Code = -1
					_rsp.Status.Message = err.Error()
					return nil
				}
			}
		}
	} else {
		// 当值不一致时，更新数据库值
		changed := deviceInCache.Model.Name != device.Model.Name ||
			deviceInCache.Model.OperatingSystem != device.Model.OperatingSystem ||
			deviceInCache.Model.SystemVersion != device.Model.SystemVersion ||
			deviceInCache.Model.Shape != device.Model.Shape
		if changed {
			daoDevice := model.NewDeviceDAO(nil)
			err := daoDevice.Update(&device.Model)
			if nil != err {
				_rsp.Status.Code = -1
				_rsp.Status.Message = err.Error()
				return nil
			}
		}
	}

	//在缓存中更新设备
	caoDevice.Save(deviceUUID, device)

	// 查找缓存是否包含此简介
	caoProfile := cache.NewProfileCAO()
	profileInCache, err := caoProfile.Find(profileUUID)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	// 缓存没有值时，从数据库中取值
	if nil == profileInCache {
		daoProfile := model.NewProfileDAO(nil)
		profile, err := daoProfile.Get(profileUUID)
		if nil != err {
			_rsp.Status.Code = -1
			_rsp.Status.Message = err.Error()
			return nil
		}
		profileInCache = &cache.Profile{
			Access: profile.Access,
			Alias:  profile.Alias,
		}
	}
	//在缓存中更新简介
	caoProfile.Save(profileUUID, profileInCache)

	//在缓存中更新属性
	caoProperty := cache.NewPropertyCAO()
	if nil != _req.UpProperty {
		for k, v := range _req.UpProperty {
			caoProperty.Save(k, v)
		}
	}

	// 赋值回复
	_rsp.Access = profileInCache.Access
	_rsp.Alias = profileInCache.Alias

	_rsp.Property = make(map[string]string)
	if nil != _req.DownProperty {
		for _, k := range _req.DownProperty {
			v, ok, err := caoProperty.Find(k)
			if nil != err {
				_rsp.Status.Code = -1
				_rsp.Status.Message = err.Error()
				return nil
			}
			// 仅返回存在的属性
			if !ok {
				continue
			}
			_rsp.Property[k] = v
		}
	}

	return nil
}

func (this *Sync) Pull(_ctx context.Context, _req *proto.SyncPullRequest, _rsp *proto.SyncPullResponse) error {
	logger.Infof("Received Sync.Pull request: %v", _req)
	_rsp.Status = &proto.Status{}

	if "" == _req.Domain {
		_rsp.Status.Code = 1
		_rsp.Status.Message = "domain is required"
		return nil
	}

	dao := model.NewJoinDAO(nil)
	device, err := dao.ListDeviceByDomain(_req.Domain)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Device = make([]*proto.DeviceEntity, len(device))
	for i := 0; i < len(device); i++ {
		_rsp.Device[i] = &proto.DeviceEntity{
			SerialNumber:    device[i].SerialNumber,
			Name:            device[i].Name,
			OperatingSystem: device[i].OperatingSystem,
			SystemVersion:   device[i].SystemVersion,
			Shape:           device[i].Shape,
		}
	}

	_rsp.Property = make(map[string]string)
	caoProperty := cache.NewPropertyCAO()
	if nil != _req.DownProperty {
		for _, k := range _req.DownProperty {
			v, ok, err := caoProperty.Find(k)
			if nil != err {
				_rsp.Status.Code = -1
				_rsp.Status.Message = err.Error()
				return nil
			}
			// 仅返回存在的属性
			if !ok {
				continue
			}
			_rsp.Property[k] = v
		}
	}

	return nil
}
