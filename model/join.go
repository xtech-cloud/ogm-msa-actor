package model

import (
)

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


func (this *JoinDAO) ListDeviceByDomain(_uuid string) ([]*Device, error){
    //TODO 关联表查询
    /*
	db := this.conn.DB
	var device []*Device
    res := db.Preload(Profile{}.TableName()).Find(&device)
    return device, res.Error
    */
	db := this.conn.DB
	var device []*Device
	res := db.Order("created_at desc").Find(&device)
	return device, res.Error
}
