import { Link } from "react-router-dom";

export default function Header() {
  return (
    <header className="header">
      <div className="container">
        <nav className="header-navigation">
          <Link to="/dashboard" className="logo">
            <svg width="28" height="32" className="icon">
              <use href="/img/icons.svg#icon-Logo-1-1"></use>
            </svg>
            Duckademic
          </Link>

          <div className="search-bar">
            <svg width="23" height="23" className="icon icon-search">
              <use href="/img/icons.svg#icon-search-1"></use>
            </svg>
            <input type="text" placeholder="Search" className="search-bar-input" />
          </div>

          <div className="header-actions">
            <button className="notifications" type="button" aria-label="Notifications">
              <svg width="30" height="31" className="icon-bell">
                <use href="/img/icons.svg#icon-bell-1"></use>
              </svg>
              <svg width="9" height="9" className="icon icon-bell-has-new">
                <use href="/img/icons.svg#icon-Ellipse-1"></use>
              </svg>
            </button>

            <Link to="/dashboard" className="header-actions-profile" aria-label="Profile">
              <img src="/img/profile_pic.png" alt="profile" />
            </Link>
          </div>
        </nav>
      </div>
    </header>
  );
}
