package model

import (
	"fmt"
	"vote/app/tools"
)

func GetUser(Name string) *User {
	var ret User
	if err := Conn.Table("user").Where("name=?", Name).Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return &ret
}

// 传入一个指针
func CreateUser(user *User) error {
	if err := Conn.Create(user).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return err
	}
	return nil
}

// 原生sql优化
func GetUserV1(name string) *User {
	var ret User
	err := Conn.Raw("select * from user where name = ? limit 1", name).Scan(&ret).Error
	if err != nil {
		tools.Logger.Errorf("[GetUserV1]err:%s", err.Error())
	}
	return &ret
}
