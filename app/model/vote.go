package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"time"
	"vote/app/tools"
)

// GetVotes 从数据库中检索所有投票。
func GetVotes() []Vote {
	ret := make([]Vote, 0)
	if err := Conn.Table("vote").Find(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

// GetVote 根据ID检索单个投票以及其选项。
func GetVote(id int64) VoteWithOpt {
	var ret Vote
	if err := Conn.Table("vote").Where("id=?", id).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	opt := make([]VoteOpt, 0)
	if err := Conn.Table("vote_opt").Where("vote_id = ?", id).Find(&opt).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}

	return VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}
}

func GetVoteV2(id int64) (*VoteWithOptV1, error) {
	var ret VoteWithOptV1
	err := Conn.Preload("vote_opt", "vote_id = ?", id).Raw("select * from vote where id = ? ", id).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		return &ret, err
	}
	return &ret, nil
}

func GetVoteV5(id int64) (*VoteWithOpt, error) {
	var ret Vote

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Conn.Raw("select * from vote where id = ?", id).Scan(&ret).Error
		if err != nil {
			fmt.Printf("err:%s", err.Error())
		}
	}()
	wg.Add(1)
	opt := make([]VoteOpt, 0)
	go func() {
		defer wg.Done()
		err1 := Conn.Raw("select * from vote_opt where vote_id = ?", id).Scan(&opt).Error
		if err1 != nil {
			fmt.Printf("err:%s", err1.Error())
		}
	}()
	wg.Wait()

	return &VoteWithOpt{
		Vote: ret,
		Opt:  opt,
	}, nil
}

// DoVote 使用GORM执行事务性投票操作。
func DoVote(userId, voteId int64, optIDs []int64) bool {
	tx := Conn.Begin()
	var ret Vote
	if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
	}
	var oldVoteUser VoteOptUser
	if err := tx.Table("vote_opt_user").Where("vote_id = ? and user_id = ?", voteId, userId).First(&oldVoteUser).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
	}
	if oldVoteUser.Id > 0 {
		fmt.Printf("用户已投票")
		tx.Rollback()
	}

	for _, value := range optIDs {
		if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
		user := VoteOptUser{
			VoteId:      voteId,
			UserId:      userId,
			VoteOptId:   value,
			CreatedTime: time.Now(),
		}
		if err := tx.Create(&user).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
	}

	tx.Commit()

	return true
}

// DoVoteV1 原生SQL
func DoVoteV1(userId, voteId int64, optIDs []int64) bool {
	Conn.Exec("begin").
		Exec("select * from vote where id = ?", voteId).
		Exec("commit")
	return false
}

// DoVoteV2 匿名函数最常用的写法 利用了匿名函数实现事务
func DoVoteV2(userId, voteId int64, optIDs []int64, context *gin.Context) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		var ret Vote
		if err := tx.Table("vote").Where("id = ?", voteId).First(&ret).Error; err != nil {
			fmt.Printf("没有这个投票选项 err:%s", err.Error())
			context.JSON(http.StatusOK, tools.ECode{
				Message: "没有这个投票选项",
			})
			return err
		}

		for _, value := range optIDs {
			if err := tx.Table("vote_opt").Where("id = ?", value).Update("count", gorm.Expr("count + ?", 1)).Error; err != nil {
				fmt.Printf("err:%s", err.Error())
				return err
			}
			user := VoteOptUser{
				VoteId:      voteId,
				UserId:      userId,
				VoteOptId:   value,
				CreatedTime: time.Now(),
			}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}

	return true
}

// DovoteV3 原生sql的优化
func DoVoteV3(userId, voteId int64, optIDs []int64) bool {
	tx := Conn.Begin()
	var ret Vote
	if err := tx.Raw("select * from vote where id = ?", voteId).Scan(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
	}

	var oldVoteUser VoteOptUser
	if err := tx.Raw("select * from vote_opt_user where vote_id = ? and user_id = ?", voteId, userId).Scan(&oldVoteUser).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
		tx.Rollback()
	}
	if oldVoteUser.Id > 0 {
		fmt.Printf("用户已投票")
		tx.Rollback()
	}

	for _, value := range optIDs {
		if err := tx.Exec("update vote_opt set count = count+1 where id = ? limit 1", value).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
		user := VoteOptUser{
			VoteId:      voteId,
			UserId:      userId,
			VoteOptId:   value,
			CreatedTime: time.Now(),
		}
		if err := tx.Create(&user).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			tx.Rollback()
		}
	}

	tx.Commit()

	return true
}

// AddVote 在事务内向数据库添加新的投票及其选项。
func Addvote(vote Vote, opt []VoteOpt) error {
	//事务
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&vote).Error; err != nil {
			return err
		}
		for _, voteOpt := range opt {
			voteOpt.VoteId = vote.Id
			if err := tx.Create(&voteOpt).Error; err != nil {
				return err
			}
		}
		return nil

	})
	return err
}

// UpdateVote 在事务内更新数据库中现有投票及其选项。
func Updatevote(vote Vote, opt []VoteOpt) error {
	//事务
	err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&vote).Error; err != nil {
			return err
		}
		for _, voteOpt := range opt {
			if err := tx.Save(&voteOpt).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// DelVote 在事务内从数据库中删除一个投票及其相关数据。
func DelVote(id int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Vote{}, id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		if err := tx.Where("vote_id = ?", id).Delete(&VoteOpt{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		if err := tx.Where("vote_id = ?", id).Delete(&VoteOptUser{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}
	return true
}

// 删除版本V1原生sql
func DelVoteV1(id int64) bool {
	if err := Conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("delete from vote where id = ? limit 1", id).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		if err := tx.Exec("delete from vote_opt where vote_id = ?", id).Delete(&VoteOpt{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}

		if err := tx.Exec("delete from vote_opt_user where vote_id = ?", id).Delete(&VoteOptUser{}).Error; err != nil {
			fmt.Printf("err:%s", err.Error())
			return err
		}
		return nil
	}); err != nil {
		fmt.Printf("err:%s", err.Error())
		return false
	}
	return true
}

// GetVoteHistory 检索用户对特定投票的投票历史。
func GetVoteHistory(userId, voteId int64) []VoteOptUser {
	ret := make([]VoteOptUser, 0)
	if err := Conn.Table("vote_opt_user").Where("vote_id = ? and user_id = ?", voteId, userId).First(&ret).Error; err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	return ret
}

// EndVote 根据过期时间将已过期的投票的状态更新为0。
func EndVote() {
	votes := make([]Vote, 0)
	if err := Conn.Table("vote").Where("status = ?", 1).Find(&votes).Error; err != nil {
		return
	}

	now := time.Now().Unix()
	for _, vote := range votes {
		if vote.Time+vote.CreatedTime.Unix() <= now {
			Conn.Table("vote").Where("id = ?", vote.Id).Update("status", 0)
		}
	}

	return
}

// EndVoteV1原生sql优化
func EndVoteV1() {
	votes := make([]Vote, 0)

	if err := Conn.Raw("select * from vote where status = ? ", 1).Scan(&votes).Error; err != nil {
		return
	}

	now := time.Now().Unix()
	for _, vote := range votes {
		if vote.Time+vote.CreatedTime.Unix() <= now {
			Conn.Exec("update vote set status = 0 where id = ? limit 1", vote.Id)
		}
	}

	return
}

// GetVoteByName 根据名称检索投票。
func GetVoteByName(name string) *Vote {
	var ret Vote
	if err := Conn.Table("vote").Where("title = ?", name).First(&ret).Error; err != nil {
		fmt.Printf("")
	}
	return &ret
}
