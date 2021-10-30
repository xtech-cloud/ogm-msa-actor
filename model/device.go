package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Device struct {
	UUID            string `gorm:"column:uuid;type:char(32);primaryKey"`
	SerialNumber    string `gorm:"column:sn;type:varchar(256);not null;unique"`
	Name            string `gorm:"column:name;type:varchar(256);not null"`
	OperatingSystem string `gorm:"column:os;type:varchar(256);not null"`
	SystemVersion   string `gorm:"column:ver;type:varchar(256);not null"`
	Shape           string `gorm:"column:shape;type:varchar(256);not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

var ErrDeviceExists = errors.New("device exists")

func (Device) TableName() string {
	return "ogm_actor_device"
}

type DeviceDAO struct {
	conn *Conn
}

func NewDeviceDAO(_conn *Conn) *DeviceDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &DeviceDAO{
		conn: conn,
	}
}

func (this *DeviceDAO) Exists(_uuid string) bool {
	db := this.conn.DB
	var count int64
	db.Model(&Device{}).Where("uuid = ?", _uuid).Count(&count)
	return count > 0
}

func (this *DeviceDAO) Insert(_device *Device) error {
	db := this.conn.DB
	return db.Create(_device).Error
}

func (this *DeviceDAO) Update(_device *Device) error {
	db := this.conn.DB
	res := db.Model(&Device{}).Where("uuid = ?", _device.UUID).Updates(
		map[string]interface{}{
			"name":  _device.Name,
			"os":    _device.OperatingSystem,
			"ver":   _device.SystemVersion,
			"shape": _device.Shape,
		})
	return res.Error
}

func (this *DeviceDAO) FindBySN(_serialnumber string) (*Device, error) {
	db := this.conn.DB
	var device Device
	res := db.Where("sn = ?", _serialnumber).First(&device)
	// 未找到时，返回空值
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &device, res.Error
}