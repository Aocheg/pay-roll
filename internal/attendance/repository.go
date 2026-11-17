package attendance

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(a Attendance) (int64, error) {
	query := `
		INSERT INTO attendance (
			student_id,
			attendance_date,
			status,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		a.StudentID,
		a.AttendanceDate,
		a.Status,
		a.CreatedAt,
		a.UpdatedAt,
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

func (r *Repository) GetByID(id int) (Attendance, error) {
	query := `
		SELECT
			id,
			student_id,
			attendance_date,
			status,
			created_at,
			updated_at
		FROM attendance
		WHERE id = ?
	`

	var a Attendance

	err := r.db.QueryRow(query, id).Scan(
		&a.ID,
		&a.StudentID,
		&a.AttendanceDate,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return Attendance{}, err
	}

	return a, nil
}

func (r *Repository) List() ([]Attendance, error) {
	query := `
		SELECT
			id,
			student_id,
			attendance_date,
			status,
			created_at,
			updated_at
		FROM attendance
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attendanceList := make([]Attendance, 0)

	for rows.Next() {
		var a Attendance
		err := rows.Scan(
			&a.ID,
			&a.StudentID,
			&a.AttendanceDate,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		attendanceList = append(attendanceList, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendanceList, nil
}

func (r *Repository) ListByStudent(studentID int) ([]Attendance, error) {
	query := `
		SELECT
			id,
			student_id,
			attendance_date,
			status,
			created_at,
			updated_at
		FROM attendance
		WHERE student_id = ?
		ORDER BY attendance_date ASC, id ASC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attendanceList := make([]Attendance, 0)

	for rows.Next() {
		var a Attendance
		err := rows.Scan(
			&a.ID,
			&a.StudentID,
			&a.AttendanceDate,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		attendanceList = append(attendanceList, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendanceList, nil
}

func (r *Repository) GetByStudentAndDate(studentID int, attendanceDate string) (Attendance, error) {
	query := `
		SELECT
			id,
			student_id,
			attendance_date,
			status,
			created_at,
			updated_at
		FROM attendance
		WHERE student_id = ? AND attendance_date = ?
	`

	var a Attendance

	err := r.db.QueryRow(query, studentID, attendanceDate).Scan(
		&a.ID,
		&a.StudentID,
		&a.AttendanceDate,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return Attendance{}, err
	}

	return a, nil
}

func (r *Repository) Update(a Attendance) error {
	query := `
		UPDATE attendance
		SET
			student_id = ?,
			attendance_date = ?,
			status = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		a.StudentID,
		a.AttendanceDate,
		a.Status,
		a.UpdatedAt,
		a.ID,
	)
	return err
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM attendance WHERE id = ?", id)
	return err
}
