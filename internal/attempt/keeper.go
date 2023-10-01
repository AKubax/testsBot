package attempt

import "testsBot/internal/quiz"

func NewKeeper() *Keeper {
	return &Keeper{currentAttemps: make(map[int64]*Attempt)}
}

type Keeper struct {
	currentAttemps map[int64]*Attempt
}

func (k *Keeper) StartAttempt(uid int64, quiz quiz.Quiz, marathon bool) *Attempt {
	numberOfQuestions := quiz.QuestionsInAttempt
	if marathon {
		numberOfQuestions = len(quiz.Questions)
	}

	attempt := NewAttempt(quiz.Questions, numberOfQuestions, marathon)
	k.currentAttemps[uid] = attempt
	return attempt
}

func (k *Keeper) GetAttempt(uid int64) *Attempt {
	return k.currentAttemps[uid]
}

func (k *Keeper) EndAttempt(uid int64) bool {
	_, ok := k.currentAttemps[uid]
	delete(k.currentAttemps, uid)
	return ok
}
