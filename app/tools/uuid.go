package tools

import (
	"fmt"
	"github.com/google/uuid"
)

// 使用雪花算法生成uuid
func GetUUID() string {
	id := uuid.New() //默认版本，基于一个随机数
	fmt.Printf("uuid:%s,version:%s", id.String(), id.Version().String())
	return id.String()
}
