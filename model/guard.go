package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Guard struct {
	UUID       string `gorm:"column:uuid;type:char(32);primaryKey"`
	Domain     Domain `gorm:"ForeignKey:DomainUUID;AssociationForeignKey:uuid"`
	DomainUUID string `gorm:"column:domain_uuid;type:char(32);not null"`
	Device     Device `gorm:"ForeignKey:DeviceUUID;AssociationForeignKey:uuid"`
	DeviceUUID string `gorm:"column:device_uuid;type:char(32);not null"`
    Alias      string `gorm:"column:alias;type:varchar(128);not null;default:''"`
	Access     int32  `gorm:"column:access;type:tinyint;not null;default:0"`
}

func (Guard) TableName() string {
	return "ogm_actor_guard"
}

type GuardDAO struct {
	conn *Conn
}

func ToGuardUUID(_domainUUID string, _deviceUUID string) string {
	return ToUUID(_domainUUID + _deviceUUID)
}

func NewGuardDAO(_conn *Conn) *GuardDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &GuardDAO{
		conn: conn,
	}
}

func (this *GuardDAO) Insert(_guard *Guard) error {
	db := this.conn.DB
	return db.Create(_guard).Error
}

func (this *GuardDAO) Upsert(_guard *Guard) error {
	db := this.conn.DB
	// 在冲突时，不做任何操作
	return db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(_guard).Error
}

func (this *GuardDAO) Get(_uuid string) (*Guard, error) {
	db := this.conn.DB
	var guard Guard
	res := db.Where("uuid = ?", _uuid).First(&guard)
	// 未找到时，返回空值
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &guard, res.Error
}

func (this *GuardDAO) FindByDomain(_uuid string) ([]Guard, error) {
	db := this.conn.DB
	var guard []Guard
	res := db.Where("domain_uuid = ?", _uuid).Find(&guard)
	return guard, res.Error
}

func (this *GuardDAO) Exists(_uuid string) bool {
	db := this.conn.DB
	var count int64
	db.Model(&Guard{}).Where("uuid = ?", _uuid).Count(&count)
	return count > 0
}

func (this *GuardDAO) Update(_uuid string, _access int32, _alias string) error {
	db := this.conn.DB
    // 忽略零值
	res := db.Model(&Guard{}).Where("uuid = ?", _uuid).Updates(Guard{Access: _access, Alias: _alias})
	return res.Error
}
