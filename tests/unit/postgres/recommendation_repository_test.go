package postgres_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/project-go-sender-recommendation-agent/internal/adapters/secondary/postgres"
	"github.com/project-go-sender-recommendation-agent/internal/core/domain"
)

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db, mock
}

// --- Save ---

func TestRecommendationRepository_Save_Success(t *testing.T) {
	db, mock := newMockDB(t)
	repo := postgres.NewRecommendationRepository(db)

	mock.ExpectExec(`INSERT INTO recommendations`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rec := &domain.Recommendation{ID: "1", SenderID: "s1", Payload: "data", Score: 0.9}
	if err := repo.Save(rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRecommendationRepository_Save_DBError(t *testing.T) {
	db, mock := newMockDB(t)
	repo := postgres.NewRecommendationRepository(db)

	mock.ExpectExec(`INSERT INTO recommendations`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("unique violation"))

	err := repo.Save(&domain.Recommendation{ID: "1", SenderID: "s1"})
	if err == nil {
		t.Fatal("expected error from db")
	}
}

// --- FindByID ---

func TestRecommendationRepository_FindByID_Found(t *testing.T) {
	db, mock := newMockDB(t)
	repo := postgres.NewRecommendationRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "sender_id", "payload", "score", "created_at"}).
		AddRow("1", "s1", "data", 0.9, now)

	mock.ExpectQuery(`SELECT id, sender_id, payload, score, created_at FROM recommendations WHERE id`).
		WithArgs("1").
		WillReturnRows(rows)

	got, err := repo.FindByID("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected a record, got nil")
	}
	if got.ID != "1" || got.SenderID != "s1" {
		t.Fatalf("unexpected record: %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRecommendationRepository_FindByID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	repo := postgres.NewRecommendationRepository(db)

	mock.ExpectQuery(`SELECT id, sender_id, payload, score, created_at FROM recommendations WHERE id`).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	got, err := repo.FindByID("missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatal("expected nil for missing record")
	}
}

func TestRecommendationRepository_FindByID_DBError(t *testing.T) {
	db, mock := newMockDB(t)
	repo := postgres.NewRecommendationRepository(db)

	mock.ExpectQuery(`SELECT id, sender_id, payload, score, created_at FROM recommendations WHERE id`).
		WithArgs("1").
		WillReturnError(errors.New("connection error"))

	if _, err := repo.FindByID("1"); err == nil {
		t.Fatal("expected db error to propagate")
	}
}

// --- FindAll ---

func TestRecommendationRepository_FindAll_Success(t *testing.T) {
	db, mock := newMockDB(t)
	repo := postgres.NewRecommendationRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "sender_id", "payload", "score", "created_at"}).
		AddRow("1", "s1", "data1", 0.9, now).
		AddRow("2", "s2", "data2", 0.8, now)

	mock.ExpectQuery(`SELECT id, sender_id, payload, score, created_at FROM recommendations ORDER BY`).
		WillReturnRows(rows)

	got, err := repo.FindAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRecommendationRepository_FindAll_Empty(t *testing.T) {
	db, mock := newMockDB(t)
	repo := postgres.NewRecommendationRepository(db)

	mock.ExpectQuery(`SELECT id, sender_id, payload, score, created_at FROM recommendations ORDER BY`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sender_id", "payload", "score", "created_at"}))

	got, err := repo.FindAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 records, got %d", len(got))
	}
}

func TestRecommendationRepository_FindAll_DBError(t *testing.T) {
	db, mock := newMockDB(t)
	repo := postgres.NewRecommendationRepository(db)

	mock.ExpectQuery(`SELECT id, sender_id, payload, score, created_at FROM recommendations ORDER BY`).
		WillReturnError(errors.New("db down"))

	if _, err := repo.FindAll(); err == nil {
		t.Fatal("expected db error to propagate")
	}
}
