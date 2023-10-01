package controller

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"math/rand"
	"testsBot/internal/quiz"
)

const (
	startQuizPrefix     = "Начать тест: "
	startMarathonPrefix = "Начать марафон: "
	backFromQuiz        = "Назад"

	endAttempt = "Закончить тест"
)

type sender func(ctx telebot.Context) error

func (c *Controller) sendDefaultGreetingsMessage(ctx telebot.Context) error {
	return ctx.Send("<em>Выберите тест:</em>", telebot.ModeHTML, &telebot.ReplyMarkup{
		ReplyKeyboard: c.getMappedQuizesButtons(),
	})
}

func (c *Controller) getMappedQuizesButtons() [][]telebot.ReplyButton {
	quizes := c.quizKeeper.GetAllQuizzes()
	res := make([][]telebot.ReplyButton, 0, (len(quizes)+2)/3)
	var row []telebot.ReplyButton = make([]telebot.ReplyButton, 0, min(len(quizes)-len(res)*3, 3))
	for name := range quizes {
		row = append(row, telebot.ReplyButton{
			Text: name,
		})
		if len(row) == cap(row) {
			res = append(res, row)
			if len(quizes)-len(res)*3 > 0 {
				row = make([]telebot.ReplyButton, 0, min(len(quizes)-len(res)*3, 3))
			}
		}
	}

	return res
}

func (c *Controller) quizMessageSender(quizName string) sender {
	return func(ctx telebot.Context) error {
		return ctx.Send(fmt.Sprintf("<em>Выбран тест:</em> %s", quizName), telebot.ModeHTML, &telebot.ReplyMarkup{
			ReplyKeyboard: c.getQuizMessageButtons(quizName),
		})
	}
}

func (c *Controller) getQuizMessageButtons(quizName string) [][]telebot.ReplyButton {
	return [][]telebot.ReplyButton{
		{
			{
				Text: fmt.Sprintf("%s%s", startQuizPrefix, quizName),
			},
			{
				Text: fmt.Sprintf("%s%s", startMarathonPrefix, quizName),
			},
		},
		{
			{
				Text: backFromQuiz,
			},
		},
	}
}

func (*Controller) sendHelpQuizChosen(ctx telebot.Context) error {
	return ctx.Send("Некорректный ввод. Используйте кнопки, чтобы начать тест или вернуться назад")
}

func (c *Controller) attemptQuestionSender(question quiz.Question, current, total int) sender {
	answerButtons := c.getAnswersButtons(question)
	msg := fmt.Sprintf("<u>Вопрос %d/%d</u>\n\n<em>%s:</em>\n", current, total, question.Prompt)
	for i, button := range answerButtons {
		msg = fmt.Sprintf("%s\n%d. %s", msg, i+1, button.Text)
	}

	var toSend any = msg
	if question.Image != "" {
		toSend = &telebot.Photo{
			File:    telebot.FromDisk(fmt.Sprintf("data/images/%s.png", question.Image)),
			Caption: msg,
		}
	}

	return func(ctx telebot.Context) error {
		return ctx.Send(toSend, telebot.ModeHTML, &telebot.ReplyMarkup{
			ReplyKeyboard: [][]telebot.ReplyButton{
				answerButtons,
				{
					{
						Text: endAttempt,
					},
				},
			},
		})
	}
}

func (*Controller) getAnswersButtons(question quiz.Question) []telebot.ReplyButton {
	answers := make([]string, 0, 3)
	answers = append(answers, question.CorrectAnswer)
	for _, num := range rand.Perm(len(question.FakeAnswers)) {
		answers = append(answers, question.FakeAnswers[num])
		if len(answers) == cap(answers) {
			break
		}
	}

	rand.Shuffle(len(answers), func(i, j int) {
		answers[i], answers[j] = answers[j], answers[i]
	})

	buttons := make([]telebot.ReplyButton, 0, len(answers))
	for _, answer := range answers {
		buttons = append(buttons, telebot.ReplyButton{Text: answer})
	}
	return buttons
}

func (c *Controller) testResultsSender(uid int64, correct, total int) sender {
	return func(ctx telebot.Context) error {
		quizName, _ := c.quizKeeper.GetCurrentQuiz(uid)
		return ctx.Send(
			fmt.Sprintf("<em>Тест закончен.</em> Результат: <strong>%d</strong>/<strong>%d</strong>", correct, total),
			telebot.ModeHTML,
			&telebot.ReplyMarkup{
				ReplyKeyboard: c.getQuizMessageButtons(quizName),
			},
		)
	}
}

func (*Controller) sendHelpInAttempt(ctx telebot.Context) error {
	return ctx.Send("Некорректный ввод. Используйте кнопки, чтобы ответить на вопрос или закончить тест")
}

func (*Controller) questionResultSender(question quiz.Question, correct bool) sender {
	if correct {
		return func(ctx telebot.Context) error {
			return ctx.Send("✅<strong>Верно!</strong>", telebot.ModeHTML)
		}
	}

	return func(ctx telebot.Context) error {
		return ctx.Send(fmt.Sprintf("❌<strong>Неверно</strong>\n\n<em>Правильный ответ:</em> \"%s\"", question.CorrectAnswer), telebot.ModeHTML)
	}
}
