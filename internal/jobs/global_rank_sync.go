package jobs

import (
	"context"

	"gorm.io/gorm"
)

// GlobalRankSync persists global_rank on student_stats from the same ordering as the leaderboard
// (active students only, by total_points DESC). Clears global_rank for inactive students.
type GlobalRankSync struct {
	DB *gorm.DB
}

// Run applies two UPDATE statements in order.
func (j *GlobalRankSync) Run(ctx context.Context) (int64, error) {
	if j.DB == nil {
		return 0, nil
	}
	db := j.DB.WithContext(ctx)
	res1 := db.Exec(`
UPDATE student_stats ss
SET global_rank = NULL, updated_at = NOW()
FROM students s
WHERE s.id = ss.student_id AND s.is_active = FALSE AND ss.global_rank IS NOT NULL`)
	if res1.Error != nil {
		return 0, res1.Error
	}
	res2 := db.Exec(`
WITH ranked AS (
	SELECT ss.student_id,
	       (RANK() OVER (ORDER BY ss.total_points DESC))::int AS rnk
	FROM student_stats ss
	INNER JOIN students s ON s.id = ss.student_id AND s.is_active = TRUE
)
UPDATE student_stats ss
SET global_rank = ranked.rnk, updated_at = NOW()
FROM ranked
WHERE ranked.student_id = ss.student_id`)
	if res2.Error != nil {
		return 0, res2.Error
	}
	return res1.RowsAffected + res2.RowsAffected, nil
}
