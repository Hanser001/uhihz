package service

import (
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	g "zhihu/app/global"
	"zhihu/app/internal/model"
)

// SelectUserExists 注册时检索用户名是否已存在
func SelectUserExists(username string) error {
	sqlStr := "select username from user where id>?"
	//进行sql预处理，提高效率并防止sql注入
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error("prepare failed")
		return fmt.Errorf("internal error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(0)

	if err != nil {
		g.Logger.Error("query mysql record failed", zap.Error(err))
		return fmt.Errorf("internal error")
	}

	// 关闭rows释放持有的数据库链接
	defer rows.Close()

	for rows.Next() {
		var u model.User
		rows.Scan(&u.Username)
		if username == u.Username {
			return fmt.Errorf("user already exists")
		}
	}

	return nil
}

// EncryptPasswordWithSalt 加密密码
func EncryptPasswordWithSalt(password string) string {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		g.Logger.Error("encrypt failed")
		return "" //加密失败则返回空字符串
	} else {
		//加密成功返回加密后的字符串和空异常
		return string(encryptedPassword)
	}
}

// CheckUsername 登录时检查用户是否存在
func CheckUsername(username string) bool {
	sqlStr := "select username from user where id>?"
	//进行sql预处理，提高效率并防止sql注入
	stmt, _ := g.Mysql.Prepare(sqlStr)

	defer stmt.Close()

	rows, _ := stmt.Query(0)

	defer rows.Close()

	for rows.Next() {
		var u model.User
		rows.Scan(&u.Username)
		if username == u.Username {
			return true
		}
	}

	return false
}

// GetPassword 登录时从数据库取出用户密码
func GetPassword(username string) string {
	sqlStr := "select password from user where username=? "

	stmt, _ := g.Mysql.Prepare(sqlStr)

	defer stmt.Close()

	var u model.User
	stmt.QueryRow(username).Scan(&u.Password)

	return u.Password
}

// CheckPassword 登录时将用户输入密码与数据库中密码进行对比
func CheckPassword(password, encryptedPassword string) bool {
	//使用bcrypt包的CompareHashAndPassword对比密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	//对比密码是否正确回返回一个异常，第一个参数为加密后的密码，第二个参数为未加密的密码，按照官方说法，只要异常是nil，则密码正确
	return err == nil
}

// AddUser 添加用户
func AddUser(username, password string) {
	sqlStr := "insert into user(username,password) values (?,?) "
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}
	defer stmt.Close()

	stmt.Exec(username, password)
}

// SetPersonalSign 设置个性签名
func SetPersonalSign(personalSignature string, id int) {
	sqlStr := "update user set personalSignature=? where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}
	defer stmt.Close()

	stmt.Exec(personalSignature, id)
}

// ChangePassword 更改密码
func ChangePassword(newPassword string, username string) {
	sqlStr := "update user set password=? where username=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}

	defer stmt.Close()

	stmt.Exec(newPassword, username)
}

// ChangeUsername 更改用户名
func ChangeUsername(newUsername string, id int) {
	sqlStr := "update user set username=? where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
		return
	}

	defer stmt.Close()

	stmt.Exec(newUsername, id)
}

// GetUserInfo 查询单个用户信息
func GetUserInfo(id int) model.User {
	sqlStr := "select id,username,createTime,personalSignature from user where id=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
	}

	defer stmt.Close()
	var u model.User

	stmt.QueryRow(id).Scan(&u.Id, &u.Username, &u.CreateTime, &u.PersonalSignature)

	return u
}

// GetUid 获取用户id用于颁发token
func GetUid(username string) int {
	sqlStr := "select id from user where username=?"
	stmt, err := g.Mysql.Prepare(sqlStr)

	if err != nil {
		g.Logger.Error("prepare failed")
	}

	defer stmt.Close()

	var u model.User

	stmt.QueryRow(username).Scan(&u.Id)

	return u.Id
}

// GetUserArticles 取得用户发布的文章
func GetUserArticles(uid int) ([]model.Article, error) {
	sqlStr := "select * from article where uid=?"
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error(err.Error())
		return []model.Article{}, err
	}

	defer stmt.Close()

	var Articles []model.Article

	rows, err := stmt.Query(uid)

	defer rows.Close()

	for rows.Next() {
		var a model.Article
		err := rows.Scan(&a.Id, &a.Uid, &a.Title, &a.Content, &a.CreateTime, &a.UpdateTime)

		if err != nil {
			g.Logger.Error(err.Error())
			return []model.Article{}, err
		}

		Articles = append(Articles, a)
	}

	return Articles, nil
}

// GetUserArticleCollection 获取用户文章收藏表
func GetUserArticleCollection(uid int) ([]model.Article, error) {
	sqlStr := "select * from article where id = (select aid from article_collection where uid=?)"
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error(err.Error())
		return []model.Article{}, err
	}

	defer stmt.Close()

	var ArticleCollections []model.Article

	rows, err := stmt.Query(uid)
	if err != nil {
		g.Logger.Error(err.Error())
		return []model.Article{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Article
		err := rows.Scan(&a.Id, &a.Uid, &a.Title, &a.Content, &a.CreateTime, &a.UpdateTime)

		if err != nil {
			g.Logger.Error(err.Error())
			return []model.Article{}, err
		}

		ArticleCollections = append(ArticleCollections, a)
	}

	return ArticleCollections, nil
}

// GetUserAnswer 取得用户的回答
func GetUserAnswer(uid int) ([]model.AnswerComment, error) {
	sqlStr := "select qid,content,createTime,updateTime from answer_comment where uid=? and pid is null"
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error(err.Error())
		return []model.AnswerComment{}, err
	}

	defer stmt.Close()

	var answers []model.AnswerComment

	rows, err := stmt.Query(uid)
	if err != nil {
		g.Logger.Error(err.Error())
		return []model.AnswerComment{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.AnswerComment
		err := rows.Scan(&a.Qid, &a.Content, &a.Content, &a.UpdateTime)

		if err != nil {
			g.Logger.Error(err.Error())
			return []model.AnswerComment{}, err
		}

		answers = append(answers, a)
	}

	return answers, nil
}

// GetUserQuestions 取得用户所有提问
func GetUserQuestions(uid int) ([]model.Question, error) {
	sqlStr := "select title,content,createTime,updateTime from question where uid=?"
	stmt, err := g.Mysql.Prepare(sqlStr)
	if err != nil {
		g.Logger.Error(err.Error())
		return []model.Question{}, err

	}

	defer stmt.Close()

	var questions []model.Question

	rows, err := stmt.Query(uid)
	if err != nil {
		g.Logger.Error(err.Error())
		return []model.Question{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var q model.Question
		err := rows.Scan(&q.Title, &q.Content, &q.CreateTime, &q.UpdateTime)

		if err != nil {
			g.Logger.Error(err.Error())
			return []model.Question{}, err
		}

		questions = append(questions, q)
	}

	return questions, nil
}
