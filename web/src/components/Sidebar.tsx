import { NavLink } from "react-router-dom";

type NavItem = {
  to: string;
  label: string;
  iconHref: string;
  end?: boolean;
  iconClassName?: string;
};

const NAV: NavItem[] = [
  { to: "/", label: "Home", iconHref: "/img/icons.svg#icon-SVG", end: true },
  { to: "/courses", label: "My Courses", iconHref: "/img/icons.svg#icon-SVG-1", iconClassName: "icon" },
  { to: "/messaging", label: "Messaging", iconHref: "/img/icons.svg#icon-SVG-2" },
  { to: "/schedule", label: "Schedule", iconHref: "/img/icons.svg#icon-SVG-3" },
  { to: "/grades", label: "Grades", iconHref: "/img/icons.svg#icon-SVG-4" },
];

function Icon({ href, className }: { href: string; className?: string }) {
  return (
    <svg width="16" height="16" className={className} aria-hidden="true">
      <use href={href} />
    </svg>
  );
}

export default function Sidebar() {
  return (
    <aside className="sidebar-left">
      <nav className="sidebar-navigation">
        <ul className="sidebar-navigation-list">
          {NAV.map((item) => (
            <li className="sidebar-navigation-list-item" key={item.to}>
              <NavLink
                to={item.to}
                end={item.end}
                className={({ isActive }) => (isActive ? "highlighted" : "")}
              >
                <Icon href={item.iconHref} className={item.iconClassName} />
                {item.label}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
