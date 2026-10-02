package student

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(s Student) (int64, error) {
	query := `
		INSERT INTO students (
			student_no,
			full_name,
			email,
			phone,
			program,
			daily_rate,
			active
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		s.StudentNo,
		s.FullName,
		s.Email,
		s.Phone,
		s.Program,
		s.DailyRate,
		s.Active,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetByID(id int) (Student, error) {
	query := `
		SELECT
			id,
			student_no,
			full_name,
			email,
			phone,
			program,
			daily_rate,
			active
		FROM students
		WHERE id = ?
	`

	var student Student

	err := r.db.QueryRow(query, id).Scan(
		&student.ID,
		&student.StudentNo,
		&student.FullName,
		&student.Email,
		&student.Phone,
		&student.Program,
		&student.DailyRate,
		&student.Active,
	)

	if err != nil {
		return Student{}, err
	}

	return student, nil
}

func (r *Repository) List() ([]Student, error) {
	query := `
		SELECT
			id,
			student_no,
			full_name,
			email,
			phone,
			program,
			daily_rate,
			active
		FROM students
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		var student Student

		err := rows.Scan(
			&student.ID,
			&student.StudentNo,
			&student.FullName,
			&student.Email,
			&student.Phone,
			&student.Program,
			&student.DailyRate,
			&student.Active,
		)

		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *Repository) Update(s Student) error {
	query := `
		UPDATE students
		SET
			student_no = ?,
			full_name = ?,
			email = ?,
			phone = ?,
			program = ?,
			daily_rate = ?,
			active = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		s.StudentNo,
		s.FullName,
		s.Email,
		s.Phone,
		s.Program,
		s.DailyRate,
		s.Active,
		s.ID,
	)

	return err
}

func (r *Repository) Deactivate(id int) error {
	query := `
		UPDATE students
		SET active = 0
		WHERE id = ?
	`

	_, err := r.db.Exec(query, id)

	return err
}
