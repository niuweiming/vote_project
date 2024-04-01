package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// GetVoteCache 从缓存中获取投票信息
func GetVoteCache(c context.Context, id int64) VoteWithOpt {
	var ret VoteWithOpt
	key := fmt.Sprintf("key_vote_%d", id)
	voteStr, err := Rdb.Get(c, key).Result()
	if err == nil && len(voteStr) > 0 {
		//存在数据
		_ = json.Unmarshal([]byte(voteStr), &ret)
		return ret
	}

	fmt.Printf("err:%s", err.Error())
	ret = GetVote(id)
	if ret.Vote.Id > 0 {
		//更新缓存
		s, _ := json.Marshal(ret)
		err1 := Rdb.Set(c, key, s, 36000*time.Second).Err()
		if err1 != nil {
			fmt.Printf("err1:%s", err1.Error())
		}
	}

	return ret
}

// GetVoteHistoryV1 获取用户投票历史信息（版本1）
func GetVoteHistoryV1(c context.Context, userId, voteId int64) []VoteOptUser {
	ret := make([]VoteOptUser, 0)

	//先查下Redis
	k := fmt.Sprintf("vote_user_%d_%d", userId, voteId)
	str, err := Rdb.Get(c, k).Result()
	fmt.Printf("我看看str是什么%s", str)
	if err == nil && len(str) > 0 {
		fmt.Printf("不回溯数据库！\n")
		_ = json.Unmarshal([]byte(str), &ret)
		return ret
	}

	fmt.Printf("回溯数据库！\n")
	//回溯数据库
	if err := Conn.Table("vote_opt_user").Where("vote_id = ? and user_id = ?", voteId, userId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		return ret
	}

	retStr, _ := json.Marshal(ret)
	fmt.Printf("我看看str是什么%s", string(retStr))
	err = Rdb.Set(c, k, retStr, 36000*time.Second).Err()
	if err != nil {
		fmt.Printf("err1:%s\n", err.Error())
	}

	return ret
}
