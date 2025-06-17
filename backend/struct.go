package main

type Score struct {
	Id int	`json:"Id"`
	Name string `json:"Name"`
	Wpm int `json:"Wpm"`
	Errors int `json:"Errors"`
	Characters int `json:"Characters"`
	Accuracy float32 `json:"Accuracy"`
	WordCount int `json:"WordCount"`
	Time int  `json:"Time"`

}

