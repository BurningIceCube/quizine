package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BurningIceCube/quizine/internal/quiz"

	_ "github.com/mattn/go-sqlite3"
)

type QuestionStore struct {
	db *sql.DB
}

func NewQuestionStore(dbPath string) (*QuestionStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return &QuestionStore{db: db}, nil
}

func createTables(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS questions (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		prompt TEXT NOT NULL,
		difficulty INTEGER NOT NULL,
		answer TEXT NOT NULL,
		hint TEXT,
		time_limit INTEGER NOT NULL,
		options TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSQL)
	return err
}

func (qs *QuestionStore) Close() error {
	return qs.db.Close()
}

func (qs *QuestionStore) SaveQuestion(q quiz.Questioner) error {
	var optionsJSON string
	var questionType string

	switch q := q.(type) {
	case *quiz.MultiChoice:
		optionsJSONBytes, err := json.Marshal(q.Options)
		if err != nil {
			return fmt.Errorf("failed to marshal options: %v", err)
		}
		optionsJSON = string(optionsJSONBytes)
		questionType = "MULTI_CHOICE"
	case *quiz.TrueFalse:
		questionType = "TRUE_FALSE"
	case *quiz.FillIn:
		questionType = "FILL_IN"
	default:
		return fmt.Errorf("unknown question type")
	}

	query := `
	INSERT INTO questions (id, type, prompt, difficulty, answer, hint, time_limit, options)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		type = excluded.type,
		prompt = excluded.prompt,
		difficulty = excluded.difficulty,
		answer = excluded.answer,
		hint = excluded.hint,
		time_limit = excluded.time_limit,
		options = excluded.options`

	_, err := qs.db.Exec(query,
		q.GetID(),
		questionType,
		q.GetPrompt(),
		q.GetDifficulty(),
		getAnswerString(q),
		getHint(q),
		q.GetTimeLimit().Milliseconds(),
		optionsJSON,
	)

	return err
}

func (qs *QuestionStore) GetQuestion(id string) (quiz.Questioner, error) {
	query := `
	SELECT type, prompt, difficulty, answer, hint, time_limit, options
	FROM questions
	WHERE id = ?`

	var (
		questionType string
		prompt       string
		difficulty   int
		answer       string
		hint         string
		timeLimit    int64
		optionsJSON  string
	)

	err := qs.db.QueryRow(query, id).Scan(
		&questionType,
		&prompt,
		&difficulty,
		&answer,
		&hint,
		&timeLimit,
		&optionsJSON,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get question: %v", err)
	}

	switch questionType {
	case "MULTI_CHOICE":
		var options []string
		if err := json.Unmarshal([]byte(optionsJSON), &options); err != nil {
			return nil, fmt.Errorf("failed to unmarshal options: %v", err)
		}
		return &quiz.MultiChoice{
			Id:         id,
			Prompt:     prompt,
			Options:    options,
			Difficulty: difficulty,
			Answer:     answer,
			Hint:       hint,
			TimeLimit:  time.Duration(timeLimit) * time.Millisecond,
		}, nil
	case "TRUE_FALSE":
		return &quiz.TrueFalse{
			Id:         id,
			Prompt:     prompt,
			Difficulty: difficulty,
			Answer:     answer == "true",
			Hint:       hint,
			TimeLimit:  time.Duration(timeLimit) * time.Millisecond,
		}, nil
	case "FILL_IN":
		return &quiz.FillIn{
			Id:         id,
			Prompt:     prompt,
			Difficulty: difficulty,
			Answer:     answer,
			Hint:       hint,
			TimeLimit:  time.Duration(timeLimit) * time.Millisecond,
		}, nil
	default:
		return nil, fmt.Errorf("unknown question type: %s", questionType)
	}
}

func (qs *QuestionStore) DeleteQuestion(id string) error {
	query := `DELETE FROM questions WHERE id = ?`
	_, err := qs.db.Exec(query, id)
	return err
}

func (qs *QuestionStore) ListQuestions() ([]quiz.Questioner, error) {
	query := `SELECT id FROM questions ORDER BY created_at DESC`
	rows, err := qs.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list questions: %v", err)
	}
	defer rows.Close()

	var questions []quiz.Questioner
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan question ID: %v", err)
		}
		question, err := qs.GetQuestion(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get question %s: %v", id, err)
		}
		questions = append(questions, question)
	}

	return questions, nil
}

// Helper functions
func getAnswerString(q quiz.Questioner) string {
	switch q := q.(type) {
	case *quiz.MultiChoice:
		return q.Answer
	case *quiz.TrueFalse:
		if q.Answer {
			return "true"
		}
		return "false"
	case *quiz.FillIn:
		return q.Answer
	default:
		return ""
	}
}

func getHint(q quiz.Questioner) string {
	switch q := q.(type) {
	case *quiz.MultiChoice:
		return q.Hint
	case *quiz.TrueFalse:
		return q.Hint
	case *quiz.FillIn:
		return q.Hint
	default:
		return ""
	}
}
