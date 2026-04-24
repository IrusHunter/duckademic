package main

import (
	"context"
	"fmt"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/jmoiron/sqlx"
)

func Migrate(database *sqlx.DB) error {
	migrationsF := []func(*sqlx.Tx) error{
		teacherMigrations,
		studentMigrations,
		courseMigrations,
		studentCourseMigrations,
		teacherCourseMigrations,
	}

	for _, f := range migrationsF {
		tx, err := database.Beginx()
		if err != nil {
			return err
		}

		if err := f(tx); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func studentMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS students (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create students table: %w", err)
	}

	indexSlug := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_students_slug
	ON students (slug);
	`
	if _, err := tx.Exec(indexSlug); err != nil {
		return fmt.Errorf("failed to create students slug index: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "students"); err != nil {
		return fmt.Errorf("failed to create on update trigger for students: %w", err)
	}

	return nil
}
func teacherMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teachers (
		id UUID PRIMARY KEY,
		slug TEXT NOT NULL,
		name TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teachers table: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "teachers"); err != nil {
		return fmt.Errorf("failed to create on update trigger for teachers: %w", err)
	}

	return nil
}
func courseMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS courses (
		id UUID PRIMARY KEY,
		manager_id UUID,
		slug TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create courses table: %w", err)
	}

	indexManager := `
	CREATE INDEX IF NOT EXISTS idx_courses_manager_id
	ON courses (manager_id);
	`
	if _, err := tx.Exec(indexManager); err != nil {
		return fmt.Errorf("failed to create courses manager_id index: %w", err)
	}

	if err := db.EnsureForeignKeyTx(
		context.Background(),
		tx,
		"fk_courses_manager",
		"courses",
		"manager_id",
		"teachers",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create manager_id foreign key: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "courses"); err != nil {
		return fmt.Errorf("failed to create on update trigger for courses: %w", err)
	}

	return nil
}
func studentCourseMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS student_courses (
		id UUID PRIMARY KEY,
		student_id UUID NOT NULL,
		course_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create student_courses table: %w", err)
	}

	indexUnique := `
	CREATE UNIQUE INDEX IF NOT EXISTS uq_student_courses_student_course
	ON student_courses (student_id, course_id);
	`
	if _, err := tx.Exec(indexUnique); err != nil {
		return fmt.Errorf("failed to create student_courses unique index: %w", err)
	}

	if err := db.EnsureForeignKeyTx(
		context.Background(),
		tx,
		"fk_student_courses_student",
		"student_courses",
		"student_id",
		"students",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create student_id foreign key: %w", err)
	}

	if err := db.EnsureForeignKeyTx(
		context.Background(),
		tx,
		"fk_student_courses_course",
		"student_courses",
		"course_id",
		"courses",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create course_id foreign key: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "student_courses"); err != nil {
		return fmt.Errorf("failed to create on update trigger for student_courses: %w", err)
	}

	return nil
}
func teacherCourseMigrations(tx *sqlx.Tx) error {
	schema := `
	CREATE TABLE IF NOT EXISTS teacher_courses (
		id UUID PRIMARY KEY,
		teacher_id UUID NOT NULL,
		course_id UUID NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
	);
	`

	if _, err := tx.Exec(schema); err != nil {
		return fmt.Errorf("failed to create teacher_courses table: %w", err)
	}

	indexUnique := `
	CREATE UNIQUE INDEX IF NOT EXISTS uq_teacher_courses_teacher_course
	ON teacher_courses (teacher_id, course_id);
	`
	if _, err := tx.Exec(indexUnique); err != nil {
		return fmt.Errorf("failed to create teacher_courses unique index: %w", err)
	}

	if err := db.EnsureForeignKeyTx(
		context.Background(),
		tx,
		"fk_teacher_courses_teacher",
		"teacher_courses",
		"teacher_id",
		"teachers",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create teacher_id foreign key: %w", err)
	}

	if err := db.EnsureForeignKeyTx(
		context.Background(),
		tx,
		"fk_teacher_courses_course",
		"teacher_courses",
		"course_id",
		"courses",
		"id",
	); err != nil {
		return fmt.Errorf("failed to create course_id foreign key: %w", err)
	}

	if err := db.EnsureUpdatedAtTriggerTx(context.Background(), tx, "teacher_courses"); err != nil {
		return fmt.Errorf("failed to create updated_at trigger for teacher_courses: %w", err)
	}

	return nil
}
