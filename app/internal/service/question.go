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

// CollectQuestion 对问题的收藏
func CollectQuestion(uid, qid int) {
	sqlStr := "insert into question_collection(uid,qid) values (?,?)"

	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}

	defer stmt.Close()

	stmt.Exec(uid, qid)
}

// SelectQuestion 查看问题详细内容
func SelectQuestion(id int) model.Question {
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
	key := fmt.Sprintf("questionClick:%s", strconv.Itoa(qid))

	g.Redis.Incr(ctx, key)
}

// DeleteQuestion 删除问题
func DeleteQuestion(id int) {
	sqlStr := "delete from question where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	stmt.Exec(id)
}

// GetUserQuestionCollection 获取问题收藏表
func GetUserQuestionCollection(uid int) []int {
	sqlStr := "select * from article_collection where uid=?"
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	var QuestionFavorites []int

	rows, err := stmt.Query(uid)
	for rows.Next() {
		var a model.ArticleCollection
		err := rows.Scan(&a.Aid)
		if err != nil {
			g.Logger.Error(err.Error())
		}
		QuestionFavorites = append(QuestionFavorites, a.Aid)
	}

	return QuestionFavorites
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

// LikeToQuestion 点赞问题
func LikeToQuestion(ctx context.Context, qid int, uid int) {
	//以question:id为key,userLike:id为value，将所以用户id存入一个set
	key := fmt.Sprintf("question:%s", strconv.Itoa(qid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// LikeToAnswer 点赞回答
func LikeToAnswer(ctx context.Context, aid int, uid int) {
	key := fmt.Sprintf("answer:%s", strconv.Itoa(aid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// LikeAnswerComment 点赞评论
func LikeAnswerComment(ctx context.Context, cid int, uid int) {
	key := fmt.Sprintf("answerComment:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
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
