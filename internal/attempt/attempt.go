package attempt

import (
	"math/rand"
	"strings"
	"testsBot/internal/quiz"
)

func NewAttempt(listOfQuestions []quiz.Question, numberOfQuestions int, onlyAnsweredInResult bool) *Attempt {
	return &Attempt{
		results:              make([]bool, numberOfQuestions),
		questions:            questionsForAttempt(listOfQuestions, numberOfQuestions),
		onlyAnsweredInResult: onlyAnsweredInResult,
	}
}

func questionsForAttempt(listOfQuestions []quiz.Question, numberOfQuestions int) []quiz.Question {
	nums := rand.Perm(len(listOfQuestions))
	attemptQuestions := make([]quiz.Question, 0, numberOfQuestions)
	for _, num := range nums {
		if len(attemptQuestions) == cap(attemptQuestions) {
			break
		}
		attemptQuestions = append(attemptQuestions, listOfQuestions[num])
	}
	return attemptQuestions
}

type Attempt struct {
	results              []bool
	questions            []quiz.Question
	index                int
	onlyAnsweredInResult bool
}

func (a *Attempt) CurrentQuestionNumber() int {
	return a.index + 1
}

func (a *Attempt) TotalQuestions() int {
	return len(a.questions)
}

func (a *Attempt) CurrentQuestion() (quiz.Question, bool) {
	if a.index < len(a.questions) {
		return a.questions[a.index], true
	}
	return quiz.Question{}, false
}

func (a *Attempt) SubmitAnswer(answer string) *bool {
	question := a.questions[a.index]
	if buttonEquals(answer, question.CorrectAnswer) {
		return a.processAnswerResult(true)
	}
	for _, fake := range question.FakeAnswers {
		if buttonEquals(answer, fake) {
			return a.processAnswerResult(false)
		}
	}
	return nil
}

func (a *Attempt) processAnswerResult(isCorrect bool) *bool {
	a.results[a.index] = isCorrect
	a.index++
	return &isCorrect
}

func (a *Attempt) Results() (int, int) {
	result := 0
	for _, correct := range a.results {
		if correct {
			result += 1
		}
	}
	total := a.TotalQuestions()
	if a.onlyAnsweredInResult {
		total = a.index
	}
	return result, total
}

func buttonEquals(buttonText, text string) bool {
	if len(buttonText) > 100 && strings.HasPrefix(text, buttonText[:len(buttonText)-3]) {
		return true
	}

	return buttonText == text
}
