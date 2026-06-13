import { NavLink } from 'react-router-dom'
import { LuHouse, LuBookOpen, LuMessageSquare, LuCalendarDays, LuGraduationCap } from 'react-icons/lu'
import type { ReactNode } from 'react'
import type { User } from '../../store/authStore'
import css from './Sidebar.module.css'

type Role = User['role']

interface NavItem {
  to: string
  label: string
  icon: ReactNode
  end?: boolean
  roles: Role[]
}

const NAV: NavItem[] = [
  { to: '/home', label: 'Home', icon: <LuHouse />, end: true, roles: ['student', 'teacher'] },
  { to: '/classroom', label: 'My Courses', icon: <LuBookOpen />, roles: ['student', 'teacher'] },
  { to: '/messenger', label: 'Messaging', icon: <LuMessageSquare />, roles: ['student', 'teacher'] },
  { to: '/schedule', label: 'Schedule', icon: <LuCalendarDays />, roles: ['student', 'teacher'] },
  { to: '/grades',   label: 'Grades',   icon: <LuGraduationCap />,      roles: ['student', 'teacher'] },
]

interface SidebarProps {
  role: Role
}

export default function Sidebar({ role }: SidebarProps) {
  const items = NAV.filter((item) => item.roles.includes(role))

  return (
    <aside className={css.sidebar}>
      <nav>
        <ul className={css.list}>
          {items.map((item) => (
            <li className={css.item} key={item.to}>
              <NavLink
                to={item.to}
                end={item.end}
                className={({ isActive }) =>
                  `${css.link} ${isActive ? css.active : ''}`
                }
              >
                <span className={css.icon}>{item.icon}</span>
                {item.label}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  )
}
