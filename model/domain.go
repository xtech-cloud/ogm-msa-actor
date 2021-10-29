package model

import (
	"errors"
	"gorm.io/gorm"
)

type Domain struct {
	Embedded gorm.Model `gorm:"embedded"`
	UUID     string     `gorm:"column:uuid;type:char(32);not null;unique"`
	Name     string     `gorm:"column:name;type:varchar(256);not null;unique"`
}

var ErrDomainExists = errors.New("domain exists")

func (Domain) TableName() string {
	return "ogm_actor_domain"
}

type DomainDAO struct {
	conn *Conn
}

func NewDomainDAO(_conn *Conn) *DomainDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &DomainDAO{
		conn: conn,
	}
}

func (this *DomainDAO) Insert(_domain *Domain) error {
	db := this.conn.DB
	return db.Create(_domain).Error
}

func (this *DomainDAO) Get(_uuid string) (*Domain, error) {
	db := this.conn.DB
	var domain Domain
	res := db.Where("uuid= ?", _uuid).First(&domain)
	// 未找到时，返回空值
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &domain, res.Error
}

func (this *DomainDAO) Exists(_name string) bool {
	db := this.conn.DB
	var count int64
	db.Model(&Domain{}).Where("name = ?", _name).Count(&count)
	return count > 0
}

func (this *DomainDAO) Count() int64 {
	db := this.conn.DB
	var count int64
	db.Model(&Domain{}).Count(&count)
	return count
}

func (this *DomainDAO) List(_offset int64, _count int64) ([]*Domain, error) {
	db := this.conn.DB
	var domain []*Domain
	res := db.Offset(int(_offset)).Limit(int(_count)).Order("created_at desc").Find(&domain)
	return domain, res.Error
}

func (this *DomainDAO) Delete(_uuid string) error {
	db := this.conn.DB
	return db.Unscoped().Where("uuid = ?", _uuid).Delete(&Domain{}).Error
}
