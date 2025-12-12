package model

import (
	"time"

	"github.com/google/uuid"
)

// AchievementStatistics response untuk statistik prestasi (FR-011)
type AchievementStatistics struct {
	TotalByType        []TypeStatistic   `json:"total_by_type"`
	TotalByPeriod      []PeriodStatistic `json:"total_by_period"`
	TopStudents        []TopStudent      `json:"top_students"`
	CompetitionLevels  []LevelStatistic  `json:"competition_levels"`
	Summary            StatisticSummary  `json:"summary"`
}

// TypeStatistic statistik berdasarkan tipe prestasi
type TypeStatistic struct {
	AchievementType string `json:"achievement_type"`
	Count           int64  `json:"count"`
	Percentage      float64 `json:"percentage"`
}

// PeriodStatistic statistik berdasarkan periode (bulan/tahun)
type PeriodStatistic struct {
	Period string `json:"period"` // Format: "2024-01" atau "2024"
	Count  int64  `json:"count"`
	Year   int    `json:"year"`
	Month  *int   `json:"month,omitempty"` // null jika yearly
}

// TopStudent mahasiswa dengan prestasi terbanyak
type TopStudent struct {
	StudentID       uuid.UUID `json:"student_id"`
	StudentIDNumber string    `json:"student_id_number"`
	FullName        string    `json:"full_name"`
	ProgramStudy    string    `json:"program_study"`
	AcademicYear    string    `json:"academic_year"`
	TotalCount      int64     `json:"total_count"`
	VerifiedCount   int64     `json:"verified_count"`
}

// LevelStatistic statistik berdasarkan tingkat kompetisi
type LevelStatistic struct {
	Level      string  `json:"level"` // internasional, nasional, regional, lokal
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

// StatisticSummary ringkasan statistik
type StatisticSummary struct {
	TotalAchievements    int64     `json:"total_achievements"`
	VerifiedAchievements int64     `json:"verified_achievements"`
	PendingAchievements  int64     `json:"pending_achievements"`
	RejectedAchievements int64     `json:"rejected_achievements"`
	TotalStudents        int64     `json:"total_students"`
	DateRange            DateRange `json:"date_range"`
}

// DateRange rentang tanggal data
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// StatisticsRequest request untuk filter statistik
type StatisticsRequest struct {
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	AchievementType *string    `json:"achievement_type,omitempty"`
	Status          *string    `json:"status,omitempty"`
	PeriodType      string     `json:"period_type"` // "monthly" atau "yearly"
	TopLimit        int        `json:"top_limit"`   // limit untuk top students
}