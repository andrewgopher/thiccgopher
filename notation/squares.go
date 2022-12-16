package notation

func FileToInt(file byte) int {
	return int(file - 'a')
}

func RankToInt(rank byte) int {
	return int(rank - '1')
}
