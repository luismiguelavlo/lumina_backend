package main

import (
	"context"
	"library_back/internal/services/analytics"
	"library_back/internal/services/author"
	"library_back/internal/services/auth"
	"library_back/internal/services/badgedef"
	"library_back/internal/services/book"
	"library_back/internal/services/fine"
	"library_back/internal/services/loan"
	"library_back/internal/services/ranking"
	"library_back/internal/services/reputation"
	"library_back/internal/services/sanction"
	"library_back/internal/services/student"
	"log"
	"os"
	"strings"
	"time"

	"library_back/internal/handlers"
	"library_back/internal/jobs"
	"library_back/internal/pkg/database"
	"library_back/internal/pkg/ratelimit"
	"library_back/internal/repositories"
	"library_back/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET is required (min 32 characters)")
	}

	orm, err := database.Open(databaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	sqlDB, err := orm.DB()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()
	log.Println("database: connected (PostgreSQL via GORM)")
	if err := database.Migrate(orm); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("database: schema synchronized (GORM AutoMigrate)")

	userRepo := repositories.NewUserRepositoryGorm(orm)
	revokedRepo := repositories.NewRevokedTokenRepositoryGorm(orm)
	tokenSvc, err := services.NewTokenService(jwtSecret, revokedRepo)
	if err != nil {
		log.Fatalf("token service: %v", err)
	}
	authSvc := auth.NewAuthService(userRepo, tokenSvc)

	catalogBooks := repositories.NewCatalogBookRepositoryGorm(orm)
	catalogAuthors := repositories.NewAuthorCatalogRepositoryGorm(orm)
	catalogGenres := repositories.NewGenreCatalogRepositoryGorm(orm)
	authorSvc := author.NewService(catalogAuthors)
	bookSvc := book.NewBookService(catalogBooks, catalogAuthors, catalogGenres)

	deptRepo := repositories.NewStudentDepartmentRepositoryGorm(orm)
	studentRepo := repositories.NewStudentRepositoryGorm(orm)
	studentProfileRepo := repositories.NewStudentProfileRepositoryGorm(orm)
	badgeRepo := repositories.NewBadgeRepositoryGorm(orm)
	studentBadgeRepo := repositories.NewStudentBadgeRepositoryGorm(orm)

	loanRepo := repositories.NewLoanRepositoryGorm(orm)
	sanctionRepo := repositories.NewSanctionRepositoryGorm(orm)
	repStore := repositories.NewReputationStoreGorm(orm)
	repSvc := reputation.NewService(repStore)
	studentSvc := student.NewStudentService(studentRepo, deptRepo, studentProfileRepo, badgeRepo, studentBadgeRepo, repSvc)

	activityLogRepo := repositories.NewActivityLogRepositoryGorm(orm)
	fineRepo := repositories.NewFineRepositoryGorm(orm)
	fineSvc := fine.NewService(fineRepo, loanRepo, repSvc)
	loanSvc := loan.NewService(loanRepo, catalogBooks, studentRepo, repSvc, sanctionRepo, activityLogRepo, fineSvc)
	sanctionSvc := sanction.NewService(sanctionRepo, studentRepo, repSvc)
	studentSvc.SetActivityLog(activityLogRepo)
	fineSvc.SetActivityLog(activityLogRepo)
	sanctionSvc.SetActivityLog(activityLogRepo)

	analyticsRepo := repositories.NewAnalyticsRepositoryGorm(orm)
	analyticsSvc := analytics.NewService(analyticsRepo)

	rankingRepo := repositories.NewRankingRepositoryGorm(orm)
	rankingSvc := ranking.NewService(rankingRepo, studentRepo)

	healthHandler := handlers.NewHealthHandler(orm)
	authHandler := handlers.NewAuthHandler(authSvc)
	bookHandler := handlers.NewBookHandler(bookSvc)
	catalogListHandler := handlers.NewCatalogListHandler(catalogAuthors, catalogGenres)
	authorHandler := handlers.NewAuthorHandler(authorSvc)
	departmentHandler := handlers.NewDepartmentHandler(deptRepo)
	studentHandler := handlers.NewStudentHandler(studentSvc)
	loanHandler := handlers.NewLoanHandler(loanSvc)
	fineHandler := handlers.NewFineHandler(fineSvc)
	sanctionHandler := handlers.NewSanctionHandler(sanctionSvc)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsSvc)
	rankingHandler := handlers.NewRankingHandler(rankingSvc)
	badgeDefSvc := badgedef.NewService(badgeRepo)
	badgeAdminHandler := handlers.NewBadgeAdminHandler(badgeDefSvc)
	staffAPI := handlers.StaffBearerMiddleware(tokenSvc, userRepo)

	startBlacklistCleanup(revokedRepo)
	startLoanOverdueSync(loanRepo)
	startGlobalRankSync(orm)

	seedHandler := handlers.NewSeedHandler(orm)

	r := gin.Default()
	r.Use(handlers.AllowAllCORS())
	rlCfg := ratelimit.ConfigFromEnv()
	rlStore := ratelimit.NewStore(rlCfg)
	rlStore.StartSweeper(5*time.Minute, 30*time.Minute)
	if rlCfg.Enabled {
		log.Printf("rate limit: enabled (general %d req/min per IP, ranking %d req/min per IP; RATE_LIMIT_DISABLED=1 to disable)", rlCfg.GeneralRPM, rlCfg.RankingRPM)
	}
	r.Use(handlers.IPRateLimit(rlStore))

	r.GET("/health", healthHandler.Check)
	r.POST("/setup/seed", seedHandler.Run)
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/refresh", authHandler.Refresh)
	r.POST("/auth/logout", authHandler.Logout)
	r.GET("/auth/me", handlers.AuthBearerMiddleware(tokenSvc), authHandler.Me)

	publicAPI := r.Group("/api")
	{
		publicAPI.GET("/ranking/students/:studentId", rankingHandler.StudentRank)
		publicAPI.GET("/ranking/top3", rankingHandler.Top3)
		publicAPI.GET("/ranking/leaderboard", rankingHandler.Leaderboard)
		publicAPI.GET("/books", bookHandler.List)
		publicAPI.GET("/books/:id", bookHandler.Get)
		publicAPI.GET("/authors", catalogListHandler.ListAuthors)
		publicAPI.GET("/genres", catalogListHandler.ListGenres)
		publicAPI.GET("/students", studentHandler.List)
	}

	api := r.Group("/api", staffAPI)
	{
		api.GET("/departments", departmentHandler.List)

		api.POST("/students", studentHandler.Create)
		api.GET("/students/:id/profile", studentHandler.Profile)
		api.POST("/students/:id/badges", studentHandler.AwardBadge)
		api.DELETE("/students/:id/badges/:badgeId", studentHandler.RevokeBadge)
		api.GET("/students/:id", studentHandler.Get)
		api.PATCH("/students/:id", studentHandler.Update)
		api.DELETE("/students/:id", studentHandler.Deactivate)

		api.POST("/authors", authorHandler.Create)

		api.POST("/books", bookHandler.Create)
		api.PUT("/books/:id", bookHandler.UpdatePut)
		api.PATCH("/books/:id", bookHandler.UpdatePatch)
		api.DELETE("/books/:id", bookHandler.Delete)

		api.GET("/badges", badgeAdminHandler.List)
		api.POST("/badges", badgeAdminHandler.Create)
		api.GET("/badges/:id", badgeAdminHandler.Get)
		api.PATCH("/badges/:id", badgeAdminHandler.Patch)
		api.DELETE("/badges/:id", badgeAdminHandler.Delete)

		api.GET("/loans", loanHandler.List)
		api.POST("/loans", loanHandler.Create)
		api.PATCH("/loans/:id/return", loanHandler.Return)

		api.GET("/fines", fineHandler.List)
		api.POST("/fines", fineHandler.Create)
		api.GET("/fines/:id", fineHandler.Get)
		api.PATCH("/fines/:id/paid", fineHandler.MarkPaid)
		api.PATCH("/fines/:id/waived", fineHandler.MarkWaived)

		api.GET("/sanctions", sanctionHandler.List)
		api.POST("/sanctions", sanctionHandler.Create)
		api.PATCH("/sanctions/:id", sanctionHandler.Lift)

		api.GET("/analytics/dashboard", analyticsHandler.Dashboard)
	}

	addr := ":8080"
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		if !strings.HasPrefix(p, ":") {
			p = ":" + p
		}
		addr = p
	}
	log.Printf("http: listening on %s", addr)
	_ = r.Run(addr)
}

func startBlacklistCleanup(repo repositories.RevokedTokenRepository) {
	job := &jobs.BlacklistCleanup{Repo: repo}
	go func() {
		ticker := time.NewTicker(48 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			n, err := job.Run(context.Background())
			if err != nil {
				log.Printf("blacklist cleanup: %v", err)
				continue
			}
			if n > 0 {
				log.Printf("blacklist cleanup: removed %d rows", n)
			}
		}
	}()
}

// startGlobalRankSync persists student_stats.global_rank from leaderboard ordering (active students).
func startGlobalRankSync(orm *gorm.DB) {
	job := &jobs.GlobalRankSync{DB: orm}
	run := func() {
		n, err := job.Run(context.Background())
		if err != nil {
			log.Printf("global rank sync: %v", err)
			return
		}
		if n > 0 {
			log.Printf("global rank sync: updated %d row(s)", n)
		}
	}
	run()
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			run()
		}
	}()
}

// startLoanOverdueSync runs a DB sync job: active loans with due_date < today become overdue.
// Uses PostgreSQL CURRENT_DATE in the repository UPDATE. Interval is conservative for load;
// adjust if you need near-real-time mora (see functionalities/loan/design.md).
func startLoanOverdueSync(loans repositories.LoanRepository) {
	job := &jobs.LoanOverdueSync{Loans: loans}
	run := func() {
		n, err := job.Run(context.Background())
		if err != nil {
			log.Printf("loan overdue sync: %v", err)
			return
		}
		if n > 0 {
			log.Printf("loan overdue sync: marked %d loan(s) as overdue", n)
		}
	}
	run()
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			run()
		}
	}()
}
