package payroll

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(p Payroll) (int64, error) {
	query := `
		INSERT INTO payroll (
			student_id,
			period_start,
			period_end,
			expected_days,
			present_days,
			attendance_percentage,
			eligible,
			daily_rate,
			payable_days,
			gross_amount,
			payment_status,
			processed_at,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		p.StudentID,
		p.PeriodStart,
		p.PeriodEnd,
		p.ExpectedDays,
		p.PresentDays,
		p.AttendancePercentage,
		p.Eligible,
		p.DailyRate,
		p.PayableDays,
		p.GrossAmount,
		p.PaymentStatus,
		p.ProcessedAt,
		p.CreatedAt,
		p.UpdatedAt,
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

func (r *Repository) GetByID(id int) (Payroll, error) {
	query := `
		SELECT
			id,
			student_id,
			period_start,
			period_end,
			expected_days,
			present_days,
			attendance_percentage,
			eligible,
			daily_rate,
			payable_days,
			gross_amount,
			payment_status,
			processed_at,
			created_at,
			updated_at
		FROM payroll
		WHERE id = ?
	`

	var p Payroll
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.StudentID,
		&p.PeriodStart,
		&p.PeriodEnd,
		&p.ExpectedDays,
		&p.PresentDays,
		&p.AttendancePercentage,
		&p.Eligible,
		&p.DailyRate,
		&p.PayableDays,
		&p.GrossAmount,
		&p.PaymentStatus,
		&p.ProcessedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return Payroll{}, err
	}

	return p, nil
}

func (r *Repository) List() ([]Payroll, error) {
	query := `
		SELECT
			id,
			student_id,
			period_start,
			period_end,
			expected_days,
			present_days,
			attendance_percentage,
			eligible,
			daily_rate,
			payable_days,
			gross_amount,
			payment_status,
			processed_at,
			created_at,
			updated_at
		FROM payroll
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payrolls := make([]Payroll, 0)
	for rows.Next() {
		var p Payroll
		err := rows.Scan(
			&p.ID,
			&p.StudentID,
			&p.PeriodStart,
			&p.PeriodEnd,
			&p.ExpectedDays,
			&p.PresentDays,
			&p.AttendancePercentage,
			&p.Eligible,
			&p.DailyRate,
			&p.PayableDays,
			&p.GrossAmount,
			&p.PaymentStatus,
			&p.ProcessedAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payrolls = append(payrolls, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payrolls, nil
}

func (r *Repository) ListByStudent(studentID int) ([]Payroll, error) {
	query := `
		SELECT
			id,
			student_id,
			period_start,
			period_end,
			expected_days,
			present_days,
			attendance_percentage,
			eligible,
			daily_rate,
			payable_days,
			gross_amount,
			payment_status,
			processed_at,
			created_at,
			updated_at
		FROM payroll
		WHERE student_id = ?
		ORDER BY period_start DESC, id DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payrolls := make([]Payroll, 0)
	for rows.Next() {
		var p Payroll
		err := rows.Scan(
			&p.ID,
			&p.StudentID,
			&p.PeriodStart,
			&p.PeriodEnd,
			&p.ExpectedDays,
			&p.PresentDays,
			&p.AttendancePercentage,
			&p.Eligible,
			&p.DailyRate,
			&p.PayableDays,
			&p.GrossAmount,
			&p.PaymentStatus,
			&p.ProcessedAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payrolls = append(payrolls, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payrolls, nil
}

func (r *Repository) GetByStudentAndPeriod(studentID int, periodStart, periodEnd string) (Payroll, error) {
	query := `
		SELECT
			id,
			student_id,
			period_start,
			period_end,
			expected_days,
			present_days,
			attendance_percentage,
			eligible,
			daily_rate,
			payable_days,
			gross_amount,
			payment_status,
			processed_at,
			created_at,
			updated_at
		FROM payroll
		WHERE student_id = ? AND period_start = ? AND period_end = ?
	`

	var p Payroll
	err := r.db.QueryRow(query, studentID, periodStart, periodEnd).Scan(
		&p.ID,
		&p.StudentID,
		&p.PeriodStart,
		&p.PeriodEnd,
		&p.ExpectedDays,
		&p.PresentDays,
		&p.AttendancePercentage,
		&p.Eligible,
		&p.DailyRate,
		&p.PayableDays,
		&p.GrossAmount,
		&p.PaymentStatus,
		&p.ProcessedAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return Payroll{}, err
	}

	return p, nil
}

func (r *Repository) UpdateStatus(id int, status string) error {
	query := `
		UPDATE payroll
		SET payment_status = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, status, time.Now().UTC().Format(time.RFC3339), id)
	return err
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM payroll WHERE id = ?", id)
	return err
}
