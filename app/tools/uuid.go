package tools

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

var snowNode *snowflake.Node

// 使用雪花算法生成uuid
func GetUUID() string {
	id := uuid.New() //默认版本，基于一个随机数
	fmt.Printf("uuid:%s,version:%s", id.String(), id.Version().String())
	return id.String()
}

func GetUid() int64 {
	if snowNode == nil {
		snowNode, _ = snowflake.NewNode(1)
	}

	return snowNode.Generate().Int64()
}
