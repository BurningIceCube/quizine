package db

import (
	"os"
	"testing"
	"time"

	"github.com/BurningIceCube/quizine/pkg/quiz"
)

func TestQuestionStore(t *testing.T) {
	// Skip test if CGO is disabled
	if os.Getenv("CGO_ENABLED") == "0" {
		t.Skip("Skipping test because CGO is disabled")
	}

	// Create a temporary database file
	dbPath := "test_questions.db"
	defer os.Remove(dbPath)

	// Create a new store
	store, err := NewQuestionStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create question store: %v", err)
	}
	defer store.Close()

	// Test MultiChoice question
	mc := &quiz.MultiChoice{
		Id:         "mc1",
		Prompt:     "What is 2+2?",
		Options:    []string{"3", "4", "5", "6"},
		Difficulty: 1,
		Answer:     "4",
		Hint:       "It's an even number",
		TimeLimit:  30 * time.Second,
	}

	// Test SaveQuestion
	if err := store.SaveQuestion(mc); err != nil {
		t.Errorf("Failed to save MultiChoice question: %v", err)
	}

	// Test GetQuestion
	retrievedMC, err := store.GetQuestion("mc1")
	if err != nil {
		t.Errorf("Failed to get MultiChoice question: %v", err)
	}

	// Verify MultiChoice data
	if mcQ, ok := retrievedMC.(*quiz.MultiChoice); ok {
		if mcQ.Id != "mc1" {
			t.Errorf("Expected ID 'mc1', got '%s'", mcQ.Id)
		}
		if mcQ.Prompt != "What is 2+2?" {
			t.Errorf("Expected prompt 'What is 2+2?', got '%s'", mcQ.Prompt)
		}
		if len(mcQ.Options) != 4 {
			t.Errorf("Expected 4 options, got %d", len(mcQ.Options))
		}
		if mcQ.Answer != "4" {
			t.Errorf("Expected answer '4', got '%s'", mcQ.Answer)
		}
		if mcQ.Hint != "It's an even number" {
			t.Errorf("Expected hint 'It's an even number', got '%s'", mcQ.Hint)
		}
		if mcQ.TimeLimit != 30*time.Second {
			t.Errorf("Expected time limit 30s, got %v", mcQ.TimeLimit)
		}
	} else {
		t.Error("Expected MultiChoice type")
	}

	// Test TrueFalse question
	tf := &quiz.TrueFalse{
		Id:         "tf1",
		Prompt:     "The sky is blue",
		Difficulty: 2,
		Answer:     true,
		Hint:       "Think about daytime",
		TimeLimit:  15 * time.Second,
	}

	if err := store.SaveQuestion(tf); err != nil {
		t.Errorf("Failed to save TrueFalse question: %v", err)
	}

	retrievedTF, err := store.GetQuestion("tf1")
	if err != nil {
		t.Errorf("Failed to get TrueFalse question: %v", err)
	}

	// Verify TrueFalse data
	if tfQ, ok := retrievedTF.(*quiz.TrueFalse); ok {
		if tfQ.Id != "tf1" {
			t.Errorf("Expected ID 'tf1', got '%s'", tfQ.Id)
		}
		if tfQ.Prompt != "The sky is blue" {
			t.Errorf("Expected prompt 'The sky is blue', got '%s'", tfQ.Prompt)
		}
		if !tfQ.Answer {
			t.Error("Expected answer true")
		}
		if tfQ.Hint != "Think about daytime" {
			t.Errorf("Expected hint 'Think about daytime', got '%s'", tfQ.Hint)
		}
		if tfQ.TimeLimit != 15*time.Second {
			t.Errorf("Expected time limit 15s, got %v", tfQ.TimeLimit)
		}
	} else {
		t.Error("Expected TrueFalse type")
	}

	// Test FillIn question
	fi := &quiz.FillIn{
		Id:         "fi1",
		Prompt:     "The capital of France is ___",
		Difficulty: 3,
		Answer:     "Paris",
		Hint:       "It starts with P",
		TimeLimit:  45 * time.Second,
	}

	if err := store.SaveQuestion(fi); err != nil {
		t.Errorf("Failed to save FillIn question: %v", err)
	}

	retrievedFI, err := store.GetQuestion("fi1")
	if err != nil {
		t.Errorf("Failed to get FillIn question: %v", err)
	}

	// Verify FillIn data
	if fiQ, ok := retrievedFI.(*quiz.FillIn); ok {
		if fiQ.Id != "fi1" {
			t.Errorf("Expected ID 'fi1', got '%s'", fiQ.Id)
		}
		if fiQ.Prompt != "The capital of France is ___" {
			t.Errorf("Expected prompt 'The capital of France is ___', got '%s'", fiQ.Prompt)
		}
		if fiQ.Answer != "Paris" {
			t.Errorf("Expected answer 'Paris', got '%s'", fiQ.Answer)
		}
		if fiQ.Hint != "It starts with P" {
			t.Errorf("Expected hint 'It starts with P', got '%s'", fiQ.Hint)
		}
		if fiQ.TimeLimit != 45*time.Second {
			t.Errorf("Expected time limit 45s, got %v", fiQ.TimeLimit)
		}
	} else {
		t.Error("Expected FillIn type")
	}

	// Test ListQuestions
	questions, err := store.ListQuestions()
	if err != nil {
		t.Errorf("Failed to list questions: %v", err)
	}
	if len(questions) != 3 {
		t.Errorf("Expected 3 questions, got %d", len(questions))
	}

	// Test DeleteQuestion
	if err := store.DeleteQuestion("mc1"); err != nil {
		t.Errorf("Failed to delete question: %v", err)
	}

	// Verify question was deleted
	_, err = store.GetQuestion("mc1")
	if err == nil {
		t.Error("Expected error when getting deleted question")
	}
}
