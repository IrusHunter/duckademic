import { useState } from 'react'
import { useNavigate, useParams, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider, useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'
import css from './App.module.css'
import { LuUser, LuBookCopy, LuGraduationCap, LuUsers, LuNotebookText, LuUniversity, LuSheet } from 'react-icons/lu';

const queryClient = new QueryClient()

// ─── ТИПИ ────────────────────────────────────────────────────────────────────

type FieldDef = {
  key: string
  label: string
  required?: boolean
  format?: 'time-ns' | 'duration-ns' | 'weekday-ua'
  relation?: {
    serviceKey: string
    tableKey: string
    labelKey: string
  }
}

type TableDef = {
  key: string
  label: string
  listEndpoint: string
  itemEndpoint: string
  fields: FieldDef[]
  editFields?: FieldDef[]   // окремі поля для редагування (якщо відрізняються від fields)
  columns: FieldDef[]
  readOnly?: boolean
  numericKeys?: string[]
}

type ServiceDef = {
  key: string
  label: string
  icon: React.ReactNode | string
  baseURL: string
  tables: TableDef[]
}

// ─── КОНФІГ СЕРВІСІВ ─────────────────────────────────────────────────────────

const SERVICES: ServiceDef[] = [
  // ── EMPLOYEE ──────────────────────────────────────────────────────────────
  {
    key: 'employee',
    label: 'Employee Service',
    icon: <LuUser />,
    baseURL: '/api/employee',
    tables: [
      {
        key: 'employees',
        label: 'Employees',
        listEndpoint: '/employees',
        itemEndpoint: '/employee',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'first_name', label: 'First Name' },
          { key: 'last_name', label: 'Last Name' },
          { key: 'middle_name', label: 'Middle Name' },
          { key: 'phone_number', label: 'Phone' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
          { key: 'deleted_at', label: 'Deleted At' },
        ],
        fields: [
          { key: 'first_name', label: 'First Name', required: true },
          { key: 'last_name', label: 'Last Name', required: true },
          { key: 'middle_name', label: 'Middle Name' },
          { key: 'phone_number', label: 'Phone Number' },
        ],
      },
      {
        key: 'teachers',
        label: 'Teachers',
        listEndpoint: '/teachers',
        itemEndpoint: '/teacher',
        columns: [
          { key: 'employee_id', label: 'Employee ID' },
          { key: 'email', label: 'Email' },
          { key: 'academic_rank_id', label: 'Academic Rank ID' },
          { key: 'academic_degree_id', label: 'Academic Degree ID' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
          { key: 'deleted_at', label: 'Deleted At' },
        ],
        fields: [
          {
            key: 'employee_id', label: 'Employee', required: true,
            relation: { serviceKey: 'employee', tableKey: 'employees', labelKey: 'last_name' },
          },
          { key: 'email', label: 'Email', required: true },
          {
            key: 'academic_rank_id', label: 'Academic Rank', required: true,
            relation: { serviceKey: 'employee', tableKey: 'academic-ranks', labelKey: 'title' },
          },
          {
            key: 'academic_degree_id', label: 'Academic Degree',
            relation: { serviceKey: 'employee', tableKey: 'academic-degrees', labelKey: 'title' },
          },
        ],
      },
      {
        key: 'academic-ranks',
        label: 'Academic Ranks',
        listEndpoint: '/academic-ranks',
        itemEndpoint: '/academic-rank',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'title', label: 'Title' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [{ key: 'title', label: 'Title', required: true }],
      },
      {
        key: 'academic-degrees',
        label: 'Academic Degrees',
        listEndpoint: '/academic-degrees',
        itemEndpoint: '/academic-degree',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'title', label: 'Title' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [{ key: 'title', label: 'Title', required: true }],
      },
    ],
  },

  // ── CURRICULUM ────────────────────────────────────────────────────────────
  {
    key: 'curriculum',
    label: 'Curriculum Service',
    icon: <LuBookCopy />,
    baseURL: '/api/curriculum',
    tables: [
      {
        key: 'curriculums',
        label: 'Curriculums',
        listEndpoint: '/curriculums',
        itemEndpoint: '/curriculum',
        numericKeys: ['duration_years'],
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'duration_years', label: 'Duration (years)' },
          { key: 'effective_from', label: 'Effective From' },
          { key: 'effective_to', label: 'Effective To' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
          { key: 'duration_years', label: 'Duration (years)', required: true },
          { key: 'effective_from', label: 'Effective From (YYYY-MM-DD)', required: true },
          { key: 'effective_to', label: 'Effective To (YYYY-MM-DD)' },
        ],
      },
      {
        key: 'semesters',
        label: 'Semesters',
        listEndpoint: '/semesters',
        itemEndpoint: '/semester',
        numericKeys: ['number'],
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'curriculum_id', label: 'Curriculum ID' },
          { key: 'number', label: 'Number' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          {
            key: 'curriculum_id', label: 'Curriculum', required: true,
            relation: { serviceKey: 'curriculum', tableKey: 'curriculums', labelKey: 'name' },
          },
          { key: 'number', label: 'Semester Number', required: true },
        ],
      },
      {
        key: 'disciplines',
        label: 'Disciplines',
        listEndpoint: '/disciplines',
        itemEndpoint: '/discipline',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [{ key: 'name', label: 'Name', required: true }],
      },
      {
        key: 'lesson-types',
        label: 'Lesson Types',
        listEndpoint: '/lesson-types',
        itemEndpoint: '/lesson-type',
        numericKeys: ['hours_value'],
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'hours_value', label: 'Hours per lesson' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
          { key: 'hours_value', label: 'Hours per lesson', required: true },
        ],
      },
      {
        key: 'semester-disciplines',
        label: 'Semester Disciplines',
        listEndpoint: '/semester-disciplines',
        itemEndpoint: '/semester-discipline',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'semester_id', label: 'Semester ID' },
          { key: 'discipline_id', label: 'Discipline ID' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          {
            key: 'semester_id', label: 'Semester', required: true,
            relation: { serviceKey: 'curriculum', tableKey: 'semesters', labelKey: 'number' },
          },
          {
            key: 'discipline_id', label: 'Discipline', required: true,
            relation: { serviceKey: 'curriculum', tableKey: 'disciplines', labelKey: 'name' },
          },
        ],
      },
      {
        key: 'lesson-type-assignments',
        label: 'Lesson Type Assignments',
        listEndpoint: '/lesson-type-assignments',
        itemEndpoint: '/lesson-type-assignment',
        numericKeys: ['required_hours'],
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'lesson_type_id', label: 'Lesson Type ID' },
          { key: 'discipline_id', label: 'Discipline ID' },
          { key: 'required_hours', label: 'Required Hours' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          {
            key: 'lesson_type_id', label: 'Lesson Type', required: true,
            relation: { serviceKey: 'curriculum', tableKey: 'lesson-types', labelKey: 'name' },
          },
          {
            key: 'discipline_id', label: 'Discipline', required: true,
            relation: { serviceKey: 'curriculum', tableKey: 'disciplines', labelKey: 'name' },
          },
          { key: 'required_hours', label: 'Required Hours', required: true },
        ],
      },
    ],
  },

  // ── STUDENT ───────────────────────────────────────────────────────────────
  {
    key: 'student',
    label: 'Student Service',
    icon: <LuGraduationCap />,
    baseURL: '/api/student',
    tables: [
      {
        key: 'students',
        label: 'Students',
        listEndpoint: '/students',
        itemEndpoint: '/student',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'first_name', label: 'First Name' },
          { key: 'last_name', label: 'Last Name' },
          { key: 'middle_name', label: 'Middle Name' },
          { key: 'email', label: 'Email' },
          { key: 'phone_number', label: 'Phone' },
          { key: 'semester_id', label: 'Semester ID' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
          { key: 'deleted_at', label: 'Deleted At' },
        ],
        fields: [
          { key: 'first_name', label: 'First Name', required: true },
          { key: 'last_name', label: 'Last Name', required: true },
          { key: 'email', label: 'Email', required: true },
          {
            key: 'semester_id', label: 'Semester', required: true,
            relation: { serviceKey: 'curriculum', tableKey: 'semesters', labelKey: 'number' },
          },
          { key: 'middle_name', label: 'Middle Name' },
          { key: 'phone_number', label: 'Phone Number' },
        ],
      },
      {
        key: 'semesters',
        label: 'Semesters (read)',
        listEndpoint: '/semesters',
        itemEndpoint: '/semester',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'curriculum_id', label: 'Curriculum ID' },
          { key: 'number', label: 'Number' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
    ],
  },

  // ── STUDENT GROUP ─────────────────────────────────────────────────────────
  {
    key: 'student-group',
    label: 'Student Group Service',
    icon: <LuUsers  />,
    baseURL: '/api/student-group',
    tables: [
      {
        key: 'group-cohorts',
        label: 'Group Cohorts',
        listEndpoint: '/group-cohorts',
        itemEndpoint: '/group-cohort',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'semester_id', label: 'Semester ID' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
          {
            key: 'semester_id', label: 'Semester', required: true,
            relation: { serviceKey: 'student-group', tableKey: 'semesters', labelKey: 'number' },
          },
        ],
      },
      {
        key: 'student-groups',
        label: 'Student Groups',
        listEndpoint: '/student-groups',
        itemEndpoint: '/student-group',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'group_cohort_id', label: 'Group Cohort ID' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          { key: 'name', label: 'Name', required: true },
          {
            key: 'group_cohort_id', label: 'Group Cohort', required: true,
            relation: { serviceKey: 'student-group', tableKey: 'group-cohorts', labelKey: 'name' },
          },
        ],
      },
      {
        key: 'group-members',
        label: 'Group Members',
        listEndpoint: '/group-members',
        itemEndpoint: '/group-member',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'studentId', label: 'Student ID' },
          { key: 'group_cohort_id', label: 'Group Cohort ID' },
          { key: 'student_group_id', label: 'Student Group ID' },
          { key: 'createdAt', label: 'Created At' },
          { key: 'updatedAt', label: 'Updated At' },
        ],
        fields: [
          {
            key: 'studentId', label: 'Student', required: true,
            relation: { serviceKey: 'student-group', tableKey: 'students', labelKey: 'name' },
          },
          {
            key: 'student_group_id', label: 'Student Group', required: true,
            relation: { serviceKey: 'student-group', tableKey: 'student-groups', labelKey: 'name' },
          },
        ],
      },
      {
        key: 'group-cohort-assignments',
        label: 'Cohort Assignments',
        listEndpoint: '/group-cohort-assignments',
        itemEndpoint: '/group-cohort-assignment',
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'group_cohort_id', label: 'Group Cohort ID' },
          { key: 'discipline_id', label: 'Discipline ID' },
          { key: 'lesson_type_id', label: 'Lesson Type ID' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          {
            key: 'group_cohort_id', label: 'Group Cohort', required: true,
            relation: { serviceKey: 'student-group', tableKey: 'group-cohorts', labelKey: 'name' },
          },
          {
            key: 'discipline_id', label: 'Discipline', required: true,
            relation: { serviceKey: 'student-group', tableKey: 'disciplines', labelKey: 'name' },
          },
          {
            key: 'lesson_type_id', label: 'Lesson Type', required: true,
            relation: { serviceKey: 'student-group', tableKey: 'lesson-types', labelKey: 'name' },
          },
        ],
      },
      {
        key: 'semesters',
        label: 'Semesters (read)',
        listEndpoint: '/semesters',
        itemEndpoint: '/semester',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'curriculum_id', label: 'Curriculum ID' },
          { key: 'number', label: 'Number' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'students',
        label: 'Students (read)',
        listEndpoint: '/students',
        itemEndpoint: '/student',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'semester_id', label: 'Semester ID' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'disciplines',
        label: 'Disciplines (read)',
        listEndpoint: '/disciplines',
        itemEndpoint: '/discipline',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'lesson-types',
        label: 'Lesson Types (read)',
        listEndpoint: '/lesson-types',
        itemEndpoint: '/lesson-type',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
    ],
  },

  // ── TEACHER LOAD ──────────────────────────────────────────────────────────
  {
    key: 'teacher-load',
    label: 'Teacher Load Service',
    icon: <LuNotebookText />,
    baseURL: '/api/teacher-load',
    tables: [
      {
        key: 'teacher-loads',
        label: 'Teacher Loads',
        listEndpoint: '/teacher-loads',
        itemEndpoint: '/teacher-load',
        numericKeys: ['group_count'],
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'teacher_id', label: 'Teacher ID' },
          { key: 'discipline_id', label: 'Discipline ID' },
          { key: 'lesson_type_id', label: 'Lesson Type ID' },
          { key: 'group_count', label: 'Group Count' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          {
            key: 'teacher_id', label: 'Teacher', required: true,
            relation: { serviceKey: 'teacher-load', tableKey: 'teachers', labelKey: 'name' },
          },
          {
            key: 'discipline_id', label: 'Discipline', required: true,
            relation: { serviceKey: 'teacher-load', tableKey: 'disciplines', labelKey: 'name' },
          },
          {
            key: 'lesson_type_id', label: 'Lesson Type', required: true,
            relation: { serviceKey: 'teacher-load', tableKey: 'lesson-types', labelKey: 'name' },
          },
          { key: 'group_count', label: 'Group Count', required: true },
        ],
      },
      {
        key: 'teachers',
        label: 'Teachers (read)',
        listEndpoint: '/teachers',
        itemEndpoint: '/teacher',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'disciplines',
        label: 'Disciplines (read)',
        listEndpoint: '/disciplines',
        itemEndpoint: '/discipline',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'lesson-types',
        label: 'Lesson Types (read)',
        listEndpoint: '/lesson-types',
        itemEndpoint: '/lesson-type',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
    ],
  },

  // ── ASSET ─────────────────────────────────────────────────────────────────
  {
    key: 'asset',
    label: 'Asset Service',
    icon: <LuUniversity />,
    baseURL: '/api/asset',
    tables: [
      {
        key: 'classrooms',
        label: 'Classrooms',
        listEndpoint: '/classrooms',
        itemEndpoint: '/classroom',
        numericKeys: ['capacity'],
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'number', label: 'Number' },
          { key: 'capacity', label: 'Capacity' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [
          { key: 'number', label: 'Classroom Number', required: true },
          { key: 'capacity', label: 'Capacity', required: true },
        ],
      },
    ],
  },

  // ── SCHEDULE ──────────────────────────────────────────────────────────────
  {
    key: 'schedule',
    label: 'Schedule Service',
    icon: <LuSheet  />,
    baseURL: '/api/schedule',
    tables: [
      {
        key: 'lesson-slots',
        label: 'Lesson Slots',
        listEndpoint: '/lesson-slots',
        itemEndpoint: '/lesson-slot',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slot', label: 'Slot' },
          { key: 'weekday', label: 'Weekday', format: 'weekday-ua' },
          { key: 'start_time', label: 'Start Time', format: 'time-ns' },
          { key: 'duration', label: 'Duration', format: 'duration-ns' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'lesson-occurrences',
        label: 'Lesson Occurrences',
        listEndpoint: '/lesson-occurrences',
        itemEndpoint: '/lesson-occurrence',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'study_load_id', label: 'Study Load ID' },
          { key: 'teacher_id', label: 'Teacher ID' },
          { key: 'student_group_id', label: 'Student Group ID' },
          { key: 'lesson_slot_id', label: 'Lesson Slot ID' },
          { key: 'classroom_id', label: 'Classroom ID' },
          { key: 'date', label: 'Date' },
          { key: 'status', label: 'Status' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        // Schedule academic-ranks: тільки PUT /priority, не POST/DELETE
        // fields порожні → кнопки Add/Delete відсутні
        // editFields → лише priority → кнопка Edit є
        key: 'academic-ranks',
        label: 'Academic Ranks',
        listEndpoint: '/academic-ranks',
        itemEndpoint: '/academic-rank',
        readOnly: false,
        numericKeys: ['priority'],
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'title', label: 'Title' },
          { key: 'priority', label: 'Priority' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
        editFields: [
          { key: 'priority', label: 'Priority', required: true },
        ],
      },
      {
        key: 'lesson-types',
        label: 'Lesson Types (read)',
        listEndpoint: '/lesson-types',
        itemEndpoint: '/lesson-type',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'hours_value', label: 'Hours' },
          { key: 'reserved_weeks', label: 'Reserved Weeks' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'teachers',
        label: 'Teachers (read)',
        listEndpoint: '/teachers',
        itemEndpoint: '/teacher',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'academic_rank_id', label: 'Academic Rank ID' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'classrooms',
        label: 'Classrooms (read)',
        listEndpoint: '/classrooms',
        itemEndpoint: '/classroom',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'number', label: 'Number' },
          { key: 'capacity', label: 'Capacity' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
      {
        key: 'student-groups',
        label: 'Student Groups (read)',
        listEndpoint: '/student-groups',
        itemEndpoint: '/student-group',
        readOnly: true,
        columns: [
          { key: 'id', label: 'ID' },
          { key: 'slug', label: 'Slug' },
          { key: 'name', label: 'Name' },
          { key: 'created_at', label: 'Created At' },
          { key: 'updated_at', label: 'Updated At' },
        ],
        fields: [],
      },
    ],
  },
]

// ─── ХЕЛПЕРИ ─────────────────────────────────────────────────────────────────

function makeApi(baseURL: string) {
  return axios.create({ baseURL, withCredentials: true })
}

function getServiceByKey(key: string): ServiceDef | undefined {
  return SERVICES.find(s => s.key === key)
}

function normalizeArray(result: unknown): Record<string, unknown>[] {
  if (Array.isArray(result)) return result
  if (result && typeof result === 'object') {
    const obj = result as Record<string, unknown>
    if (Array.isArray(obj.data)) return obj.data as Record<string, unknown>[]
    if (Array.isArray(obj.items)) return obj.items as Record<string, unknown>[]
    if (Array.isArray(obj.results)) return obj.results as Record<string, unknown>[]
    return [obj]
  }
  return []
}

function nsToTime(ns: unknown): string {
  if (ns === null || ns === undefined || ns === '') return ''
  const totalSeconds = Math.floor(Number(ns) / 1_000_000_000)
  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}`
}

function nsToDuration(ns: unknown): string {
  if (ns === null || ns === undefined || ns === '') return ''
  const totalMinutes = Math.floor(Number(ns) / 60_000_000_000)
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60
  if (hours === 0) return `${minutes}хв`
  if (minutes === 0) return `${hours} год`
  return `${hours} год ${minutes} хв`
}

function nsToWeekday(n: unknown): string {
  const days: Record<number, string> = {
    1: 'Понеділок',
    2: 'Вівторок',
    3: 'Середа',
    4: 'Четвер',
    5: 'Пʼятниця',
    6: 'Субота',
    7: 'Неділя',
  }
  return days[Number(n)] ?? String(n ?? '')
}

function formatCell(value: unknown, format?: FieldDef['format']): string {
  if (format === 'time-ns') return nsToTime(value)
  if (format === 'duration-ns') return nsToDuration(value)
  if (format === 'weekday-ua') return nsToWeekday(value)
  return String(value ?? '')
}

// ─── RELATION SELECT ─────────────────────────────────────────────────────────

function RelationSelect({ field, value, onChange }: {
  field: FieldDef
  value: string
  onChange: (v: string) => void
}) {
  const rel = field.relation!
  const service = getServiceByKey(rel.serviceKey)
  const table = service?.tables.find(t => t.key === rel.tableKey)

  const { data: raw, isLoading } = useQuery({
    queryKey: ['relation', rel.serviceKey, rel.tableKey],
    queryFn: async () => {
      const api = makeApi(service!.baseURL)
      const r = await api.get(table!.listEndpoint)
      return normalizeArray(r.data)
    },
    enabled: !!service && !!table,
  })

  const options = raw ?? []

  return (
    <select
      value={value}
      onChange={e => onChange(e.target.value)}
      required={field.required}
      style={{ padding: '6px 8px', borderRadius: 4, border: '1px solid #ccc', width: 194 }}
    >
      <option value="">— select —</option>
      {isLoading && <option disabled>Loading...</option>}
      {options.map((item, idx) => {
        const id = String(item.id ?? idx)
        const label = String(item[rel.labelKey] ?? id)
        return (
          <option key={id || `opt-${idx}`} value={id}>{label}</option>
        )
      })}
    </select>
  )
}

// ─── ADD FORM ─────────────────────────────────────────────────────────────────

function AddForm({ fields, onSubmit }: {
  fields: FieldDef[]
  onSubmit: (data: Record<string, string>) => void
}) {
  const [values, setValues] = useState<Record<string, string>>({})
  const [open, setOpen] = useState(false)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(values)
    setValues({})
    setOpen(false)
  }

  return (
    <div style={{ marginBottom: 20 }}>
      <button
        onClick={() => setOpen(!open)}
        style={{ padding: '6px 16px', background: '#4A6CF7', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer', marginBottom: 10 }}
      >
        {open ? '✕ Cancel' : '+ Add'}
      </button>
      {open && (
        <form onSubmit={handleSubmit} style={{ display: 'flex', gap: 12, flexWrap: 'wrap', padding: 16, background: '#f9f9f9', borderRadius: 8 }}>
          {fields.map(field => (
            <div key={field.key} style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
              <label style={{ fontSize: 12, color: '#666' }}>{field.label}{field.required && ' *'}</label>
              {field.relation ? (
                <RelationSelect
                  field={field}
                  value={values[field.key] ?? ''}
                  onChange={v => setValues(prev => ({ ...prev, [field.key]: v }))}
                />
              ) : (
                <input
                  value={values[field.key] ?? ''}
                  onChange={e => setValues(prev => ({ ...prev, [field.key]: e.target.value }))}
                  required={field.required}
                  style={{ padding: '6px 8px', borderRadius: 4, border: '1px solid #ccc', width: 180 }}
                />
              )}
            </div>
          ))}
          <div style={{ display: 'flex', alignItems: 'flex-end' }}>
            <button type="submit" style={{ padding: '6px 16px', background: '#333', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
              Save
            </button>
          </div>
        </form>
      )}
    </div>
  )
}

// ─── DATA TABLE ───────────────────────────────────────────────────────────────

function DataTable({ data, columns, editFields, onDelete, onEditClick, readOnly, canDelete }: {
  data: Record<string, unknown>[]
  columns: FieldDef[]
  editFields: FieldDef[]
  onDelete: (id: string) => void
  onEditClick: (id: string) => void
  readOnly?: boolean
  canDelete: boolean
}) {
  const [selected, setSelected] = useState<string[]>([])
  const [action, setAction] = useState('')

  const toggleSelect = (id: string) =>
    setSelected(prev => prev.includes(id) ? prev.filter(i => i !== id) : [...prev, id])

  const handleGo = () => {
    if (action === 'delete') {
      selected.forEach(id => onDelete(id))
      setSelected([])
    }
  }

  // Показувати колонку Actions якщо є кнопка Edit або Delete
  const showActions = !readOnly && (editFields.length > 0 || canDelete)

  return (
    <div>
      {canDelete && (
        <div style={{ marginBottom: 12, display: 'flex', gap: 8, alignItems: 'center' }}>
          <span style={{ fontSize: 14 }}>Action:</span>
          <select value={action} onChange={e => setAction(e.target.value)} style={{ padding: '4px 8px', borderRadius: 4, border: '1px solid #ccc' }}>
            <option value="">---</option>
            <option value="delete">Delete selected</option>
          </select>
          <button onClick={handleGo} style={{ padding: '4px 12px', borderRadius: 4, border: '1px solid #ccc', cursor: 'pointer' }}>Go</button>
        </div>
      )}
      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: 14 }}>
          <thead>
            <tr style={{ background: '#f5f5f5' }}>
              {canDelete && <th style={{ padding: '8px 12px', borderBottom: '2px solid #e0e0e0', width: 40 }} />}
              {columns.map(col => (
                <th key={col.key} style={{ padding: '8px 12px', textAlign: 'left', borderBottom: '2px solid #e0e0e0', color: '#555', whiteSpace: 'nowrap' }}>
                  {col.label}
                </th>
              ))}
              {showActions && <th style={{ padding: '8px 12px', borderBottom: '2px solid #e0e0e0', color: '#555' }}>Actions</th>}
            </tr>
          </thead>
          <tbody>
            {data.length === 0 && (
              <tr>
                <td colSpan={columns.length + (showActions ? 1 : 0) + (canDelete ? 1 : 0)} style={{ padding: 24, textAlign: 'center', color: '#aaa' }}>
                  No data
                </td>
              </tr>
            )}
            {data.map((row, idx) => {
              const id = String(row.id ?? '')
              const rowKey = id || `row-${idx}`
              return (
                <tr key={rowKey} style={{ background: selected.includes(id) ? '#e8f0fe' : 'white' }}>
                  {canDelete && (
                    <td style={{ padding: '8px 12px', borderBottom: '1px solid #eee' }}>
                      <input type="checkbox" checked={selected.includes(id)} onChange={() => toggleSelect(id)} />
                    </td>
                  )}
                  {columns.map(col => (
                    <td key={col.key} style={{ padding: '8px 12px', borderBottom: '1px solid #eee', maxWidth: 220, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                      {formatCell(row[col.key], col.format)}
                    </td>
                  ))}
                  {showActions && (
                    <td style={{ padding: '8px 12px', borderBottom: '1px solid #eee', whiteSpace: 'nowrap' }}>
                      {editFields.length > 0 && (
                        <button
                          onClick={() => onEditClick(id)}
                          style={{ padding: '3px 10px', color: '#4A6CF7', border: '1px solid #4A6CF7', borderRadius: 4, background: 'white', cursor: 'pointer', fontSize: 13, marginRight: 6 }}
                        >
                          Edit
                        </button>
                      )}
                      {canDelete && (
                        <button
                          onClick={() => onDelete(id)}
                          style={{ padding: '3px 10px', color: 'red', border: '1px solid red', borderRadius: 4, background: 'white', cursor: 'pointer', fontSize: 13 }}
                        >
                          Delete
                        </button>
                      )}
                    </td>
                  )}
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}

// ─── EDIT PAGE ────────────────────────────────────────────────────────────────

function EditPage({ service, table }: { service: ServiceDef; table: TableDef }) {
  const navigate = useNavigate()
  const { itemId } = useParams<{ itemId: string }>()
  const qc = useQueryClient()
  const api = makeApi(service.baseURL)
  const queryKey = [service.key, table.key]

  const editFields = table.editFields ?? table.fields

  const { data: raw, isLoading } = useQuery({
    queryKey: [service.key, table.key, itemId],
    queryFn: async () => {
      const r = await api.get(`${table.itemEndpoint}/${itemId}`)
      return r.data as Record<string, unknown>
    },
    enabled: !!itemId,
  })

  const [values, setValues] = useState<Record<string, string>>({})
  const initialized = Object.keys(values).length > 0
  if (raw && !initialized) {
    const initial = Object.fromEntries(editFields.map(f => [f.key, String(raw[f.key] ?? '')]))
    setValues(initial)
  }

  const getErrorMessage = (e: any): string => e?.response?.data?.error || 'Unknown error'

  const updateMutation = useMutation({
    mutationFn: (body: Record<string, string>) => {
      const converted: Record<string, unknown> = { ...body }
      for (const key of table.numericKeys ?? []) {
        if (converted[key] !== undefined && converted[key] !== '') {
          converted[key] = Number(converted[key])
        }
      }
      return api.put(`${table.itemEndpoint}/${itemId}`, converted).then(r => r.data)
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey })
      navigate(-1)
    },
    onError: (e) => console.error('Update error:', getErrorMessage(e)),
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    updateMutation.mutate(values)
  }

  return (
    <div style={{ padding: 24, flex: 1, maxWidth: 600 }}>
      <button
        onClick={() => navigate(-1)}
        style={{ marginBottom: 20, padding: '6px 14px', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', background: 'white' }}
      >
        ← Back
      </button>

      <h2 style={{ marginBottom: 4 }}>Edit — {table.label}</h2>
      <p style={{ color: '#999', fontSize: 13, marginBottom: 24 }}>
        {service.label} · {service.baseURL}{table.itemEndpoint}/{itemId}
      </p>

      {isLoading && <p>Loading...</p>}

      {!isLoading && raw && (
        <>
          {/* Інформаційні (read-only) поля запису */}
          <div style={{ marginBottom: 24, padding: 16, background: '#f9f9f9', borderRadius: 8, fontSize: 13 }}>
            <p style={{ fontWeight: 600, marginBottom: 8, color: '#555' }}>Record info</p>
            {table.columns
              .filter(col => !editFields.find(f => f.key === col.key))
              .map(col => (
                <div key={col.key} style={{ display: 'flex', gap: 8, marginBottom: 4 }}>
                  <span style={{ color: '#999', minWidth: 140 }}>{col.label}:</span>
                  <span style={{ color: '#333' }}>{formatCell(raw[col.key], col.format)}</span>
                </div>
              ))
            }
          </div>

          {/* Форма редагування */}
          <form onSubmit={handleSubmit}>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16, marginBottom: 24 }}>
              {editFields.map(field => (
                <div key={field.key}>
                  <label style={{ display: 'block', fontSize: 13, fontWeight: 600, color: '#444', marginBottom: 6 }}>
                    {field.label}{field.required && ' *'}
                  </label>
                  {field.relation ? (
                    <RelationSelect
                      field={field}
                      value={values[field.key] ?? ''}
                      onChange={v => setValues(prev => ({ ...prev, [field.key]: v }))}
                    />
                  ) : (
                    <input
                      value={values[field.key] ?? ''}
                      onChange={e => setValues(prev => ({ ...prev, [field.key]: e.target.value }))}
                      required={field.required}
                      style={{ padding: '8px 12px', borderRadius: 4, border: '1px solid #ccc', width: '100%', boxSizing: 'border-box', fontSize: 14 }}
                    />
                  )}
                </div>
              ))}
            </div>

            <div style={{ display: 'flex', gap: 10 }}>
              <button
                type="submit"
                disabled={updateMutation.isPending}
                style={{ padding: '8px 24px', background: '#4A6CF7', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer', fontSize: 14 }}
              >
                {updateMutation.isPending ? 'Saving...' : '✓ Save'}
              </button>
              <button
                type="button"
                onClick={() => navigate(-1)}
                style={{ padding: '8px 24px', background: 'white', color: '#666', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', fontSize: 14 }}
              >
                Cancel
              </button>
            </div>

            {updateMutation.isError && (
              <p style={{ color: 'red', marginTop: 12, fontSize: 13 }}>
                Error: {(updateMutation.error as any)?.response?.data?.error || 'Unknown error'}
              </p>
            )}
          </form>
        </>
      )}
    </div>
  )
}

// ─── TABLE PAGE ───────────────────────────────────────────────────────────────

function TablePage({ service, table }: { service: ServiceDef; table: TableDef }) {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const api = makeApi(service.baseURL)
  const queryKey = [service.key, table.key]

  const { data: rawData, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const r = await api.get(table.listEndpoint)
      return normalizeArray(r.data)
    },
  })
  const data = rawData ?? []

  const getErrorMessage = (e: any): string => e?.response?.data?.error || 'Unknown error'

  const convertNumeric = (body: Record<string, string>): Record<string, unknown> => {
    const converted: Record<string, unknown> = { ...body }
    for (const key of table.numericKeys ?? []) {
      if (converted[key] !== undefined && converted[key] !== '') {
        converted[key] = Number(converted[key])
      }
    }
    return converted
  }

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.delete(`${table.itemEndpoint}/${id}`).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey }),
    onError: (e) => console.error('Delete error:', getErrorMessage(e)),
  })

  const createMutation = useMutation({
    mutationFn: (body: Record<string, string>) =>
      api.post(table.listEndpoint, convertNumeric(body)).then(r => r.data),
    onSuccess: () => qc.invalidateQueries({ queryKey }),
    onError: (e) => console.error('Create error:', getErrorMessage(e)),
  })

  const editFields = table.editFields ?? table.fields
  // Delete доступний лише якщо є fields (є POST → є й DELETE)
  const canDelete = !table.readOnly && table.fields.length > 0

  return (
    <div style={{ padding: 24, flex: 1 }}>
      <button
        onClick={() => navigate(`/admin/${service.key}`)}
        style={{ marginBottom: 16, padding: '6px 14px', border: '1px solid #ccc', borderRadius: 4, cursor: 'pointer', background: 'white' }}
      >
        ← Back
      </button>
      <h2 style={{ marginBottom: 4 }}>{table.label}</h2>
      <p style={{ color: '#999', fontSize: 13, marginBottom: 20 }}>
        {service.label} · {service.baseURL}{table.listEndpoint}
      </p>

      {!table.readOnly && table.fields.length > 0 && (
        <AddForm fields={table.fields} onSubmit={body => createMutation.mutate(body)} />
      )}

      {isLoading && <p>Loading...</p>}
      {error && <p style={{ color: 'red' }}>Error loading data</p>}
      {!isLoading && !error && (
        <DataTable
          data={data}
          columns={table.columns}
          editFields={editFields}
          onDelete={id => deleteMutation.mutate(id)}
          onEditClick={id => navigate(`/admin/${service.key}/${table.key}/edit/${id}`)}
          readOnly={table.readOnly}
          canDelete={canDelete}
        />
      )}
    </div>
  )
}

// ─── SERVICE HOME ─────────────────────────────────────────────────────────────

function ServiceHome({ service }: { service: ServiceDef }) {
  const navigate = useNavigate()
  return (
    <div style={{ padding: 24 }}>
      <h2 style={{ marginBottom: 4 }}>{service.label}</h2>
      <p style={{ color: '#999', fontSize: 13, marginBottom: 24 }}>{service.baseURL}</p>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 10, maxWidth: 400 }}>
        {service.tables.map(table => (
          <button
            key={table.key}
            onClick={() => navigate(`/admin/${service.key}/${table.key}`)}
            style={{ padding: '14px 18px', textAlign: 'left', border: '1px solid #e0e0e0', borderRadius: 8, cursor: 'pointer', background: 'white', fontSize: 15 }}
          >
            {table.label}
            {table.readOnly && (
              <span style={{ marginLeft: 8, fontSize: 11, color: '#aaa', background: '#f0f0f0', padding: '2px 6px', borderRadius: 4 }}>
                read-only
              </span>
            )}
          </button>
        ))}
      </div>
    </div>
  )
}

// ─── DYNAMIC ROUTE ────────────────────────────────────────────────────────────

function DynamicServiceRoute() {
  const { serviceKey, tableKey, itemId } = useParams<{ serviceKey: string; tableKey: string; itemId: string }>()
  const service = SERVICES.find(s => s.key === serviceKey)
  if (!service) return <div style={{ padding: 24 }}>Service not found</div>
  if (!tableKey) return <ServiceHome service={service} />
  const table = service.tables.find(t => t.key === tableKey)
  if (!table) return <div style={{ padding: 24 }}>Table not found</div>
  if (itemId) return <EditPage service={service} table={table} />
  return <TablePage service={service} table={table} />
}

// ─── SIDEBAR ──────────────────────────────────────────────────────────────────

function Sidebar() {
  const navigate = useNavigate()
  const currentPath = window.location.pathname

  return (
    <aside className={css.aside}>
      <p className={css.servicesHeader}>Services</p>
      <ul className={css.servicesList}>
        {SERVICES.map(service => {
          const isActive = currentPath.includes(`/admin/${service.key}`)
          return (
            <li key={service.key}>
              <button
                onClick={() => navigate(`/admin/${service.key}`)}
                className={`${css.button} ${isActive ? css.active : ''}`}
              >
                <span className={css.serviceIcon}>{service.icon}</span>
                <span className={css.serviceTitle}>{service.label}</span>
              </button>
            </li>
          )
        })}
      </ul>
    </aside>
  )
}

// ─── LAYOUT ───────────────────────────────────────────────────────────────────

function AdminLayout() {
  return (
    <div style={{ display: 'flex', minHeight: '100vh' }}>
      <Sidebar />
      <div style={{ flex: 1 }} className={css.mainDiv}>
        <Routes>
          <Route path="/" element={<ServiceHome service={SERVICES[0]} />} />
          <Route path=":serviceKey" element={<DynamicServiceRoute />} />
          <Route path=":serviceKey/:tableKey" element={<DynamicServiceRoute />} />
          <Route path=":serviceKey/:tableKey/edit/:itemId" element={<DynamicServiceRoute />} />
        </Routes>
      </div>
    </div>
  )
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AdminLayout />
    </QueryClientProvider>
  )
}