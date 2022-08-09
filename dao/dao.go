package dao

import (
	"fwallet/model"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Dao struct {
	orm *xorm.Engine
}

func New(dataSource string) *Dao {
	orm, err := xorm.NewEngine("sqlite3", dataSource)
	if err != nil {
		panic(err)
	}
	return &Dao{
		orm: orm,
	}
}

func (d *Dao) Close() {
	if d.orm != nil {
		_ = d.orm.Close()
	}
}

func (d *Dao) Sync2() {
	if err := d.orm.Sync2(new(model.Network),
		new(model.Account),
		new(model.Asset)); err != nil {
		panic(err)
	}
}

func (d *Dao) Insert(sess *xorm.Session, bean interface{}) error {
	var (
		affected int64
		err      error
	)
	if sess != nil {
		affected, err = sess.Insert(bean)
	} else {
		affected, err = d.orm.Insert(bean)
	}

	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.Errorf("wrong affected:%d", affected)
	}
	return nil
}

func (d *Dao) Inserts(sess *xorm.Session, beans interface{}) (int64, error) {
	var (
		affected int64
		err      error
	)
	if sess != nil {
		affected, err = sess.Insert(beans)
	} else {
		affected, err = d.orm.Insert(beans)
	}

	if err != nil {
		return affected, err
	}
	return affected, nil
}

func (d *Dao) Delete(sess *xorm.Session, id int64, bean interface{}) error {
	var (
		affected int64
		err      error
	)
	if sess != nil {
		affected, err = sess.ID(id).Delete(bean)
	} else {
		affected, err = d.orm.ID(id).Delete(bean)
	}

	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.Errorf("wrong affected:%d", affected)
	}
	return nil
}
