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

// GetAuthorId 查询作者id
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

// ReadArticle 查看文章
func ReadArticle(id int) (model.Article, error) {
	sqlStr := "select * from article where id=?"

	stmt, err := g.Mysql.Prepare(sqlStr)

	var a model.Article

	if err != nil {
		g.Logger.Error(err.Error())
		return a, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&a.Id, &a.Uid, &a.Title, &a.Content, &a.CreateTime, &a.UpdateTime)
	if err != nil {
		g.Logger.Error(err.Error())
		return a, err
	}

	return a, nil
}

// AddArticleClick 增加文章点击量
func AddArticleClick(ctx context.Context, aid int) {
	key := fmt.Sprintf("article:%s", strconv.Itoa(aid))
	filed := "clickNum"

	g.Redis.HIncrBy(ctx, key, filed, 1)
}

// DeleteArticle 删除文章
func DeleteArticle(id int) error {
	//用事务进行删除
	tx, err := g.Mysql.Begin()
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	//既删除文章本身，还要删除附属于它的评论
	sqlStr1 := "delete from article where id=?"
	sqlStr2 := "delete from article_comment where aid=?"

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

// CollectArticle 对文章的收藏
func CollectArticle(uid, aid int) error {
	sqlStr := "insert into article_collection(uid,aid) values (?,?)"

	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(uid, aid)
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// CancelCollection 取消对文章收藏
func CancelCollection(aid, uid int) error {
	sqlStr := "delete from article_collection where aid=? and uid=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(aid, uid)
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// JudgeCollect 判断一篇文章是否已被用户收藏
func JudgeCollect(uid, aid int) (bool, error) {
	sqlStr := "select count(*) from article_collection where uid=? and aid=?"
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error(err.Error())
		return false, err
	}

	defer stmt.Close()

	var n int

	err = stmt.QueryRow(uid, aid).Scan(&n)
	if err != nil {
		g.Logger.Error(err.Error())
		return false, err
	}

	if n > 0 {
		return true, nil
	} else {
		return false, nil
	}
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

// ReadComment 获取评论内容
func ReadComment(id int) (error, model.ArticleComment) {
	sqlStr := "select * from article_comment where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err, model.ArticleComment{}
	}

	defer stmt.Close()

	var a model.ArticleComment

	err = stmt.QueryRow(id).Scan(&a.Id, &a.Aid, &a.Uid, &a.Pid, &a.ToUid, &a.Content, &a.CreateTime, &a.UpdateTime)
	if err != nil {
		g.Logger.Error(err.Error())
		return err, model.ArticleComment{}
	}

	return nil, a
}

// DeleteComment 删除评论
func DeleteComment(cid int) error {
	sqlStr := "delete from article_comment where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(cid)
	if err != nil {
		g.Logger.Error(err.Error())
		return err
	}

	return nil
}

// GetReviewerId 取得评论者id
func GetReviewerId(cid int) (error, int) {
	sqlStr := "select uid from article_comment where id=?"
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

// NewArticleCommentNum 更新文章评论数
func NewArticleCommentNum(ctx context.Context, aid int, incr int) {
	//发布评论或回复incr为1，删除时incr为-1
	key := fmt.Sprintf("article:%s", strconv.Itoa(aid))
	filed := "commentNum"

	g.Redis.HIncrBy(ctx, key, filed, int64(incr))
}

// LikeToArticle 点赞文章
func LikeToArticle(ctx context.Context, aid int, uid int) {
	key := fmt.Sprintf("articleLike:%s", strconv.Itoa(aid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	//把点赞用户id放入set
	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UnlikeToArticle 取消对文章点赞
func UnlikeToArticle(ctx context.Context, aid int, uid int) {
	key := fmt.Sprintf("articleLike:%s", strconv.Itoa(aid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	//把用户id从set中删除
	intCmd := g.Redis.SRem(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UpdateArticleLikeNum 更新文章点赞数
func UpdateArticleLikeNum(ctx context.Context, aid int, incr int) {
	key := fmt.Sprintf("article:%s", strconv.Itoa(aid))
	filed := "likeNum"

	//点赞时incr为1，取消点赞时incr为-1
	g.Redis.HIncrBy(ctx, key, filed, int64(incr))
}

// LikeArticleComment 点赞文章下评论（回复）
func LikeArticleComment(ctx context.Context, cid int, uid int) {
	key := fmt.Sprintf("commentLike:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	//把点赞用户id放入set
	intCmd := g.Redis.SAdd(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// UnlikeArticleComment 取消对评论（回复）点赞
func UnlikeArticleComment(ctx context.Context, cid int, uid int) {
	key := fmt.Sprintf("commentLike:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	//把用户id从set中删除
	intCmd := g.Redis.SRem(ctx, key, value)
	_, err := intCmd.Result()

	if err != nil {
		g.Logger.Error(err.Error())
		return
	}
}

// JudgeLikeToArticle  判断用户是否对某一文章点赞
func JudgeLikeToArticle(ctx context.Context, aid, uid int) bool {
	key := fmt.Sprintf("articleLike:%s", strconv.Itoa(aid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	strCmd := g.Redis.SIsMember(ctx, key, value)

	flag, _ := strCmd.Result()

	return flag
}

// JudgeLikeToComment 判断用户是否对某一评论（回复）点赞
func JudgeLikeToComment(ctx context.Context, cid int, uid int) bool {
	key := fmt.Sprintf("commentLike:%s", strconv.Itoa(cid))
	value := fmt.Sprintf("userLike:%s", strconv.Itoa(uid))

	strCmd := g.Redis.SIsMember(ctx, key, value)

	flag, _ := strCmd.Result()

	return flag
}

// UpdateArticleHot 更新文章的总热度值
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
