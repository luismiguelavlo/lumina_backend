package jobs

import (
	"context"

	"library_back/internal/repositories"
)

// LoanOverdueSync marks loans as overdue when due_date is before today (DB date)
// and status is still active. Runs periodically from the application (not pg_cron).
type LoanOverdueSync struct {
	Loans repositories.LoanRepository
}

// Run applies UPDATE loans SET status = overdue WHERE status = active AND due_date < CURRENT_DATE.
func (j *LoanOverdueSync) Run(ctx context.Context) (int64, error) {
	return j.Loans.MarkActiveLoansOverdue(ctx)
}
