package controller

import (
	"gopkg.in/telebot.v3"
	"strings"
)

func (c *Controller) handleDefault(uid int64, msg string) sender {
	if _, exists := c.quizKeeper.GetQuiz(msg); exists {
		c.quizKeeper.SetCurrentQuiz(uid, msg)
		c.switchState(uid, QuizChosen)
		return c.quizMessageSender(msg)
	}
	return c.sendDefaultGreetingsMessage
}

func (c *Controller) handleQuizChosen(uid int64, msg string) sender {
	quizName, exists := c.quizKeeper.GetCurrentQuiz(uid)
	if !exists || msg == backFromQuiz {
		c.quizKeeper.UnsetCurrentQuiz(uid)
		c.switchState(uid, DefaultState)
		return c.sendDefaultGreetingsMessage
	}
	if marathon := strings.HasPrefix(msg, startMarathonPrefix); marathon || strings.HasPrefix(msg, startQuizPrefix) {
		if quiz, exists := c.quizKeeper.GetQuiz(quizName); exists {
			c.switchState(uid, AttemptStarted)
			attempt := c.attempsKeeper.StartAttempt(uid, quiz, marathon)
			question, _ := attempt.CurrentQuestion()
			return c.attemptQuestionSender(question, attempt.CurrentQuestionNumber(), attempt.TotalQuestions())
		} else {
			c.quizKeeper.UnsetCurrentQuiz(uid)
			c.switchState(uid, DefaultState)
			return c.sendDefaultGreetingsMessage
		}
	}
	return c.sendHelpQuizChosen
}

func (c *Controller) handleAttempt(uid int64, msg string) sender {
	attempt := c.attempsKeeper.GetAttempt(uid)

	if msg == endAttempt {
		c.switchState(uid, QuizChosen)
		return func(ctx telebot.Context) error {
			correct, total := attempt.Results()
			return c.testResultsSender(uid, correct, total)(ctx)
		}
	}

	questionAnswered, _ := attempt.CurrentQuestion()
	res := attempt.SubmitAnswer(msg)
	if res == nil {
		return c.sendHelpInAttempt
	}

	sendResult := c.questionResultSender(questionAnswered, *res)
	var sendNext sender
	nextQuestion, exists := attempt.CurrentQuestion()
	if exists {
		sendNext = c.attemptQuestionSender(nextQuestion, attempt.CurrentQuestionNumber(), attempt.TotalQuestions())
	} else {
		c.switchState(uid, QuizChosen)
		correct, total := attempt.Results()
		sendNext = c.testResultsSender(uid, correct, total)
	}

	return func(ctx telebot.Context) error {
		err := sendResult(ctx)
		if err != nil {
			return err
		}
		return sendNext(ctx)
	}
}
