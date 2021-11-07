package model

import ()

type JoinDAO struct {
	conn *Conn
}

func NewJoinDAO(_conn *Conn) *JoinDAO {
	conn := DefaultConn
	if nil != _conn {
		conn = _conn
	}
	return &JoinDAO{
		conn: conn,
	}
}

func (this *JoinDAO) ListDeviceByDomain(_uuid string) ([]Device, error) {
	db := this.conn.DB
	var device []Device
	subQuery1 := db.Model(&Device{})
	subQuery2 := db.Model(&Guard{}).Select("device_uuid").Where("domain_uuid = ?", _uuid)
	res := db.Table("(?) as d, (?) as p", subQuery1, subQuery2).Find(&device)
	return device, res.Error
}
