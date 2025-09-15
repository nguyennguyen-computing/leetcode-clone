package repository

import (
	"database/sql"

	"leetcode-clone-backend/pkg/models"
)

// userProgressRepository implements UserProgressRepository interface
type userProgressRepository struct {
	db *sql.DB
}

// NewUserProgressRepository creates a new user progress repository
func NewUserProgressRepository(db *sql.DB) UserProgressRepository {
	return &userProgressRepository{db: db}
}

// Create creates a new user progress record
func (r *userProgressRepository) Create(progress *models.UserProgress) (*models.UserProgress, error) {
	query := `
		INSERT INTO user_progress (user_id, problem_id, is_solved, best_submission_id, attempts, first_solved_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, problem_id, is_solved, best_submission_id, attempts, first_solved_at`

	var created models.UserProgress
	err := r.db.QueryRow(
		query,
		progress.UserID,
		progress.ProblemID,
		progress.IsSolved,
		progress.BestSubmissionID,
		progress.Attempts,
		progress.FirstSolvedAt,
	).Scan(
		&created.UserID,
		&created.ProblemID,
		&created.IsSolved,
		&created.BestSubmissionID,
		&created.Attempts,
		&created.FirstSolvedAt,
	)

	if err != nil {
		return nil, NewRepositoryError("Create", err, "database_error")
	}

	return &created, nil
}

// GetByUserAndProblem retrieves user progress for a specific user and problem
func (r *userProgressRepository) GetByUserAndProblem(userID, problemID int) (*models.UserProgress, error) {
	query := `
		SELECT user_id, problem_id, is_solved, best_submission_id, attempts, first_solved_at
		FROM user_progress
		WHERE user_id = $1 AND problem_id = $2`

	var progress models.UserProgress
	err := r.db.QueryRow(query, userID, problemID).Scan(
		&progress.UserID,
		&progress.ProblemID,
		&progress.IsSolved,
		&progress.BestSubmissionID,
		&progress.Attempts,
		&progress.FirstSolvedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetByUserAndProblem", ErrNotFound, "progress_not_found")
		}
		return nil, NewRepositoryError("GetByUserAndProblem", err, "database_error")
	}

	return &progress, nil
}

// GetByUserID retrieves all progress records for a user
func (r *userProgressRepository) GetByUserID(userID int) ([]*models.UserProgress, error) {
	query := `
		SELECT user_id, problem_id, is_solved, best_submission_id, attempts, first_solved_at
		FROM user_progress
		WHERE user_id = $1
		ORDER BY first_solved_at DESC NULLS LAST, problem_id ASC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, NewRepositoryError("GetByUserID", err, "database_error")
	}
	defer rows.Close()

	var progressList []*models.UserProgress
	for rows.Next() {
		var progress models.UserProgress
		err := rows.Scan(
			&progress.UserID,
			&progress.ProblemID,
			&progress.IsSolved,
			&progress.BestSubmissionID,
			&progress.Attempts,
			&progress.FirstSolvedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("GetByUserID", err, "scan_error")
		}
		progressList = append(progressList, &progress)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("GetByUserID", err, "rows_error")
	}

	return progressList, nil
}

// Update updates an existing user progress record
func (r *userProgressRepository) Update(progress *models.UserProgress) (*models.UserProgress, error) {
	query := `
		UPDATE user_progress
		SET is_solved = $3, best_submission_id = $4, attempts = $5, first_solved_at = $6
		WHERE user_id = $1 AND problem_id = $2
		RETURNING user_id, problem_id, is_solved, best_submission_id, attempts, first_solved_at`

	var updated models.UserProgress
	err := r.db.QueryRow(
		query,
		progress.UserID,
		progress.ProblemID,
		progress.IsSolved,
		progress.BestSubmissionID,
		progress.Attempts,
		progress.FirstSolvedAt,
	).Scan(
		&updated.UserID,
		&updated.ProblemID,
		&updated.IsSolved,
		&updated.BestSubmissionID,
		&updated.Attempts,
		&updated.FirstSolvedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("Update", ErrNotFound, "progress_not_found")
		}
		return nil, NewRepositoryError("Update", err, "database_error")
	}

	return &updated, nil
}

// Delete deletes a user progress record
func (r *userProgressRepository) Delete(userID, problemID int) error {
	query := `DELETE FROM user_progress WHERE user_id = $1 AND problem_id = $2`

	result, err := r.db.Exec(query, userID, problemID)
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	if rowsAffected == 0 {
		return NewRepositoryError("Delete", ErrNotFound, "progress_not_found")
	}

	return nil
}

// GetSolvedCount returns the total number of problems solved by a user
func (r *userProgressRepository) GetSolvedCount(userID int) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM user_progress
		WHERE user_id = $1 AND is_solved = true`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, NewRepositoryError("GetSolvedCount", err, "database_error")
	}

	return count, nil
}

// GetSolvedCountByDifficulty returns the number of problems solved by difficulty
func (r *userProgressRepository) GetSolvedCountByDifficulty(userID int) (map[string]int, error) {
	query := `
		SELECT p.difficulty, COUNT(*)
		FROM user_progress up
		JOIN problems p ON up.problem_id = p.id
		WHERE up.user_id = $1 AND up.is_solved = true
		GROUP BY p.difficulty`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, NewRepositoryError("GetSolvedCountByDifficulty", err, "database_error")
	}
	defer rows.Close()

	counts := make(map[string]int)
	// Initialize with zero counts for all difficulties
	counts[models.DifficultyEasy] = 0
	counts[models.DifficultyMedium] = 0
	counts[models.DifficultyHard] = 0

	for rows.Next() {
		var difficulty string
		var count int
		err := rows.Scan(&difficulty, &count)
		if err != nil {
			return nil, NewRepositoryError("GetSolvedCountByDifficulty", err, "scan_error")
		}
		counts[difficulty] = count
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("GetSolvedCountByDifficulty", err, "rows_error")
	}

	return counts, nil
}