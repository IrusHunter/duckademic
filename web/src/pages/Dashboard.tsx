import { useMemo, useState } from "react";

type ProfileData = {
  name: string;
  role: string;
  avatar: string;
  courses: number;
  assignments: number;
  groups: number;
};

type UpcomingEvent = {
  id: string;
  title: string;
  time: string;
  type: string;
  iconHref: string;
};

type Post = {
  id: string;
  avatar: string;
  author: string;
  role: string;
  time: string;
  content: string;
  likes: number;
  comments: number;
  shares: number;
};

type CourseProgress = {
  id: string;
  title: string;
  mark: string;
  percent: number;
  markClass?: string; // "green"
  barClass?: string;  // "data-science" | "web-dev"
};

type IconProps = {
  href: string;
  w?: number;
  h?: number;
  className?: string;
};

function Icon({ href, w = 16, h = 16, className = "" }: IconProps) {
  return (
    <svg width={w} height={h} className={className} aria-hidden="true">
      <use href={href} />
    </svg>
  );
}

/* -------- Left -------- */

function ProfileCard({ profile }: { profile: ProfileData }) {
  return (
    <article className="profile-card">
      <a href="#" className="profile-card-avatar-link">
        <img src={profile.avatar} alt={profile.name} className="avatar" />
      </a>

      <h3 className="name">{profile.name}</h3>
      <p className="role">{profile.role}</p>

      <div className="stats">
        <div>
          <Icon href="/img/icons.svg#icon-SVG-1" w={20} h={20} />
          <h4 className="courses">{profile.courses}</h4>
          <p className="profile-template-paragraf">Courses</p>
        </div>

        <div>
          <Icon href="/img/icons.svg#icon-SVG-3" w={20} h={20} />
          <h4 className="assignments">{profile.assignments}</h4>
          <p className="profile-template-paragraf">Assignments Due</p>
        </div>

        <div>
          <Icon href="/img/icons.svg#icon-SVG-9" w={20} h={20} />
          <h4 className="groups">{profile.groups}</h4>
          <p className="profile-template-paragraf">Study Groups</p>
        </div>
      </div>
    </article>
  );
}

function QuickActions() {
  return (
    <section className="quick-actions">
      <h2 className="title">Quick Actions</h2>
      <ul>
        <li>
          <button type="button">
            <Icon href="/img/icons.svg#icon-SVG-4" w={16} h={16} />
            Show rating
          </button>
        </li>
        <li>
          <button type="button">
            <Icon href="/img/icons.svg#icon-SVG-3" w={16} h={16} />
            Check Schedule
          </button>
        </li>
        <li>
          <button type="button">
            <Icon href="/img/icons.svg#icon-bell-1" w={16} h={16} />
            Notifications
          </button>
        </li>
      </ul>
    </section>
  );
}

/* -------- Center -------- */

function PostCard({
  post,
  liked,
  onToggleLike,
}: {
  post: Post;
  liked: boolean;
  onToggleLike: () => void;
}) {
  const likesShown = post.likes + (liked ? 1 : 0);

  return (
    <article className="post">
      <div className="post-header">
        <a href="#">
          <img src={post.avatar} alt="author profile pic" width={48} height={48} className="avatar" />
        </a>

        <div className="post-header-info">
          <h3 className="author">{post.author}</h3>
          <p className="role">{post.role}</p>
          <h4 className="time">{post.time}</h4>
        </div>
      </div>

      <p className="content">{post.content}</p>

      <div className="stats">
        <ul>
          <li className="like-item">
            <label className="like-checkbox">
              <input
                type="checkbox"
                className="like-input"
                checked={liked}
                onChange={onToggleLike}
              />
              <svg width="16" height="16" className="like-icon" aria-hidden="true">
                <use href="/img/icons.svg#icon-SVG-6" />
              </svg>
              <svg width="16" height="16" className="like-icon-active" aria-hidden="true">
                <use href="/img/icons.svg#icon-like-active" />
              </svg>
            </label>
            <span className="likes">{likesShown}</span>
          </li>

          <li>
            <Icon href="/img/icons.svg#icon-SVG-7" w={16} h={16} />
            <span className="comments">{post.comments}</span>
          </li>

          <li>
            <Icon href="/img/icons.svg#icon-SVG-8" w={16} h={16} />
            <span className="shares">{post.shares}</span>
          </li>
        </ul>
      </div>
    </article>
  );
}

function Feed({ posts }: { posts: Post[] }) {
  const [likedMap, setLikedMap] = useState<Record<string, boolean>>({});

  return (
    <section className="feed">
      <h2 className="title hidden-title">Feed</h2>

      <div id="feed-list">
        {posts.map((p) => (
          <div key={p.id} style={{ marginBottom: 16 }}>
            <PostCard
              post={p}
              liked={!!likedMap[p.id]}
              onToggleLike={() =>
                setLikedMap((prev) => ({ ...prev, [p.id]: !prev[p.id] }))
              }
            />
          </div>
        ))}
      </div>
    </section>
  );
}

/* -------- Right -------- */

function Upcoming({ events }: { events: UpcomingEvent[] }) {
  return (
    <section className="upcoming">
      <div className="title-container">
        <Icon href="/img/icons.svg#icon-SVG-3" w={16} h={16} />
        <h2 className="title">Upcoming</h2>
      </div>

      <ul id="upcoming-list">
        {events.map((e) => {
          const isAssignment = e.type === "Assignment";
          const isLab = e.type === "Lab";
          const isLecture = e.type === "Lecture";

          const showTypePill = !isAssignment;              // Assignment — не показуємо бейдж
          const showTimeIcon = !(isLab || isLecture);      // Lab/Lecture — не показуємо svg в часі

          return (
            <li className="upcoming-item" key={e.id}>
              <h3 className="title">{e.title}</h3>

              <div className={`time-container ${isAssignment ? "time-container--assignment" : ""}`}>
                {showTimeIcon && <Icon href={e.iconHref} w={16} h={16} className="icon" />}
                <p className="time">{e.time}</p>
              </div>

              {showTypePill && (
                <h4 className={`type ${isLecture ? "type--lecture" : ""}`}>
                  {e.type}
                </h4>
              )}
            </li>
          );
        })}
      </ul>
    </section>
  );
}


function CourseProgress({ items }: { items: CourseProgress[] }) {
  return (
    <section className="course-progress">
      <div className="title-container">
        <Icon href="/img/icons.svg#icon-SVG-10" w={16} h={16} />
        <h2 className="title">Course Progress</h2>
      </div>

      <ul>
        {items.map((c) => (
          <li key={c.id}>
            <div className="course-header">
              <h3>{c.title}</h3>
              <p className={c.markClass ?? ""}>{c.mark}</p>
            </div>

            <div className="progress-bar" aria-label={`${c.title} progress`}>
              <div
                className={`progress-bar-completed ${c.barClass ?? ""}`}
                style={{ width: `${c.percent}%` }}   // важливо: щоб % були з даних
              />
            </div>

            <p className="progress-completed-description">{c.percent}% complete</p>
          </li>
        ))}
      </ul>
    </section>
  );
}

/* -------- Page -------- */

export default function Dashboard() {
  const profile: ProfileData = useMemo(
    () => ({
      name: "Emily Johnson",
      role: "Computer Science Student",
      avatar: "/img/profile_pic.png",
      courses: 3,
      assignments: 2,
      groups: 5,
    }),
    []
  );

  const upcomingEvents: UpcomingEvent[] = useMemo(
    () => [
      {
        id: "u1",
        title: "Programming Assignment Due",
        time: "Tomorrow 11:59 PM",
        type: "Assignment",
        iconHref: "/img/icons.svg#icon-SVG-11",
      },
      {
        id: "u2",
        title: "Data Science Lab",
        time: "Friday 2:00 PM",
        type: "Lab",
        iconHref: "/img/icons.svg#icon-SVG-3",
      },
      {
        id: "u3",
        title: "Study Group Meeting",
        time: "Saturday 10:00 AM",
        type: "Lecture",
        iconHref: "/img/icons.svg#icon-SVG-4",
      },
    ],
    []
  );

  const posts: Post[] = useMemo(
    () => [
      {
        id: "p1",
        avatar: "/img/profile_pic.png",
        author: "Dr. Sarah Smith",
        role: "Computer Science Professor",
        time: "2h ago",
        content:
          "Excited to announce our new Advanced Machine Learning course! Registration opens next week. This course will cover deep learning, neural networks, and practical AI applications. 🤖",
        likes: 24,
        comments: 8,
        shares: 3,
      },
      {
        id: "p2",
        avatar: "/img/profile_pic.png",
        author: "Emily Johnson",
        role: "Student • Computer Science",
        time: "4h ago",
        content:
          "Just finished my final project for Web Development class! Built a full-stack e-commerce app with React and Node.js. Thanks to all my classmates who helped along the way! 💻",
        likes: 24,
        comments: 8,
        shares: 3,
      },
      {
        id: "p3",
        avatar: "/img/profile_pic.png",
        author: "Duckucate Official",
        role: "Educational Platform",
        time: "6h ago",
        content:
          "📚 Study Tip Tuesday: Use the Pomodoro Technique! Study for 25 minutes, then take a 5-minute break. This helps maintain focus and prevents burnout. What's your favorite study method?",
        likes: 102,
        comments: 8,
        shares: 3,
      },
    ],
    []
  );

  const courseProgress: CourseProgress[] = useMemo(
    () => [
      { id: "c1", title: "Programming", mark: "4+", percent: 75 },
      { id: "c2", title: "Data Science", mark: "5-", percent: 50, markClass: "green", barClass: "data-science" },
      { id: "c3", title: "Web Dev", mark: "4+", percent: 80, barClass: "web-dev" },
    ],
    []
  );

  return (
    <main className="dashboard sidebar-space">
      <div className="left">
        <section className="profile">
          <h2 className="title hidden-title">Profile</h2>
          <div id="profile-list">
            <ProfileCard profile={profile} />
          </div>
        </section>

        <QuickActions />
      </div>

      <div className="center">
        <Feed posts={posts} />
      </div>

      <div className="right">
        <Upcoming events={upcomingEvents} />
        <CourseProgress items={courseProgress} />
      </div>
    </main>
  );
}
