package main

import (
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"testsBot/internal/attempt"
	"testsBot/internal/controller"
	"testsBot/internal/quiz"
	"time"
)

func main() {
	token, err := os.ReadFile("secrets/tg_token")
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:   string(token),
		Updates: 0,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	})

	quizKeeper := quiz.NewKeeper()
	err = quizKeeper.LoadQuizzesFromFiles("data/quizzes.json")
	if err != nil {
		log.Fatal(err)
	}

	controller := controller.NewController(attempt.NewKeeper(), quizKeeper)

	bot.Handle(telebot.OnText, controller.HandleMessage)
	bot.Start()
}
