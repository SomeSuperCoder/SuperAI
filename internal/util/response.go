package util

type Response struct {
	Choices []choice
}

type choice struct {
	Delta delta
}

type delta struct {
	Content string
}
