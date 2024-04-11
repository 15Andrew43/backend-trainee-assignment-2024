package util

var startId = 100500

func GenerateNextId() int {
	startId++
	return startId
}
