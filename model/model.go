package model

import "time"

type Network struct {
	Id        int       `xorm:"NOT NULL PK AUTOINCR INT(11)"`
	Name      string    `xorm:"VARCHAR(255) COMMENT('链名称')"`
	Rpc       string    `xorm:"VARCHAR(255) COMMENT('Rpc')"`
	ChainId   int64     `xorm:"NOT NULL UNIQUE INT(20)"`
	Symbol    string    `xorm:"VARCHAR(16) COMMENT('货币符号')"`
	Explorer  string    `xorm:"VARCHAR(255) COMMENT('区块浏览器')"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type Account struct {
	Id        int       `xorm:"NOT NULL PK AUTOINCR INT(11)"`
	Name      string    `xorm:"VARCHAR(255) COMMENT('账户名称')"`
	Address   string    `xorm:"VARCHAR(42) UNIQUE COMMENT('账户地址')"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type Asset struct {
	Id        int       `xorm:"NOT NULL PK AUTOINCR INT(11)"`
	Contract  string    `xorm:"VARCHAR(42) COMMENT('合约地址')"`
	Symbol    string    `xorm:"VARCHAR(16) COMMENT('货币符号')"`
	Decimals  int       `xorm:"NOT NULL DEFAULT 18 INT(11)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}
