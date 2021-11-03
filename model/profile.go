package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Profile struct {
	UUID       string `gorm:"column:uuid;type:char(32);primaryKey"`
	Domain     Domain `gorm:"ForeignKey:DomainUUID;AssociationForeignKey:uuid"`
	DomainUUID string `gorm:"column:domain_uuid;type:char(32);not null"`
	Device     Device `gorm:"ForeignKey:DeviceUUID;AssociationForeignKey:uuid"`
	DeviceUUID string `gorm:"column:device_uuid;type:char(32);not null"`
	Alias      string `gorm:"column:alias;type:varchar(128);not null"`
	Access     int32  `gorm:"column:access;type:tinyint;not null;default:0"`
}

func (Profile) TableName() string {
	return "ogm_actor_profile"
}

type ProfileDAO struct {
	conn *Conn
}

func ToProfileUUID(_domainUUID string, _deviceUUID string) string {
	return ToUUID(_domainUUID + _deviceUUID)
}

func NewProfileDAO(_conn *Conn) *ProfileDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &ProfileDAO{
		conn: conn,
	}
}

func (this *ProfileDAO) Insert(_profile *Profile) error {
	db := this.conn.DB
	return db.Create(_profile).Error
}

func (this *ProfileDAO) Upsert(_profile *Profile) error {
	db := this.conn.DB
	// 在冲突时，更新除主键以外的所有列到新值。
	return db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(_profile).Error
}

func (this *ProfileDAO) Get(_uuid string) (*Profile, error) {
	db := this.conn.DB
	var profile Profile
	res := db.Where("uuid = ?", _uuid).First(&profile)
	// 未找到时，返回空值
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, res.Error
}

func (this *ProfileDAO) FindByDomain(_uuid string) ([]Profile, error) {
	db := this.conn.DB
	var profile []Profile
	res := db.Where("domain_uuid = ?", _uuid).Find(&profile)
	return profile, res.Error
}

func (this *ProfileDAO) Exists(_uuid string) bool {
	db := this.conn.DB
	var count int64
	db.Model(&Profile{}).Where("uuid = ?", _uuid).Count(&count)
	return count > 0
}

func (this *ProfileDAO) Update(_uuid string, _access int32, _alias string) error {
	db := this.conn.DB
	res := db.Model(&Profile{}).Where("uuid = ?", _uuid).Updates(Profile{Access: _access, Alias: _alias})
	return res.Error
}
