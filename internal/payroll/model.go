package payroll

type Payroll struct {
	ID                   int     `json:"id"`
	StudentID            int     `json:"student_id"`
	PeriodStart          string  `json:"period_start"`
	PeriodEnd            string  `json:"period_end"`
	ExpectedDays         int     `json:"expected_days"`
	PresentDays          int     `json:"present_days"`
	AttendancePercentage float64 `json:"attendance_percentage"`
	Eligible             bool    `json:"eligible"`
	DailyRate            float64 `json:"daily_rate"`
	PayableDays          int     `json:"payable_days"`
	GrossAmount          float64 `json:"gross_amount"`
	PaymentStatus        string  `json:"payment_status"`
	ProcessedAt          string  `json:"processed_at"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}
