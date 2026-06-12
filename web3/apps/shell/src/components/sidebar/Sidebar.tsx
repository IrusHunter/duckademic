import { NavLink } from 'react-router-dom'
import { LuHouse, LuBookOpen } from 'react-icons/lu'
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

// Пункти навігації між модулями (стилістика — з оригінального web).
// Показуються лише дозволені поточній ролі.
// Щоб додати модуль — підключіть його як remote у shell/vite.config.ts,
// додайте <Route> в App.tsx і допишіть пункт сюди.
const NAV: NavItem[] = [
  { to: '/home', label: 'Головна', icon: <LuHouse />, end: true, roles: ['student', 'teacher'] },
  { to: '/classroom', label: 'Заняття', icon: <LuBookOpen />, roles: ['student', 'teacher'] },
  // Приклади з web — розкоментуйте, коли відповідний remote підключено:
  // { to: '/courses',   label: 'Мої курси',    icon: <LuBookCopy />,      roles: ['student', 'teacher'] },
  // { to: '/messaging', label: 'Повідомлення', icon: <LuMessageSquare />, roles: ['student', 'teacher'] },
  // { to: '/schedule',  label: 'Розклад',      icon: <LuCalendar />,      roles: ['student', 'teacher'] },
  // { to: '/grades',    label: 'Оцінки',       icon: <LuGraduationCap />, roles: ['student'] },
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
