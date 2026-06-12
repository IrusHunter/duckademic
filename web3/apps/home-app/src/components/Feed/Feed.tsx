import { useState } from 'react'
import Icon from '../Icon/Icon'
import Avatar from '../Avatar/Avatar'
import type { Post } from '../../types/dashboard'
import css from '../App/App.module.css'

function PostCard({ post, liked, onToggle }: { post: Post; liked: boolean; onToggle: () => void }) {
  const likesShown = post.likes + (liked ? 1 : 0)
  return (
    <article className={css.post}>
      <div className={css.postHeader}>
        <Avatar name={post.author} size={48} />
        <div className={css.postHeaderInfo}>
          <h3 className={css.author}>{post.author}</h3>
          <p className={css.role}>{post.role}</p>
          <h4 className={css.time}>{post.time}</h4>
        </div>
      </div>

      <p className={css.content}>{post.content}</p>

      <div className={css.postStats}>
        <ul>
          <li className={css.likeItem}>
            <label className={css.likeCheckbox}>
              <input
                type="checkbox"
                className={css.likeInput}
                checked={liked}
                onChange={onToggle}
              />
              <svg width="16" height="16" className={css.likeIcon} aria-hidden="true">
                <use href="/img/icons.svg#icon-SVG-6" />
              </svg>
              <svg width="16" height="16" className={css.likeIconActive} aria-hidden="true">
                <use href="/img/icons.svg#icon-like-active" />
              </svg>
            </label>
            <span className={css.likes}>{likesShown}</span>
          </li>
          <li>
            <Icon id="icon-SVG-7" size={16} />
            <span className={css.comments}>{post.comments}</span>
          </li>
          <li>
            <Icon id="icon-SVG-8" size={16} />
            <span className={css.shares}>{post.shares}</span>
          </li>
        </ul>
      </div>
    </article>
  )
}

// Стрічка поки на моках — бекенд цей функціонал ще не реалізував.
const MOCK_POSTS: Post[] = [
  {
    id: 'p1', avatar: '', author: 'Dr. Sarah Smith',
    role: 'Computer Science Professor', time: '2h ago',
    content: 'Excited to announce our new Advanced Machine Learning course! Registration opens next week. This course will cover deep learning, neural networks, and practical AI applications. 🤖',
    likes: 24, comments: 8, shares: 3,
  },
  {
    id: 'p2', avatar: '', author: 'Emily Johnson',
    role: 'Student • Computer Science', time: '4h ago',
    content: 'Just finished my final project for Web Development class! Built a full-stack e-commerce app with React and Node.js. Thanks to all my classmates who helped along the way! 💻',
    likes: 24, comments: 8, shares: 3,
  },
  {
    id: 'p3', avatar: '', author: 'Duckademic Official',
    role: 'Educational Platform', time: '6h ago',
    content: '📚 Study Tip Tuesday: Use the Pomodoro Technique! Study for 25 minutes, then take a 5-minute break. This helps maintain focus and prevents burnout. What`s your favorite study method?',
    likes: 102, comments: 8, shares: 3,
  },
]

export default function Feed() {
  const [liked, setLiked] = useState<Record<string, boolean>>({})
  return (
    <section className={css.feed}>
      <h2 className={`${css.title} ${css.hiddenTitle}`}>Feed</h2>
      <div>
        {MOCK_POSTS.map((p) => (
          <PostCard
            key={p.id}
            post={p}
            liked={!!liked[p.id]}
            onToggle={() => setLiked((prev) => ({ ...prev, [p.id]: !prev[p.id] }))}
          />
        ))}
      </div>
    </section>
  )
}
