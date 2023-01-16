package service

import (
	g "zhihu/app/global"
	"zhihu/app/internal/model"
)

// PublicQuestion 发布问题
func PublicQuestion(uid int, title string, content string) {
	sqlStr := "insert into question(uid,title,content) values(?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}

	defer stmt.Close()

	stmt.Exec(uid, title, content)
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
	stmt.QueryRow(id).Scan(&q.Uid, &q.CreateTime, &q.UpdateTime, &q.Content)
	return q
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
func PublicAnswer(qid, uid int, content string) {
	sqlStr := "insert into answer_comment(qid,uid,content) values (?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	stmt.Exec(qid, uid, content)
}

// CommentTheAnswer 评论对问题的回答
func CommentTheAnswer(qid, uid, pid int, content string) {
	sqlStr := "insert into answer_comment(qid,uid,pid,content) valuse(?,?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	stmt.Exec(qid, uid, pid, content)
}

// ReplyToComment 对评论进行回复
func ReplyToComment(qid, uid, pid, toUid int, content string) {
	sqlStr := "insert into answer_comment(qid,uid,pid,toUid,content) values (?,?,?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	stmt.Exec(qid, uid, pid, toUid, content)
}
