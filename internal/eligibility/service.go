package eligibility

func IsEligible(attendancePercentage float64) bool {
	return attendancePercentage >= 70
}
