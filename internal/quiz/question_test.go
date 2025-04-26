package quiz

import (
	"testing"
	"time"
)

func TestMultiChoice(t *testing.T) {
	mc := &MultiChoice{
		Id:         "mc1",
		Prompt:     "What is 2+2?",
		Options:    []string{"3", "4", "5", "6"},
		Difficulty: 1,
		Answer:     "4",
		Hint:       "It's an even number",
		TimeLimit:  30 * time.Second,
	}

	// Test GetID
	if mc.GetID() != "mc1" {
		t.Errorf("Expected ID 'mc1', got '%s'", mc.GetID())
	}

	// Test GetPrompt
	if mc.GetPrompt() != "What is 2+2?" {
		t.Errorf("Expected prompt 'What is 2+2?', got '%s'", mc.GetPrompt())
	}

	// Test GetDifficulty
	if mc.GetDifficulty() != 1 {
		t.Errorf("Expected difficulty 1, got %d", mc.GetDifficulty())
	}

	// Test GetTimeLimit
	if mc.GetTimeLimit() != 30*time.Second {
		t.Errorf("Expected time limit 30s, got %v", mc.GetTimeLimit())
	}

	// Test CheckAnswer
	if !mc.CheckAnswer("4") {
		t.Error("Expected correct answer '4' to return true")
	}
	if mc.CheckAnswer("3") {
		t.Error("Expected incorrect answer '3' to return false")
	}
}

func TestTrueFalse(t *testing.T) {
	tf := &TrueFalse{
		Id:         "tf1",
		Prompt:     "The sky is blue",
		Difficulty: 2,
		Answer:     true,
		Hint:       "Think about daytime",
		TimeLimit:  15 * time.Second,
	}

	// Test GetID
	if tf.GetID() != "tf1" {
		t.Errorf("Expected ID 'tf1', got '%s'", tf.GetID())
	}

	// Test GetPrompt
	if tf.GetPrompt() != "The sky is blue" {
		t.Errorf("Expected prompt 'The sky is blue', got '%s'", tf.GetPrompt())
	}

	// Test GetDifficulty
	if tf.GetDifficulty() != 2 {
		t.Errorf("Expected difficulty 2, got %d", tf.GetDifficulty())
	}

	// Test GetTimeLimit
	if tf.GetTimeLimit() != 15*time.Second {
		t.Errorf("Expected time limit 15s, got %v", tf.GetTimeLimit())
	}

	// Test CheckAnswer
	if !tf.CheckAnswer("true") {
		t.Error("Expected correct answer 'true' to return true")
	}
	if tf.CheckAnswer("false") {
		t.Error("Expected incorrect answer 'false' to return false")
	}
}

func TestFillIn(t *testing.T) {
	fi := &FillIn{
		Id:         "fi1",
		Prompt:     "The capital of France is ___",
		Difficulty: 3,
		Answer:     "Paris",
		Hint:       "It starts with P",
		TimeLimit:  45 * time.Second,
	}

	// Test GetID
	if fi.GetID() != "fi1" {
		t.Errorf("Expected ID 'fi1', got '%s'", fi.GetID())
	}

	// Test GetPrompt
	if fi.GetPrompt() != "The capital of France is ___" {
		t.Errorf("Expected prompt 'The capital of France is ___', got '%s'", fi.GetPrompt())
	}

	// Test GetDifficulty
	if fi.GetDifficulty() != 3 {
		t.Errorf("Expected difficulty 3, got %d", fi.GetDifficulty())
	}

	// Test GetTimeLimit
	if fi.GetTimeLimit() != 45*time.Second {
		t.Errorf("Expected time limit 45s, got %v", fi.GetTimeLimit())
	}

	// Test CheckAnswer
	if !fi.CheckAnswer("Paris") {
		t.Error("Expected correct answer 'Paris' to return true")
	}
	if fi.CheckAnswer("London") {
		t.Error("Expected incorrect answer 'London' to return false")
	}
}

func TestQuestionerInterface(t *testing.T) {
	questions := []Questioner{
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

	for _, q := range questions {
		// Test that all methods can be called without error
		_ = q.GetID()
		_ = q.GetPrompt()
		_ = q.GetDifficulty()
		_ = q.GetTimeLimit()
		_ = q.CheckAnswer("test")
	}
}
