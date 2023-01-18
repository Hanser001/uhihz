package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	g "zhihu/app/global"
	"zhihu/app/internal/model"
)

// PublicQuestion 发布问题
func PublicQuestion(uid int, title string, content string) error {
	sqlStr := "insert into question(uid,title,content) values (?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(uid, title, content)
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}
	return nil
}

// UpdateQuestion 更新问题描述
func UpdateQuestion(newContent string, id int) {
	sqlStr := "update question set content=? where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}

	defer stmt.Close()

	stmt.Exec(newContent, id)
}

// GetQuestionerId 取得提问者id
func GetQuestionerId(id int) int {
	sqlStr := "select uid from question where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return 0
	}

	defer stmt.Close()

	var q model.Question
	stmt.QueryRow(id).Scan(&q.Uid)

	questionerId := q.Uid
	return questionerId
}

// ReadQuestion 查看问题详细内容
func ReadQuestion(id int) model.Question {
	sqlStr := "select * from question where id=?"

	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
	}

	defer stmt.Close()

	var q model.Question
	stmt.QueryRow(id).Scan(&q.Id, &q.Uid, &q.Title, &q.Content, &q.CreateTime, &q.UpdateTime)
	return q
}

// AddQuestionClick 增加问题点击量
func AddQuestionClick(ctx context.Context, qid int) {
	//将问题的点击量，点赞量都放在一个hash里
	key := fmt.Sprintf("question:%s", strconv.Itoa(qid))
	filed := "clickNum"

	g.Redis.HIncrBy(ctx, key, filed, 1)
}

// DeleteQuestion 删除问题
func DeleteQuestion(id int) error {
	//用一个事务进行删除
	tx, err := g.Mysql.Begin()
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	//既删除问题本身，还要删除附属于它的回答和评论
	sqlStr1 := "delete from question where id=?"
	sqlStr2 := "delete from answer_comment where qid=?"

	_, err = tx.Exec(sqlStr1, id)
	if err != nil {
		tx.Rollback()
		g.Logger.Error(err.Error())
		return err
	}

	_, err = tx.Exec(sqlStr2, id)
	if err != nil {
		tx.Rollback()
		g.Logger.Error(err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// PublicAnswer 发布对问题的回答
func PublicAnswer(qid, uid int, content string) error {
	sqlStr := "insert into answer_comment(qid,uid,content) values (?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(qid, uid, content)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}
	return nil
}

// UpdateAnswer 更新自己的回答
func UpdateAnswer(content string, id int) error {
	sqlStr := "update answer_comment set content=? where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(content, id)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// CommentTheAnswer 评论对问题的回答
func CommentTheAnswer(qid, uid, pid int, content string) error {
	sqlStr := "insert into answer_comment(qid,uid,pid,content) valuse(?,?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(qid, uid, pid, content)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}
	return nil
}

// ReplyToComment 对评论进行回复
func ReplyToComment(qid, uid, pid, toUid int, content string) error {
	sqlStr := "insert into answer_comment(qid,uid,pid,toUid,content) values (?,?,?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(qid, uid, pid, toUid, content)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}
	return nil
}

// ReadReview 取得回复详细内容
func ReadReview(id int) (error, model.AnswerComment) {
	sqlStr := "select * from answer_comment where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err, model.AnswerComment{}
	}

	var a model.AnswerComment
	err = stmt.QueryRow(id).Scan(&a.Id, &a.Qid, &a.Uid, &a.Pid, &a.ToUid, &a.Content, &a.CreateTime, &a.UpdateTime)

	if err != nil {
		g.Logger.Error(err.Error())
		return err, model.AnswerComment{}
	}

	return nil, a
}

// DeleteReview 删除回复
func DeleteReview(cid int) error {
	sqlStr := "delete from answer_comment where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	_, err = stmt.Exec(cid)
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// GetAnswererId 取得答者id
func GetAnswererId(cid int) (error, int) {
	sqlStr := "select uid,pid from answer_comment where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err, 0
	}

	defer stmt.Close()

	var u model.User
	err = stmt.QueryRow(cid).Scan(&u.Id)
	if err != nil {
		g.Logger.Error(err.Error())
		return err, 0
	}

	return nil, u.Id
}

// UpdateQuestionAnswerNum 更新问题回答数
func UpdateQuestionAnswerNum(ctx context.Context, qid int, incr int) {
	//用哈希表存放问题的点击量，回答数，点赞量
	key := fmt.Sprintf("question:%s", strconv.Itoa(qid))
	filed := "answerNum"

	g.Redis.HIncrBy(ctx, key, filed, int64(incr))
}

// LikeToQuestion 点赞问题
func LikeToQuestion(ctx context.Context, qid int, uid int) {
	key := fmt.Sprintf("questionLike:%s", strconv.Itoa(qid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	//把点赞用户id放入一个set
	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UnlikeToQuestion 取消对问题点赞
func UnlikeToQuestion(ctx context.Context, qid int, uid int) {
	key := fmt.Sprintf("questionLike:%s", strconv.Itoa(qid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	//把点赞用户id从set中删除
	intCmd := g.Redis.SRem(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UpdateQuestionLikeNum 更新问题点赞数
func UpdateQuestionLikeNum(ctx context.Context, qid int, incr int) {
	key := fmt.Sprintf("question:%s", strconv.Itoa(qid))
	filed := "likeNum"

	//点赞时incr为1，取消点赞时incr为-1
	g.Redis.HIncrBy(ctx, key, filed, int64(incr))
}

// LikeToAnswer 点赞回答(评论/回复)
func LikeToAnswer(ctx context.Context, cid int, uid int) {
	key := fmt.Sprintf("answerLike:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UnlikeToAnswer 取消对回答(评论/回复)点赞
func UnlikeToAnswer(ctx context.Context, cid int, uid int) {
	key := fmt.Sprintf("answerLike:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SRem(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// JudgeLikeToQuestion 判断用户是否对某一问题点赞
func JudgeLikeToQuestion(ctx context.Context, qid int, uid int) bool {
	key := fmt.Sprintf("questionLike:%s", strconv.Itoa(qid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	strCmd := g.Redis.SIsMember(ctx, key, value)

	flag, _ := strCmd.Result()

	return flag
}

// JudgeLikeToAnswer 判断用户是否对某一回答（评论)点赞
func JudgeLikeToAnswer(ctx context.Context, cid int, uid int) bool {
	key := fmt.Sprintf("answerLike:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	strCmd := g.Redis.SIsMember(ctx, key, value)

	flag, _ := strCmd.Result()

	return flag
}

// UpdateQuestionHot 更新问题热度值
func UpdateQuestionHot(ctx context.Context, qid int, hot int) {
	key := "questionHot"
	value := &redis.Z{
		Score:  float64(hot),
		Member: qid,
	}

	_, err := g.Redis.ZAdd(ctx, key, value).Result()

	if err != nil {
		g.Logger.Error(err.Error())
	}
}
