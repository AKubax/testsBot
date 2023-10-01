package controller

type UserState byte

const (
	DefaultState UserState = iota
	QuizChosen
	AttemptStarted
)
