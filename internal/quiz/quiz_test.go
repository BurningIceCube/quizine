package quiz

import (
	"testing"
	"time"
)

func createTestQuestions() []Questioner {
	return []Questioner{
		&MultiChoice{
			Id:         "mc1",
			Prompt:     "What is 2+2?",
			Options:    []string{"3", "4", "5", "6"},
			Difficulty: 1,
			Answer:     "4",
			TimeLimit:  30 * time.Second,
		},
		&TrueFalse{
			Id:         "tf1",
			Prompt:     "The sky is blue",
			Difficulty: 2,
			Answer:     true,
			TimeLimit:  15 * time.Second,
		},
		&FillIn{
			Id:         "fi1",
			Prompt:     "The capital of France is ___",
			Difficulty: 3,
			Answer:     "Paris",
			TimeLimit:  45 * time.Second,
		},
	}
}

func TestNewQuiz(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	if quiz.Id != "quiz1" {
		t.Errorf("Expected ID 'quiz1', got '%s'", quiz.Id)
	}
	if quiz.AmountOfQuestions() != 3 {
		t.Errorf("Expected 3 questions, got %d", quiz.AmountOfQuestions())
	}
	if quiz.TotalPoints() != 0 {
		t.Errorf("Expected initial score 0, got %d", quiz.TotalPoints())
	}
	if quiz.IsCompleted() {
		t.Error("Expected quiz not to be completed initially")
	}
	if quiz.GetStatus() != STARTED {
		t.Errorf("Expected status STARTED, got %v", quiz.GetStatus())
	}
	if quiz.GetCorrectCount() != 0 {
		t.Errorf("Expected 0 correct answers, got %d", quiz.GetCorrectCount())
	}
	if len(quiz.GetQuestionHistory()) != 0 {
		t.Errorf("Expected empty question history, got %d entries", len(quiz.GetQuestionHistory()))
	}
}

func TestQuizStatus(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	// Test initial status
	if quiz.GetStatus() != STARTED {
		t.Errorf("Expected status STARTED, got %v", quiz.GetStatus())
	}

	// Test status after first question
	quiz.status = AWAITING_ANSWER
	if quiz.GetStatus() != AWAITING_ANSWER {
		t.Errorf("Expected status AWAITING_ANSWER, got %v", quiz.GetStatus())
	}

	// Test status after answer
	quiz.SubmitAnswer("4")
	if quiz.GetStatus() != ANSWERED {
		t.Errorf("Expected status ANSWERED, got %v", quiz.GetStatus())
	}

	// Test status after completion
	quiz.completed = true
	quiz.status = FINISHED
	if quiz.GetStatus() != FINISHED {
		t.Errorf("Expected status FINISHED, got %v", quiz.GetStatus())
	}

	// Test SAVED status
	quiz.status = SAVED
	if quiz.GetStatus() != SAVED {
		t.Errorf("Expected status SAVED, got %v", quiz.GetStatus())
	}

	// Test QUIT status
	quiz.status = QUIT
	if quiz.GetStatus() != QUIT {
		t.Errorf("Expected status QUIT, got %v", quiz.GetStatus())
	}
}

func TestTimeTracking(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	// Test creation date is set
	if quiz.GetCreationDate().IsZero() {
		t.Error("Expected creation date to be set")
	}

	// Test initial time taken
	if quiz.GetTimeTaken() < 0 {
		t.Error("Expected positive time taken")
	}

	// Test time taken after completion
	quiz.completed = true
	quiz.timeTaken = 30 * time.Second
	if quiz.GetTimeTaken() != 30*time.Second {
		t.Errorf("Expected 30s time taken, got %v", quiz.GetTimeTaken())
	}
}

func TestQuestionHistory(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	// Test correct answer history
	quiz.SubmitAnswer("4")
	history := quiz.GetQuestionHistory()
	if len(history) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(history))
	}
	if !history[0].Correct {
		t.Error("Expected correct answer in history")
	}
	if history[0].QuestionID != "mc1" {
		t.Errorf("Expected question ID 'mc1', got '%s'", history[0].QuestionID)
	}

	// Test incorrect answer history
	quiz.NextQuestion()
	quiz.SubmitAnswer("false")
	history = quiz.GetQuestionHistory()
	if len(history) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(history))
	}
	if history[1].Correct {
		t.Error("Expected incorrect answer in history")
	}
}

func TestCorrectCount(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	// Test initial count
	if quiz.GetCorrectCount() != 0 {
		t.Errorf("Expected 0 correct answers, got %d", quiz.GetCorrectCount())
	}

	// Test after correct answer
	quiz.SubmitAnswer("4")
	if quiz.GetCorrectCount() != 1 {
		t.Errorf("Expected 1 correct answer, got %d", quiz.GetCorrectCount())
	}

	// Test after incorrect answer
	quiz.NextQuestion()
	quiz.SubmitAnswer("false")
	if quiz.GetCorrectCount() != 1 {
		t.Errorf("Expected 1 correct answer, got %d", quiz.GetCorrectCount())
	}
}

func TestCurrentQuestion(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	// Test first question
	current := quiz.CurrentQuestion()
	if current == nil {
		t.Error("Expected first question, got nil")
	}
	if current.GetID() != "mc1" {
		t.Errorf("Expected question ID 'mc1', got '%s'", current.GetID())
	}

	// Test after completion
	quiz.currentIndex = len(questions)
	current = quiz.CurrentQuestion()
	if current != nil {
		t.Error("Expected nil after completion, got question")
	}
}

func TestNextQuestion(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	// Test first next
	if !quiz.NextQuestion() {
		t.Error("Expected successful next question")
	}
	if quiz.currentIndex != 1 {
		t.Errorf("Expected index 1, got %d", quiz.currentIndex)
	}

	// Test second next
	if !quiz.NextQuestion() {
		t.Error("Expected successful next question")
	}
	if quiz.currentIndex != 2 {
		t.Errorf("Expected index 2, got %d", quiz.currentIndex)
	}

	// Test completion
	if quiz.NextQuestion() {
		t.Error("Expected no more questions")
	}
	if !quiz.IsCompleted() {
		t.Error("Expected quiz to be completed")
	}
}

func TestSubmitAnswer(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	// Test correct answer
	if !quiz.SubmitAnswer("4") {
		t.Error("Expected correct answer to be accepted")
	}
	if quiz.TotalPoints() != 1 {
		t.Errorf("Expected score 1, got %d", quiz.TotalPoints())
	}

	// Test incorrect answer
	if quiz.SubmitAnswer("3") {
		t.Error("Expected incorrect answer to be rejected")
	}
	if quiz.TotalPoints() != 1 {
		t.Errorf("Expected score to remain 1, got %d", quiz.TotalPoints())
	}

	// Test after completion
	quiz.completed = true
	if quiz.SubmitAnswer("4") {
		t.Error("Expected no answer submission after completion")
	}
}

func TestProgress(t *testing.T) {
	questions := createTestQuestions()
	quiz := NewQuiz("quiz1", questions)

	// Test initial progress
	current, total := quiz.Progress()
	if current != 1 || total != 3 {
		t.Errorf("Expected progress 1/3, got %d/%d", current, total)
	}

	// Test after one question
	quiz.NextQuestion()
	current, total = quiz.Progress()
	if current != 2 || total != 3 {
		t.Errorf("Expected progress 2/3, got %d/%d", current, total)
	}

	// Test after completion
	quiz.NextQuestion()
	quiz.NextQuestion()
	current, total = quiz.Progress()
	if current != 3 || total != 3 {
		t.Errorf("Expected progress 3/3, got %d/%d", current, total)
	}
}
