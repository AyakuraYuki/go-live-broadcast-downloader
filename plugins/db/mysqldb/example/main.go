// wiki 文档 https://gorm.io/docs/

package main

import (
	"database/sql"
	"fmt"
	"go-live-broadcast-downloader/plugins/db"
	"go-live-broadcast-downloader/plugins/db/mysqldb"
	"gorm.io/gorm"
	"time"
)

func main() {
	conf := &Config{
		MasterDB: mysqldb.MySQLConfig{
			DSN:          "root:root@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&interpolateParams=true&parseTime=true&timeout=200ms&writeTimeout=2s&readTimeout=2s",
			MaxIdleConns: 32,
			MaxIdleTime:  30,
			MaxLifetime:  30,
			MaxOpenConns: 128,
			LogLevel:     db.LogLevelWarn,
		},
	}
	dao := NewDao(conf)
	defer dao.Close()

	for i := 0; i < 200; i++ {
		go func() {
			for {
				time.Sleep(time.Millisecond * 500)
				uinfos, _ := dao.Say()
				fmt.Println(uinfos)
			}
		}()

	}
	select {}
}

type Config struct {
	MasterDB mysqldb.MySQLConfig
}

type Dao struct {
	conf *Config
	mDB  *gorm.DB
}

func NewDao(c *Config) *Dao {
	d := &Dao{
		conf: c,
	}

	mdb, err := mysqldb.DB(c.MasterDB)
	if err != nil {
		panic(err)
	}
	d.mDB = mdb
	return d
}

func (d *Dao) Close() {
	if sqldb, err := d.mDB.DB(); err == nil {
		_ = sqldb.Close()
	}
	time.Sleep(time.Second * 3)
}

func (d *Dao) Say() (uinfos []*UserInfo, err error) {
	//a := d.mDB.Exec("UPDATE user_info SET avatar = ? WHERE id IN ?", "头像", []int64{1, 2, 3}).RowsAffected
	err = d.mDB.Raw("SELECT id, user_id,user_name,avatar FROM user_info WHERE user_id IN ?", []int64{1, 2, 3}).Find(&uinfos).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return uinfos, err

	//事务
	//tx := d.mDB.Begin()
	//if tx.Error != nil {
	//	return "", err
	//}
	//
	//defer func() {
	//	if err != nil {
	//		tx.Rollback()
	//	}
	//}()
	//
	//err = tx.Debug().Exec("update user_name=? where user_id in ?", "昵称", []int64{1, 2, 3}).Error
	//if err != nil {
	//	return "", err
	//}
	//
	//return "", nil
}

type UserInfo struct {
	UserId     uint64         `json:"userId,omitempty" gorm:"column:user_id"`
	UserName   string         `json:"userName,omitempty" gorm:"column:user_name"`
	Avatar     sql.NullString `json:"avatar,omitempty" gorm:"column:avatar"`
	CreateTime string         `json:"createTime,omitempty" gorm:"column:create_time"`
	UpdateTime string         `json:"updateTime,omitempty" gorm:"column:update_time"`
}
