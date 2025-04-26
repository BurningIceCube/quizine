package quiz

import "time"

type QuizStatus string

const (
	STARTED         QuizStatus = "STARTED"
	AWAITING_ANSWER QuizStatus = "AWAITING_ANSWER"
	ANSWERED        QuizStatus = "ANSWERED"
	FINISHED        QuizStatus = "FINISHED"
)

type QuestionResult struct {
	QuestionID string
	Correct    bool
	TimeTaken  time.Duration
}

type Quiz struct {
	questions       []Questioner
	currentIndex    int
	score           int
	completed       bool
	status          QuizStatus
	startTime       time.Time
	creationDate    time.Time
	timeTaken       time.Duration
	correctCount    int
	questionHistory []QuestionResult
}

func NewQuiz(questions []Questioner) *Quiz {
	return &Quiz{
		questions:       questions,
		currentIndex:    0,
		score:           0,
		completed:       false,
		status:          STARTED,
		startTime:       time.Now(),
		creationDate:    time.Now(),
		timeTaken:       0,
		correctCount:    0,
		questionHistory: make([]QuestionResult, 0),
	}
}

func (q *Quiz) AmountOfQuestions() int {
	return len(q.questions)
}

func (q *Quiz) CurrentQuestion() Questioner {
	if q.currentIndex >= len(q.questions) {
		return nil
	}
	return q.questions[q.currentIndex]
}

func (q *Quiz) NextQuestion() bool {
	if q.currentIndex >= len(q.questions)-1 {
		q.completed = true
		q.status = FINISHED
		q.timeTaken = time.Since(q.startTime)
		return false
	}
	q.currentIndex++
	q.status = AWAITING_ANSWER
	return true
}

func (q *Quiz) SubmitAnswer(answer string) bool {
	if q.completed || q.currentIndex >= len(q.questions) {
		return false
	}

	current := q.CurrentQuestion()
	if current == nil {
		return false
	}

	startTime := time.Now()
	isCorrect := current.CheckAnswer(answer)
	timeTaken := time.Since(startTime)

	q.questionHistory = append(q.questionHistory, QuestionResult{
		QuestionID: current.GetID(),
		Correct:    isCorrect,
		TimeTaken:  timeTaken,
	})

	if isCorrect {
		q.score += current.GetDifficulty()
		q.correctCount++
	}

	q.status = ANSWERED
	return isCorrect
}

func (q *Quiz) TotalPoints() int {
	return q.score
}

func (q *Quiz) IsCompleted() bool {
	return q.completed
}

func (q *Quiz) Progress() (int, int) {
	return q.currentIndex + 1, len(q.questions)
}

func (q *Quiz) GetStatus() QuizStatus {
	return q.status
}

func (q *Quiz) GetTimeTaken() time.Duration {
	if q.completed {
		return q.timeTaken
	}
	return time.Since(q.startTime)
}

func (q *Quiz) GetCreationDate() time.Time {
	return q.creationDate
}

func (q *Quiz) GetCorrectCount() int {
	return q.correctCount
}

func (q *Quiz) GetQuestionHistory() []QuestionResult {
	return q.questionHistory
}
