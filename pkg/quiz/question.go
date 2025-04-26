package quiz

import (
	"time"
)

type MultiChoice struct {
	Id         string        `json:"id"`
	Prompt     string        `json:"prompt"`
	Options    []string      `json:"options"`
	Difficulty int           `json:"difficulty"`
	Answer     string        `json:"answer"`
	Hint       string        `json:"hint"`
	TimeLimit  time.Duration `json:"timeLimit"`
}

func (mc *MultiChoice) GetPrompt() string {
	return mc.Prompt
}

func (mc *MultiChoice) GetID() string {
	return mc.Id
}

func (mc *MultiChoice) CheckAnswer(answer string) bool {
	return answer == mc.Answer
}

func (mc *MultiChoice) GetDifficulty() int {
	return mc.Difficulty
}

func (mc *MultiChoice) GetTimeLimit() time.Duration {
	return mc.TimeLimit
}

type TrueFalse struct {
	Id         string        `json:"id"`
	Prompt     string        `json:"prompt"`
	Difficulty int           `json:"difficulty"`
	Answer     bool          `json:"answer"`
	Hint       string        `json:"hint"`
	TimeLimit  time.Duration `json:"timeLimit"`
}

func (tf *TrueFalse) GetPrompt() string {
	return tf.Prompt
}

func (tf *TrueFalse) GetID() string {
	return tf.Id
}

func (tf *TrueFalse) CheckAnswer(answer string) bool {
	return answer == "true"
}

func (tf *TrueFalse) GetDifficulty() int {
	return tf.Difficulty
}

func (tf *TrueFalse) GetTimeLimit() time.Duration {
	return tf.TimeLimit
}

type FillIn struct {
	Id         string        `json:"id"`
	Prompt     string        `json:"prompt"`
	Difficulty int           `json:"difficulty"`
	Answer     string        `json:"answer"`
	Hint       string        `json:"hint"`
	TimeLimit  time.Duration `json:"timeLimit"`
}

func (fi *FillIn) GetPrompt() string {
	return fi.Prompt
}

func (fi *FillIn) GetID() string {
	return fi.Id
}

func (fi *FillIn) CheckAnswer(answer string) bool {
	return answer == fi.Answer
}

func (fi *FillIn) GetDifficulty() int {
	return fi.Difficulty
}

func (fi *FillIn) GetTimeLimit() time.Duration {
	return fi.TimeLimit
}

type Questioner interface {
	GetID() string
	GetPrompt() string
	GetDifficulty() int
	GetTimeLimit() time.Duration
	CheckAnswer(answer string) bool
}
