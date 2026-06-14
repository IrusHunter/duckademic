declare module 'authApp/AuthApp' {
  type AuthLoginResult = {
    id: string
    login: string
    role: 'admin' | 'student' | 'teacher'
    is_default_password: boolean
    access_token: string
    refresh_token: string
  }

  type Props = {
    onLoginSuccess?: (data: AuthLoginResult) => void
  }

  const AuthApp: import('react').ComponentType<Props>
  export default AuthApp
}

declare module 'classroomApp/ClassroomApp' {
  const ClassroomApp: import('react').ComponentType
  export default ClassroomApp
}

declare module 'homeApp/HomeApp' {
  const HomeApp: import('react').ComponentType
  export default HomeApp
}

declare module 'adminApp/AdminApp' {
  const AdminApp: import('react').ComponentType
  export default AdminApp
}

declare module 'dashboardApp/DashboardApp' {
  const DashboardApp: import('react').ComponentType
  export default DashboardApp
}

declare module 'gradesApp/GradesApp' {
  const GradesApp: import('react').ComponentType
  export default GradesApp
}

declare module 'messengerApp/MessengerApp' {
  const MessengerApp: import('react').ComponentType
  export default MessengerApp
}

declare module 'scheduleApp/ScheduleApp' {
  const ScheduleApp: import('react').ComponentType
  export default ScheduleApp
}