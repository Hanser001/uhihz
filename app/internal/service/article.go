package service

import (
	g "zhihu/app/global"
	"zhihu/app/internal/model"
)

// PublishArticle 发布文章
func PublishArticle(uid int, title string, content string) {
	sqlStr := "insert into article(uid,title,content) values (?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}
	defer stmt.Close()

	stmt.Exec(uid, title, content)
}

// UpdateArticle 更新文章
func UpdateArticle(content string, id int) {
	sqlStr := "update article set content=? where id=?"

	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}
	defer stmt.Close()

	stmt.Exec(content, id)
}

// SelectArticle 查看文章(该返回什么数据?)
func SelectArticle(id int) model.Article {
	sqlStr := "select * from article where id=?"

	stmt, err := g.Mysql.Prepare(sqlStr)

	var a model.Article

	if err != nil {
		g.Logger.Error(err.Error())
		return a
	}

	defer stmt.Close()

	stmt.QueryRow(id).Scan(&a.Id, &a.Uid, &a.CreateTime, &a.UpdateTime, &a.Title, &a.Content)

	return a
}

// DeleteArticle 删除文章
func DeleteArticle(id int) {
	sqlStr := "delete from article where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	stmt.Exec(id)
}

// CollectArticle 对文章的收藏
func CollectArticle(uid, aid int) {
	sqlStr := "insert into article_collection(uid,aid) values (?,?)"

	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}

	defer stmt.Close()

	stmt.Exec(uid, aid)
}

// GetUserArticleCollection 获取用户文章收藏表
func GetUserArticleCollection(uid int) []int {
	sqlStr := "select * from artilce_collection where uid=?"
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	var ArticleFavorites []int

	rows, err := stmt.Query(uid)
	for rows.Next() {
		var q model.QuestionCollect
		err := rows.Scan(&q.Qid)

		if err != nil {
			g.Logger.Error(err.Error())
		}

		ArticleFavorites = append(ArticleFavorites, q.Qid)
	}

	return ArticleFavorites
}

// PublicComment 发表对文章评论
func PublicComment(aid, uid int, content string) {
	sqlStr := "insert into article_comment(aid,uid,cotent) values (?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	stmt.Exec(aid, uid, content)
}

// CommentTheComment 对评论(一级评论)进行评论(二级评论)
func CommentTheComment(aid, uid, pid int, content string) {
	sqlStr := "insert into article_comment(aid,uid,pid,content) values (?,?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	stmt.Exec(aid, uid, pid, content)
}

// ReplyTheComment 对评论进行回复
func ReplyTheComment(aid, uid, pid, toUid int, content string) {
	sqlStr := "insert into article_comment(aid,uid,pid,toUid,content) values (?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
	}

	defer stmt.Close()

	stmt.Exec(aid, uid, pid, toUid, content)
}

// 点赞评论
func AddLike() {

}
