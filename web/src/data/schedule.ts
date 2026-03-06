// src/data/schedule.ts

export type LessonType = "Lecture" | "Lab";
export type LessonColor = "blue" | "green" | "rose" | "amber";

export type Audience = {
  faculty: string;
  specialty: string;
  year: string;
  group: string;
  subgroup: string;
};

export type Lesson = {
  start: string;
  end: string;
  title: string;
  teacher: string;
  type: LessonType;
  modeText: string;
  dates: string;
  color: LessonColor;

  // важливо: можемо фільтрувати по частині полів (не обов’язково всі)
  audience?: Partial<Audience>;
};

export const dayOrder = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"] as const;
export type DayName = (typeof dayOrder)[number];

export const typeIconId: Record<LessonType, string> = {
  Lecture: "icon-video",
  Lab: "icon-navigation",
};

export const VIEW_OPTIONS = ["General", "Date"] as const;

export const FACULTY_OPTIONS = [
  "Faculty of Information Technology",
  "Faculty of History",
  "Faculty of Philosophy",
  "Faculty of Law",
  "Faculty of Economics",
  "Faculty of Psychology",
] as const;

export const SPECIALTIES_BY_FACULTY: Record<string, string[]> = {
  "Faculty of Information Technology": ["Software Engineering"],
  "Faculty of History": ["History", "Archaeology", "Museum and Heritage Studies"],
  "Faculty of Philosophy": ["Philosophy", "Linguistics", "Translation and Interpreting", "English Philology"],
  "Faculty of Law": ["Law", "International Law"],
  "Faculty of Economics": ["Economics"],
  "Faculty of Psychology": ["Public Relations and Communications", "Journalism and Media Studies"],
};

export const YEAR_OPTIONS = ["1 Bachelor", "2 Bachelor", "3 Bachelor", "4 Bachelor", "1 Master", "2 Master"] as const;
export const GROUP_OPTIONS = ["1", "2", "3", "4", "5", "6", "7"] as const;
export const SUBGROUP_OPTIONS = ["1", "2", "3", "4", "5", "6", "7"] as const;

// дефолт, як у тебе в макеті
export const DEFAULT_FILTERS: Audience & { view: (typeof VIEW_OPTIONS)[number] } = {
  view: "General",
  faculty: "Faculty of Information Technology",
  specialty: "Software Engineering",
  year: "3 Bachelor",
  group: "2",
  subgroup: "4",
};

export const scheduleData: Record<DayName, Lesson[]> = {
  Monday: [
    {
      start: "09:00",
      end: "10:20",
      title: "Introduction to Programming",
      teacher: "Dr. Smith",
      type: "Lecture",
      modeText: "Online",
      dates: "[23.01–28.05]",
      color: "blue",
      audience: { faculty: "Faculty of Information Technology", specialty: "Software Engineering", year: "3 Bachelor", group: "2", subgroup: "4" },
    },
    {
      start: "10:30",
      end: "11:50",
      title: "Data Science Lab",
      teacher: "Prof. Johnson",
      type: "Lab",
      modeText: "Room 101",
      dates: "[24.01–31.01, 14.02–01.05, 15.05–29.05]",
      color: "green",
      audience: { faculty: "Faculty of Information Technology", specialty: "Software Engineering", year: "3 Bachelor" },
    },
    {
      start: "12:10",
      end: "13:30",
      title: "Web Development Workshop",
      teacher: "Dr. Williams",
      type: "Lecture",
      modeText: "Online",
      dates: "[23.01, 06.02]",
      color: "rose",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "13:40",
      end: "15:00",
      title: "Study Group – Algorithms",
      teacher: "Student Led",
      type: "Lab",
      modeText: "Library Room 3",
      dates: "[24.01–31.01, 14.02–01.05, 15.05–29.05]",
      color: "amber",
      audience: { faculty: "Faculty of Information Technology" },
    },

    // приклад для іншого факультету — щоб було видно, що фільтр працює
    {
      start: "10:30",
      end: "11:50",
      title: "Ancient History Seminar",
      teacher: "Dr. Taylor",
      type: "Lecture",
      modeText: "Room 210",
      dates: "[24.01–30.05]",
      color: "blue",
      audience: { faculty: "Faculty of History", specialty: "History", year: "2 Bachelor" },
    },
  ],

  Tuesday: [
    {
      start: "08:30",
      end: "09:50",
      title: "Discrete Mathematics",
      teacher: "Dr. Brown",
      type: "Lecture",
      modeText: "Room 204",
      dates: "[24.01–30.05]",
      color: "blue",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "10:00",
      end: "11:20",
      title: "Operating Systems",
      teacher: "Dr. Miller",
      type: "Lecture",
      modeText: "Room 305",
      dates: "[24.01–30.05]",
      color: "rose",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "11:30",
      end: "13:00",
      title: "Operating Systems Lab",
      teacher: "Assistant Team",
      type: "Lab",
      modeText: "Online",
      dates: "[31.01–31.05]",
      color: "green",
      audience: { faculty: "Faculty of Information Technology", group: "2" },
    },
  ],

  Wednesday: [
    {
      start: "09:00",
      end: "10:20",
      title: "Database Systems",
      teacher: "Dr. Smith",
      type: "Lecture",
      modeText: "Room 110",
      dates: "[25.01–31.05]",
      color: "blue",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "10:30",
      end: "11:50",
      title: "Database Lab",
      teacher: "Prof. Johnson",
      type: "Lab",
      modeText: "Online",
      dates: "[01.02–31.05]",
      color: "green",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "12:10",
      end: "13:30",
      title: "Software Engineering",
      teacher: "Dr. Williams",
      type: "Lecture",
      modeText: "Room 207",
      dates: "[25.01–31.05]",
      color: "rose",
      audience: { faculty: "Faculty of Information Technology" },
    },
  ],

  Thursday: [
    {
      start: "09:00",
      end: "10:20",
      title: "Computer Networks",
      teacher: "Dr. Green",
      type: "Lecture",
      modeText: "Room 301",
      dates: "[26.01–01.06]",
      color: "blue",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "10:30",
      end: "11:50",
      title: "Computer Networks Lab",
      teacher: "Assistant Team",
      type: "Lab",
      modeText: "Online",
      dates: "[02.02–01.06]",
      color: "green",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "12:10",
      end: "13:30",
      title: "UI/UX Workshop",
      teacher: "Guest Lecturer",
      type: "Lecture",
      modeText: "Online",
      dates: "[09.02, 23.02, 09.03]",
      color: "amber",
      audience: { faculty: "Faculty of Information Technology" },
    },
  ],

  Friday: [
    {
      start: "08:30",
      end: "09:50",
      title: "Probability & Statistics",
      teacher: "Dr. Carter",
      type: "Lecture",
      modeText: "Room 210",
      dates: "[27.01–02.06]",
      color: "blue",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "10:00",
      end: "11:20",
      title: "Machine Learning Basics",
      teacher: "Prof. Johnson",
      type: "Lecture",
      modeText: "Online",
      dates: "[27.01–02.06]",
      color: "rose",
      audience: { faculty: "Faculty of Information Technology" },
    },
    {
      start: "11:30",
      end: "13:00",
      title: "Project Work – Duckademic",
      teacher: "Mentor Team",
      type: "Lab",
      modeText: "Project Room / Online",
      dates: "[03.02–02.06]",
      color: "amber",
      audience: { faculty: "Faculty of Information Technology" },
    },
  ],
};
