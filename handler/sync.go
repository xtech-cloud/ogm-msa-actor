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
		Model: &model.Device{
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
	guardUUID := model.ToGuardUUID(_req.Domain, deviceUUID)
	guard := &cache.Guard{
		Model: &model.Guard{
			UUID:       guardUUID,
			DomainUUID: _req.Domain,
			DeviceUUID: deviceUUID,
			Access:     0,
			Alias:      "",
		},
	}

	//在缓存中更新设备
	caoDevice := cache.NewDeviceCAO()
	caoDevice.Save(device)

	//在缓存中更新守卫
	caoGuard := cache.NewGuardCAO()
	caoGuard.Save(guard)

	//在缓存中更新属性
	caoDomain := cache.NewDomainCAO()
	domain, err := caoDomain.Get(_req.Domain)
	if "" == _req.Device.SerialNumber {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	if nil != _req.UpProperty {
		for k, v := range _req.UpProperty {
			domain.Property[k] = v
		}
	}

	// 赋值回复
	_rsp.Access = guard.Model.Access
	_rsp.Alias = guard.Model.Alias

	_rsp.Property = make(map[string]string)
	if nil != _req.DownProperty {
		for _, k := range _req.DownProperty {
			if v, ok := domain.Property[k]; ok {
				_rsp.Property[k] = v
			}
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

	//TODO 仅拉取允许访问的设备

	caoGuard := cache.NewGuardCAO()
	caoDevice := cache.NewDeviceCAO()
	profileAry, err := caoGuard.Filter(_req.Domain)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}

	_rsp.Device = make([]*proto.DeviceEntity, len(profileAry))
	for i, v := range profileAry {
		profile, err := caoGuard.Get(v)
		if nil != err {
			_rsp.Status.Code = -1
			_rsp.Status.Message = err.Error()
			return nil
		}
		device, err := caoDevice.Get(profile.Model.DeviceUUID)
		if nil != err {
			_rsp.Status.Code = -1
			_rsp.Status.Message = err.Error()
			return nil
		}
		_rsp.Device[i] = &proto.DeviceEntity{
			SerialNumber:    device.Model.SerialNumber,
			Name:            device.Model.Name,
			OperatingSystem: device.Model.OperatingSystem,
			SystemVersion:   device.Model.SystemVersion,
			Shape:           device.Model.Shape,
		}
	}

	_rsp.Property = make(map[string]string)
	caoDomain := cache.NewDomainCAO()
	domain, err := caoDomain.Get(_req.Domain)
	if nil != err {
		_rsp.Status.Code = -1
		_rsp.Status.Message = err.Error()
		return nil
	}
	if nil != _req.DownProperty {
		for _, k := range _req.DownProperty {
			if v, ok := domain.Property[k]; ok {
				_rsp.Property[k] = v
			}
		}
	}

	return nil
}
