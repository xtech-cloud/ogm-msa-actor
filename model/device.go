package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (this *DeviceDAO) Upsert(_device *Device) error {
	db := this.conn.DB
	// 在冲突时，更新除主键以外的所有列到新值。
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(_device).Error
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

func (this *DeviceDAO) Get(_uuid string) (*Device, error) {
	db := this.conn.DB
	var device Device
	res := db.Where("uuid = ?", _uuid).First(&device)
	// 未找到时，返回空值
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &device, res.Error
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

func (this *DeviceDAO) Count() int64 {
	db := this.conn.DB
	var count int64
	db.Model(&Device{}).Count(&count)
	return count
}

func (this *DeviceDAO) List(_offset int64, _count int64) ([]*Device, error) {
	db := this.conn.DB
	var device []*Device
	res := db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&device)
	return device, res.Error
}
