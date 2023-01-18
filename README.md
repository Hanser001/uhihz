# uhihz

## 👀简介

基于**gin**框架实现的婴儿级知乎后端

## 功能

### 用户

- 登录与注册
- 修改个人信息（如用户名，个性签名）
- 可以对文章、问题登进行点赞
- 查看自己的提问、回答、收藏

### 问答

- 能够发布问题并进行管理
- 可以发布并修改自己的回答
- 可以对评论进行回复和点赞

### 文章

- 可以发布文章并进行管理
- 可以回复和点赞评论

## 项目结构🎨

```
├── app ----------------------------- (项目文件)
    	├── api ------------------------- (api层)
    		├── article ------------------ (关于文章的api)
    		├── user ------------------(关于用户的api)
    		├── question -------------------- (关于提问的api)
    	├── global ---------------------- (全局组件)
    	├── internal -------------------- (内部包)
    		├── middleware -------------- (中间件)
    		├── model ------------------- (模型)
    		├── service ----------------- (服务层)
    	├── router ---------------------- (路由层)
    ├── boot ---------------------------- (项目启动包)
    ├── utils --------------------------- (工具包)
    	├── jwt ---------------------------(jwt使用)
```



## 🚲技术栈

![img](https://pics0.baidu.com/feed/37d12f2eb9389b500b4171da9ee664d4e5116edc.jpeg@f_auto?token=5d58f71b7b323fc31a5bf542a3dec2c7)

- [Gin](https://learnku.com/docs/gin-gonic/1.7)

> Gin是一个golang的微框架，封装比较优雅，API友好，源码注释比较明确，已经发布了1.0版本。具有快速灵活，容错方便等特点。

只会这个😂

![img](https://github.com/StellarisW/gohu/raw/master/manifest/image/mysql.svg)

- [MySQL](MySQLMySQLhttps://www.mysql.com/)

> 由瑞典MySQL AB 公司开发，属于 Oracle 旗下产品。MySQL 是最流行的关系型数据库管理系统关系型数据库管理系统之一，在 WEB 应用方面，MySQL是最好的 RDBMS (Relational Database Management System，关系数据库管理系统) 应用软件之一

![img](https://github.com/StellarisW/gohu/raw/master/manifest/image/redis.svg)

- [Redis](redisRedishttps://redis.io/)

> 使用C语言编写的、支持网络交互的、可基于内存也可持久化的Key-Value数据库

![img](https://camo.githubusercontent.com/dd51cf3dbd56f3c69f73f26255f377384d4dec4665d884a56ae1fd6a7bda319c/687474703a2f2f6a77742e696f2f696d672f6c6f676f2d61737365742e737667)

- JWT

- > SON Web Token (JWT)是一个开放标准(RFC 7519)，它定义了一种紧凑的、自包含的方式，用于作为JSON对象在各方之间安全地传输信息。该信息可以被验证和信任，因为它是数字签名的。

## 后端

### Token

- 使用 **"github.com/golang-jwt/jwt/v4"**包，用自定义声明和自定义密钥来为用户颁发Token

  ```
  type CustomClaims struct {
     //额外保存用户id
     Uid int
     jwt.RegisteredClaims
  }
  ```

- 后续的处理处理函数可以取出Token中存储的信息

### 问答

- 缓存问题相关数据

在redis存储问题的点击数、回答数、被点赞数，减少与mysql的交互

```
key := fmt.Sprintf("question:%s", strconv.Itoa(qid))
filed := "answerNum"

g.Redis.HIncrBy(ctx, key, filed, int64(incr))
```

### 文章

- 缓存相关数据

和问答一样，把一些会高频更新的数据存在redis里

### 评论和回复

- 模仿贴吧的评论模式，分为评论和回复，回复别人的评论时会显示被回复者id（用户名）

![image-20230118155011012](https://typora-1314425967.cos.ap-nanjing.myqcloud.com/typora/image-20230118155011012.png)

### 🐱‍👓存储设计

#### 表设计

- **user：记录用户信息**

![image-20230118155414461](https://typora-1314425967.cos.ap-nanjing.myqcloud.com/typora/image-20230118155414461.png)

- **article:记录文章信息**

![image-20230118155508093](https://typora-1314425967.cos.ap-nanjing.myqcloud.com/typora/image-20230118155508093.png)

- **question:记录提问信息**

![image-20230118155539658](https://typora-1314425967.cos.ap-nanjing.myqcloud.com/typora/image-20230118155539658.png)

- **article_comment:记录文章下评论信息**

![image-20230118155704527](https://typora-1314425967.cos.ap-nanjing.myqcloud.com/typora/image-20230118155704527.png)

- **answer_comment:记录提问下回答和评论信息**

![image-20230118155721309](https://typora-1314425967.cos.ap-nanjing.myqcloud.com/typora/image-20230118155721309.png)

- **article_collection:记录用户收藏文章**

![image-20230118155759755](https://typora-1314425967.cos.ap-nanjing.myqcloud.com/typora/image-20230118155759755.png)

> 用户id和文章id组成联合主键，避免重复收藏

#### 📮缓存设计

##### 提问缓存

![image-20230118205351081](https://typora-1314425967.cos.ap-nanjing.myqcloud.com/typora/image-20230118205351081.png)

用hash来存储高频更新的字段（点赞数，点击数，评论数）

##### 文章缓存

同提问缓存

### 👓TODO

- [x] 注册登录
- [x] 发布问题
- [x] 发布文章
- [x] 回答问题
- [x] 评论系统
- [x] 获取用户基本信息
- [x] 修改个人信息（密码，用户名）
- [x] 收藏文章
- [x] 获取自己的提问、文章、回答
- [x] 管理自己的提问、文章、回答
- [ ] 绑定邮箱
- [ ] 热榜
- [ ] 关注
- [ ] 通知
- [ ] 搜索

## 📕API文档

[乎知 - Powered by Apipost V7](https://console-docs.apipost.cn/preview/35aba4c302e7ba8e/95e8f9894abf426a?target_id=b8461525-1006-47c2-bf48-13e13985ea9a)
