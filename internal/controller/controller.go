package controller

import (
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"testsBot/internal/attempt"
	"testsBot/internal/quiz"
)

func NewController(attempsKeeper *attempt.Keeper, quizKeeper *quiz.Keeper) *Controller {
	return &Controller{
		states:        map[int64]UserState{},
		attempsKeeper: attempsKeeper,
		quizKeeper:    quizKeeper,

		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

type Controller struct {
	states        map[int64]UserState
	attempsKeeper *attempt.Keeper
	quizKeeper    *quiz.Keeper

	logger *log.Logger
}

func (c *Controller) HandleMessage(ctx telebot.Context) error {
	uid := ctx.Message().Sender.ID
	msg := ctx.Message().Text
	c.logger.Printf("[INFO][uid=%d][msg='%s'] Got message, current state = %d", uid, msg, c.getState(uid))

	var sendResponse sender
	switch c.getState(uid) {
	case AttemptStarted:
		sendResponse = c.handleAttempt(uid, msg)
	case QuizChosen:
		sendResponse = c.handleQuizChosen(uid, msg)
	default:
		sendResponse = c.handleDefault(uid, msg)
	}

	err := sendResponse(ctx)
	if err != nil {
		c.logger.Printf("[ERROR][uid=%d][msg='%s'] %s", uid, msg, err.Error())
	}
	c.logger.Printf("[INFO][uid=%d][msg='%s'] Handled message, current state = %d", uid, msg, c.getState(uid))
	return nil
}

func (c *Controller) getState(uid int64) UserState {
	return c.states[uid]
}

func (c *Controller) switchState(uid int64, target UserState) {
	c.states[uid] = target
}
