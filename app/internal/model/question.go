package model

import "time"

type Question struct {
	Id         int
	Uid        int
	title      string
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

type Answer struct {
	Id         int
	Qid        int
	Uid        int
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}

type QuestionCollect struct {
	Uid        int
	Qid        int
	CreateTime time.Time
}
