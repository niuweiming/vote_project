package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//数据库操作
//数据库表之间的关系
//表之间的关系： 投票和选项之间是1：n 用户和投票之间是n:n的关系，必须三张表

var Conn *gorm.DB

// 连接数据库
func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", "root", "123456", "47.236.102.175:3307", "vote")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	Conn = conn
}

func Close() {
	db, _ := Conn.DB()
	_ = db.Close()

}
