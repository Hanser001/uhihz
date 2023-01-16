package hot

// ComputeHot 一个简陋的热度算法
func ComputeHot(likeNumber, collectNumber, readNumber int) int {
	//一个阅读一分，一个点赞五分，一个收藏十分
	readMark := 1 * readNumber
	likeMark := 5 * likeNumber
	collectMark := 10 * collectNumber

	Hot := readMark + likeMark + collectMark

	return Hot
}
