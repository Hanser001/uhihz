package model

import "time"

type Question struct {
	Id         int
	Uid        int
	Title      string
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

type AnswerComment struct {
	Id         int
	Qid        int
	Uid        int
	Pid        int
	ToUid      *int
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

type QuestionCollect struct {
	Uid        int
	Qid        int
	CreateTime time.Time
}
