package model

import (
	"time"
)

type Article struct {
	Id         int
	Uid        int
	Title      string
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

type ArticleComment struct {
	Id         int
	Aid        int
	Uid        int
	Pid        *int
	ToUid      *int
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

type ArticleCollection struct {
	Uid        int
	Aid        int
	CreateTime time.Time
}
