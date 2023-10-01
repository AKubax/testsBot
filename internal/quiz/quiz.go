package quiz

import "encoding/json"

type Quiz struct {
	QuestionsInAttempt int        `json:"questions_in_attempt"`
	Questions          []Question `json:"questions"`
}

type Question struct {
	Prompt string `json:"prompt"`
	Image  string `json:"image_name"`

	CorrectAnswer string   `json:"correct"`
	FakeAnswers   []string `json:"fakes"`
}

func (q *Question) UnmarshalJSON(bytes []byte) error {
	type question Question
	var data struct {
		question
		Fake1 string
		Fake2 string
	}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	*q = Question(data.question)
	if len(data.Fake1) > 1 {
		q.FakeAnswers = append(q.FakeAnswers, data.Fake1)
	}
	if len(data.Fake2) > 1 {
		q.FakeAnswers = append(q.FakeAnswers, data.Fake2)
	}
	return nil
}
