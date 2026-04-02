package entities

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type ScheduleGeneratorConfig struct {
	StartDate           time.Time   `json:"start_date"`
	EndDate             time.Time   `json:"end_date"`
	SlotPreference      [][]float32 `json:"slot_preference"`
	MaxDailyStudentLoad int         `json:"max_daily_student_load"`
	LessonFillRate      float64     `json:"lesson_fill_rate"`
	ClassroomOccupancy  float64     `json:"classroom_occupancy"`
}

func (cfg ScheduleGeneratorConfig) String() string {
	parts := make([]string, 0, 6)

	parts = append(parts, fmt.Sprintf("start_time: %s", cfg.StartDate.Format("15:04"))) // формат HH:MM
	parts = append(parts, fmt.Sprintf("end_time: %s", cfg.EndDate.Format("15:04")))
	parts = append(parts, fmt.Sprintf("slot_preference: %v", cfg.SlotPreference))
	parts = append(parts, fmt.Sprintf("max_daily_student_load: %d", cfg.MaxDailyStudentLoad))
	parts = append(parts, fmt.Sprintf("lesson_fill_rate: %.2f%%", cfg.LessonFillRate*100))
	parts = append(parts, fmt.Sprintf("classroom_occupancy: %.2f%%", cfg.ClassroomOccupancy*100))

	return fmt.Sprintf("ScheduleGeneratorConfig{%s}", strings.Join(parts, ", "))
}

func (cfg *ScheduleGeneratorConfig) ValidateStartTime() error {
	if cfg.StartDate.IsZero() {
		return errors.New("start_date must not be zero")
	}
	return nil
}

func (cfg *ScheduleGeneratorConfig) ValidateEndTime() error {
	if cfg.EndDate.IsZero() {
		return errors.New("end_date must not be zero")
	}
	if !cfg.EndDate.After(cfg.StartDate) {
		return errors.New("end_date must be after start_date")
	}
	return nil
}

func (cfg *ScheduleGeneratorConfig) ValidateSlotPreference() error {
	if len(cfg.SlotPreference) != 7 {
		return fmt.Errorf("slot_preference must have 7 arrays (got %d)", len(cfg.SlotPreference))
	}
	for dayIdx, day := range cfg.SlotPreference {
		for slotIdx, val := range day {
			if val <= 0 || val >= 10 {
				return fmt.Errorf("slot_preference[%d][%d] must be >0 and <10 (got %f)", dayIdx, slotIdx, val)
			}
		}
	}
	return nil
}

func (cfg *ScheduleGeneratorConfig) ValidateMaxDailyStudentLoad() error {
	if cfg.MaxDailyStudentLoad <= 0 {
		return errors.New("max_daily_student_load must be positive")
	}
	return nil
}

func (cfg *ScheduleGeneratorConfig) ValidateLessonFillRate() error {
	if cfg.LessonFillRate <= 0 {
		return errors.New("lesson_fill_rate must be positive")
	}
	return nil
}

func (cfg *ScheduleGeneratorConfig) ValidateClassroomOccupancy() error {
	if cfg.ClassroomOccupancy <= 0 {
		return errors.New("classroom_occupancy must be positive")
	}
	return nil
}
