package db

import (
	"os"
	"testing"
	"time"

	"github.com/BurningIceCube/quizine/pkg/quiz"
)

func TestQuizStore(t *testing.T) {
	// Skip test if CGO is disabled
	if os.Getenv("CGO_ENABLED") == "0" {
		t.Skip("Skipping test because CGO is disabled")
	}

	// Create a temporary database file
	dbPath := "test_quizzes.db"
	defer os.Remove(dbPath)

	// Create a new store
	store, err := NewQuizStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create quiz store: %v", err)
	}
	defer store.Close()

	// Create some test questions
	questions := []quiz.Questioner{
		&quiz.MultiChoice{
			Id:         "q1",
			Prompt:     "What is 2+2?",
			Options:    []string{"3", "4", "5", "6"},
			Difficulty: 1,
			Answer:     "4",
			TimeLimit:  30 * time.Second,
		},
		&quiz.TrueFalse{
			Id:         "q2",
			Prompt:     "The sky is blue",
			Difficulty: 2,
			Answer:     true,
			TimeLimit:  15 * time.Second,
		},
	}

	// Save questions to the database first
	for _, q := range questions {
		if err := store.questionStore.SaveQuestion(q); err != nil {
			t.Fatalf("Failed to save question: %v", err)
		}
	}

	// Create a new quiz
	quiz1 := quiz.NewQuiz("quiz1", questions)

	// Test SaveQuiz
	if err := store.SaveQuiz(quiz1); err != nil {
		t.Errorf("Failed to save quiz: %v", err)
	}

	// Test GetQuiz
	retrievedQuiz, err := store.GetQuiz("quiz1")
	if err != nil {
		t.Errorf("Failed to get quiz: %v", err)
	}

	// Verify quiz data
	if retrievedQuiz.Id != "quiz1" {
		t.Errorf("Expected quiz ID 'quiz1', got '%s'", retrievedQuiz.Id)
	}
	if retrievedQuiz.GetStatus() != "STARTED" {
		t.Errorf("Expected status STARTED, got '%s'", retrievedQuiz.GetStatus())
	}
	if retrievedQuiz.GetCurrentIndex() != 0 {
		t.Errorf("Expected current index 0, got %d", retrievedQuiz.GetCurrentIndex())
	}
	if retrievedQuiz.GetScore() != 0 {
		t.Errorf("Expected score 0, got %d", retrievedQuiz.GetScore())
	}
	if retrievedQuiz.IsCompleted() {
		t.Error("Expected quiz not completed")
	}
	if len(retrievedQuiz.GetQuestions()) != 2 {
		t.Errorf("Expected 2 questions, got %d", len(retrievedQuiz.GetQuestions()))
	}

	// Test ListQuizzes
	quizzes, err := store.ListQuizzes()
	if err != nil {
		t.Errorf("Failed to list quizzes: %v", err)
	}
	if len(quizzes) != 1 {
		t.Errorf("Expected 1 quiz, got %d", len(quizzes))
	}

	// Test DeleteQuiz
	if err := store.DeleteQuiz("quiz1"); err != nil {
		t.Errorf("Failed to delete quiz: %v", err)
	}

	// Verify quiz was deleted
	_, err = store.GetQuiz("quiz1")
	if err == nil {
		t.Error("Expected error when getting deleted quiz")
	}

	// Test quiz with history
	quiz2 := quiz.NewQuiz("quiz2", questions)

	// Answer first question
	if quiz2.CurrentQuestion() == nil {
		t.Fatal("Expected current question not to be nil")
	}
	quiz2.SubmitAnswer("4") // Correct answer
	quiz2.NextQuestion()

	// Answer second question
	if quiz2.CurrentQuestion() == nil {
		t.Fatal("Expected current question not to be nil")
	}
	quiz2.SubmitAnswer("false") // Incorrect answer
	quiz2.NextQuestion()

	if err := store.SaveQuiz(quiz2); err != nil {
		t.Errorf("Failed to save quiz with history: %v", err)
	}

	retrievedQuiz, err = store.GetQuiz("quiz2")
	if err != nil {
		t.Errorf("Failed to get quiz with history: %v", err)
	}

	history := retrievedQuiz.GetQuestionHistory()
	if len(history) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(history))
	}
	if !history[0].Correct {
		t.Error("Expected first answer to be correct")
	}
	if history[1].Correct {
		t.Error("Expected second answer to be incorrect")
	}
}
