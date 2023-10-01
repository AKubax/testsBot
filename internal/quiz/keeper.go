package quiz

import (
	"encoding/json"
	"os"
)

func NewKeeper() *Keeper {
	return &Keeper{
		quizesByName: map[string]Quiz{},
		quizesByUser: map[int64]string{},
	}
}

type Keeper struct {
	quizesByName map[string]Quiz
	quizesByUser map[int64]string
}

func (k *Keeper) LoadQuizzesFromFiles(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var quizzes map[string]string
	err = json.Unmarshal(data, &quizzes)
	if err != nil {
		return err
	}

	for quizName, fileName := range quizzes {
		err := k.LoadQuizFromFile(quizName, fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) LoadQuizFromFile(quizName, fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var quiz Quiz
	err = json.Unmarshal(data, &quiz)
	if err != nil {
		return err
	}

	k.quizesByName[quizName] = quiz
	return nil
}

func (k *Keeper) GetAllQuizzes() map[string]Quiz {
	return k.quizesByName
}

func (k *Keeper) GetQuiz(name string) (Quiz, bool) {
	quiz, ok := k.quizesByName[name]
	return quiz, ok
}

func (k *Keeper) GetCurrentQuiz(uid int64) (string, bool) {
	quizName, chosen := k.quizesByUser[uid]
	return quizName, chosen
}

func (k *Keeper) SetCurrentQuiz(uid int64, name string) {
	k.quizesByUser[uid] = name
}

func (k *Keeper) UnsetCurrentQuiz(uid int64) {
	delete(k.quizesByUser, uid)
}
