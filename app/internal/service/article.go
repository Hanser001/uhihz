package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	g "zhihu/app/global"
	"zhihu/app/internal/model"
)

// PublicArticle 发布文章
func PublicArticle(uid int, title string, content string) {
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

// GetAuthorId 查询文章作者id
func GetAuthorId(id int) int {
	sqlStr := "select uid from article where id=?"

	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return 0
	}
	defer stmt.Close()

	var a model.Article

	stmt.QueryRow(id).Scan(&a.Uid)

	authorId := a.Uid

	return authorId
}

// SelectArticle 查看文章
func SelectArticle(id int) model.Article {
	sqlStr := "select * from article where id=?"

	stmt, err := g.Mysql.Prepare(sqlStr)

	var a model.Article

	if err != nil {
		g.Logger.Error(err.Error())
		return a
	}

	defer stmt.Close()

	stmt.QueryRow(id).Scan(&a.Id, &a.Uid, &a.Title, &a.Content, &a.CreateTime, &a.UpdateTime)

	return a
}

// AddArticleClick 增加文章点击量
func AddArticleClick(ctx context.Context, aid int) {
	key := fmt.Sprintf("articleClick:%s", strconv.Itoa(aid))

	g.Redis.Incr(ctx, key)
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
func PublicComment(aid, uid int, content string) error {
	sqlStr := "insert into article_comment(aid,uid,content) values (?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(aid, uid, content)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// CommentTheComment 对评论(一级评论)进行评论(二级评论)
func CommentTheComment(aid, uid, pid int, content string) error {
	sqlStr := "insert into article_comment(aid,uid,pid,content) values (?,?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(aid, uid, pid, content)
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// ReplyTheComment 对评论进行回复
func ReplyTheComment(aid, uid, pid, toUid int, content string) error {
	sqlStr := "insert into article_comment(aid,uid,pid,toUid,content) values (?,?,?,?,?)"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(aid, uid, pid, toUid, content)
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// AddArticleCommentNum 增加文章评论数
func AddArticleCommentNum(ctx context.Context, aid int) {
	key := fmt.Sprintf("articleCommentNum:%s", strconv.Itoa(aid))

	g.Redis.Incr(ctx, key)
}

// LikeToArticle 点赞文章
func LikeToArticle(aid int, uid int, ctx context.Context) {
	//以article:id为key，userLike:id为value，将所有用户id存入一个set
	key := fmt.Sprintf("article:%s", strconv.Itoa(aid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UnlikeToArticle 取消对文章点赞
func UnlikeToArticle(aid int, uid int, ctx context.Context) {
	//以article:id为key，userLike:id为value，将所有用户id存入一个set
	key := fmt.Sprintf("article:%s", strconv.Itoa(aid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SRem(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// AddArticleLikeNum 增加文章点赞数
func AddArticleLikeNum(ctx context.Context, aid int) {
	key := fmt.Sprintf("articleLikeNum:%s", strconv.Itoa(aid))

	g.Redis.Incr(ctx, key)
}

// LikeArticleComment 点赞文章下评论（回复）
func LikeArticleComment(ctx context.Context, cid int, uid int) {
	key := fmt.Sprintf("articleComment:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UnlikeArticleComment 取消对评论（回复）点赞
func UnlikeArticleComment(ctx context.Context, cid int, uid int) {
	key := fmt.Sprintf("articleComment:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	intCmd := g.Redis.SRem(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UpdateArticleHot 更新文章的热度值
func UpdateArticleHot(c context.Context, aid int, hot int) {
	//用Zset来保存文章id及其热度
	key := "articleHot"
	value := &redis.Z{
		Score:  float64(hot),
		Member: aid,
	}

	g.Redis.ZAdd(c, key, value)

	_, err := g.Redis.ZAdd(c, key, value).Result()

	if err != nil {
		g.Logger.Error(err.Error())
	}
}
