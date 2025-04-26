package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/BurningIceCube/quizine/internal/quiz"

	_ "github.com/mattn/go-sqlite3"
)

type QuizStore struct {
	db            *sql.DB
	questionStore *QuestionStore
}

func NewQuizStore(dbPath string) (*QuizStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create tables if they don't exist
	if err := createQuizTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	questionStore, err := NewQuestionStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create question store: %v", err)
	}

	return &QuizStore{db: db, questionStore: questionStore}, nil
}

func createQuizTables(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS quizzes (
		id TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		current_index INTEGER NOT NULL,
		score INTEGER NOT NULL,
		completed BOOLEAN NOT NULL,
		start_time TIMESTAMP NOT NULL,
		creation_date TIMESTAMP NOT NULL,
		time_taken INTEGER NOT NULL,
		correct_count INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS quiz_questions (
		quiz_id TEXT NOT NULL,
		question_id TEXT NOT NULL,
		FOREIGN KEY (quiz_id) REFERENCES quizzes(id),
		FOREIGN KEY (question_id) REFERENCES questions(id),
		PRIMARY KEY (quiz_id, question_id)
	);

	CREATE TABLE IF NOT EXISTS quiz_history (
		quiz_id TEXT NOT NULL,
		question_id TEXT NOT NULL,
		correct BOOLEAN NOT NULL,
		time_taken INTEGER NOT NULL,
		FOREIGN KEY (quiz_id) REFERENCES quizzes(id),
		FOREIGN KEY (question_id) REFERENCES questions(id),
		PRIMARY KEY (quiz_id, question_id)
	);`

	_, err := db.Exec(createTableSQL)
	return err
}

func (qs *QuizStore) Close() error {
	return qs.db.Close()
}

func (qs *QuizStore) SaveQuiz(q *quiz.Quiz) error {
	tx, err := qs.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Save quiz metadata
	query := `
	INSERT INTO quizzes (id, status, current_index, score, completed, start_time, creation_date, time_taken, correct_count)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		status = excluded.status,
		current_index = excluded.current_index,
		score = excluded.score,
		completed = excluded.completed,
		time_taken = excluded.time_taken,
		correct_count = excluded.correct_count`

	_, err = tx.Exec(query,
		q.Id,
		string(q.GetStatus()),
		q.GetCurrentIndex(),
		q.GetScore(),
		q.IsCompleted(),
		q.GetStartTime(),
		q.GetCreationDate(),
		q.GetTimeTaken().Milliseconds(),
		q.GetCorrectCount(),
	)
	if err != nil {
		return fmt.Errorf("failed to save quiz: %v", err)
	}

	// Save quiz questions
	_, err = tx.Exec("DELETE FROM quiz_questions WHERE quiz_id = ?", q.Id)
	if err != nil {
		return fmt.Errorf("failed to clear quiz questions: %v", err)
	}

	for _, question := range q.GetQuestions() {
		_, err = tx.Exec("INSERT INTO quiz_questions (quiz_id, question_id) VALUES (?, ?)",
			q.Id, question.GetID())
		if err != nil {
			return fmt.Errorf("failed to save quiz question: %v", err)
		}
	}

	// Save quiz history
	_, err = tx.Exec("DELETE FROM quiz_history WHERE quiz_id = ?", q.Id)
	if err != nil {
		return fmt.Errorf("failed to clear quiz history: %v", err)
	}

	for _, result := range q.GetQuestionHistory() {
		_, err = tx.Exec("INSERT INTO quiz_history (quiz_id, question_id, correct, time_taken) VALUES (?, ?, ?, ?)",
			q.Id, result.QuestionID, result.Correct, result.TimeTaken.Milliseconds())
		if err != nil {
			return fmt.Errorf("failed to save quiz history: %v", err)
		}
	}

	return tx.Commit()
}

func (qs *QuizStore) GetQuiz(id string) (*quiz.Quiz, error) {
	// Get quiz metadata
	query := `
	SELECT status, current_index, score, completed, start_time, creation_date, time_taken, correct_count
	FROM quizzes
	WHERE id = ?`

	var (
		status       string
		currentIndex int
		score        int
		completed    bool
		startTime    time.Time
		creationDate time.Time
		timeTaken    int64
		correctCount int
	)

	err := qs.db.QueryRow(query, id).Scan(
		&status,
		&currentIndex,
		&score,
		&completed,
		&startTime,
		&creationDate,
		&timeTaken,
		&correctCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz: %v", err)
	}

	// Get quiz questions
	rows, err := qs.db.Query("SELECT question_id FROM quiz_questions WHERE quiz_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz questions: %v", err)
	}
	defer rows.Close()

	var questions []quiz.Questioner
	for rows.Next() {
		var questionID string
		if err := rows.Scan(&questionID); err != nil {
			return nil, fmt.Errorf("failed to scan question ID: %v", err)
		}
		// You'll need to implement GetQuestion in QuestionStore
		question, err := qs.questionStore.GetQuestion(questionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get question %s: %v", questionID, err)
		}
		questions = append(questions, question)
	}

	// Get quiz history
	historyRows, err := qs.db.Query("SELECT question_id, correct, time_taken FROM quiz_history WHERE quiz_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz history: %v", err)
	}
	defer historyRows.Close()

	var history []quiz.QuestionResult
	for historyRows.Next() {
		var (
			questionID string
			correct    bool
			taken      int64
		)
		if err := historyRows.Scan(&questionID, &correct, &taken); err != nil {
			return nil, fmt.Errorf("failed to scan history: %v", err)
		}
		history = append(history, quiz.QuestionResult{
			QuestionID: questionID,
			Correct:    correct,
			TimeTaken:  time.Duration(taken) * time.Millisecond,
		})
	}

	return quiz.NewQuiz(id, questions), nil
}

func (qs *QuizStore) DeleteQuiz(id string) error {
	tx, err := qs.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM quiz_history WHERE quiz_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete quiz history: %v", err)
	}

	_, err = tx.Exec("DELETE FROM quiz_questions WHERE quiz_id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete quiz questions: %v", err)
	}

	_, err = tx.Exec("DELETE FROM quizzes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete quiz: %v", err)
	}

	return tx.Commit()
}

func (qs *QuizStore) ListQuizzes() ([]*quiz.Quiz, error) {
	query := `SELECT id FROM quizzes ORDER BY created_at DESC`
	rows, err := qs.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list quizzes: %v", err)
	}
	defer rows.Close()

	var quizzes []*quiz.Quiz
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan quiz ID: %v", err)
		}
		quiz, err := qs.GetQuiz(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get quiz %s: %v", id, err)
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}
