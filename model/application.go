package model

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Application struct {
	UUID       string `gorm:"column:uuid;type:char(32);primaryKey"`
	DomainUUID string `gorm:"column:domain_uuid;type:char(32);not null"`
	Name       string `gorm:"column:name;type:varchar(256);not null"`
	Version    string `gorm:"column:version;type:varchar(64);not null"`
	Program    string `gorm:"column:program;type:varchar(512);not null"`
	Location   string `gorm:"column:location;type:varchar(512);not null"`
	Url        string `gorm:"column:url;type:varchar(1024);not null"`
}

func (Application) TableName() string {
	return "ogm_actor_application"
}

type ApplicationDAO struct {
	conn *Conn
}

func ToApplicationUUID(_domainUUID string, _name string) string {
	return ToUUID(_domainUUID + _name)
}

func NewApplicationDAO(_conn *Conn) *ApplicationDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &ApplicationDAO{
		conn: conn,
	}
}

func (this *ApplicationDAO) Upsert(_application *Application) error {
	db := this.conn.DB
	// 在冲突时，更新除主键以外的所有列到新值。
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(_application).Error
}

func (this *ApplicationDAO) Get(_uuid string) (*Application, error) {
	db := this.conn.DB
	var application Application
	res := db.Where("uuid = ?", _uuid).First(&application)
	// 未找到时，返回空值
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &application, res.Error
}

func (this *ApplicationDAO) FindByDomain(_uuid string) (int64, []Application, error) {
	db := this.conn.DB
	db = db.Where("domain_uuid = ?", _uuid)
    var total int64
    res := db.Model(&Application{}).Count(&total)
    if res.Error != nil {
        return 0, nil, res.Error
    }
	var application []Application
	res = db.Find(&application)
	return total, application, res.Error
}

func (this *ApplicationDAO) Delete(_uuid string) error {
	db := this.conn.DB
	return db.Where("uuid = ?", _uuid).Delete(&Application{}).Error
}
