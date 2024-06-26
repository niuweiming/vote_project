package model

import (
	"context"
	"fmt"
	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"vote/config"
)

//数据库操作
//数据库表之间的关系
//表之间的关系： 投票和选项之间是1：n 用户和投票之间是n:n的关系，必须三张表

var Conn *gorm.DB
var Rdb *redis.Client

// 连接mysql数据库
func NewMysql() {
	//my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", "root", "123456", "47.236.102.175:3307", "vote")
	my := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Config.Db.Username,
		config.Config.Db.Password,
		config.Config.Db.Host,
		config.Config.Db.Port,
		config.Config.Db.Db,
		config.Config.Db.Charset)
	fmt.Println("我看看有没有取到", config.Config.Db.Username)
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	Conn = conn
}

// 连接redis
func NewRdb() {
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     "47.236.102.175:6379",
	//	Password: "123456", //no password set
	//	DB:       0,        //use default DB
	//})
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Address,
		Password: config.Config.Redis.Password, //no password set
		DB:       0,                            //use default DB
	})
	fmt.Println(config.Config.Redis.Address)
	Rdb = rdb
	store, _ = redisstore.NewRedisStore(context.TODO(), Rdb)

	return
}

func Close() {
	db, _ := Conn.DB()
	_ = Rdb.Close()
	_ = db.Close()

}
