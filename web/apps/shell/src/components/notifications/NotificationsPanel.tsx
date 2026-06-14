import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getNotifications, markAllNotificationsRead } from '../../services/notificationService'
import type { Notification } from '../../services/notificationService'
import css from './NotificationsPanel.module.css'

function typeLabel(type: Notification['type']) {
  if (type === 'new_task') return 'New task'
  if (type === 'submission') return 'Submission'
  return 'Grade'
}

function typeBadgeClass(type: Notification['type'], styles: typeof css) {
  if (type === 'new_task') return styles.badgeNewTask
  if (type === 'submission') return styles.badgeSubmission
  return styles.badgeGrade
}

function timeAgo(iso: string) {
  const diff = Date.now() - new Date(iso).getTime()
  const m = Math.floor(diff / 60_000)
  if (m < 1) return 'just now'
  if (m < 60) return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h ago`
  return `${Math.floor(h / 24)}d ago`
}

function notifLink(n: Notification) {
  return `/classroom/${n.course_id}/task/${n.task_id}`
}

interface Props {
  onClose: () => void
}

export default function NotificationsPanel({ onClose }: Props) {
  const qc = useQueryClient()

  const { data: notifications = [] } = useQuery({
    queryKey: ['notifications'],
    queryFn: getNotifications,
    refetchInterval: 30_000,
  })

  const markRead = useMutation({
    mutationFn: markAllNotificationsRead,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['notifications'] }),
  })

  const hasUnread = notifications.some((n) => !n.is_read)

  return (
    <div className={css.panel}>
      <div className={css.header}>
        <span className={css.headerTitle}>Notifications</span>
        {hasUnread && (
          <button
            className={css.markReadBtn}
            onClick={() => markRead.mutate()}
            disabled={markRead.isPending}
          >
            Mark all as read
          </button>
        )}
      </div>

      <div className={css.list}>
        {notifications.length === 0 ? (
          <p className={css.empty}>No notifications yet</p>
        ) : (
          notifications.map((n) => (
            <Link
              key={n.id}
              to={notifLink(n)}
              className={`${css.item} ${!n.is_read ? css.itemUnread : ''}`}
              onClick={onClose}
            >
              <span className={`${css.dot} ${n.is_read ? css.dotRead : ''}`} />
              <div className={css.itemBody}>
                <div className={css.itemTitle}>{n.title}</div>
                <div className={css.itemMeta}>
                  <span className={`${css.typeBadge} ${typeBadgeClass(n.type, css)}`}>
                    {typeLabel(n.type)}
                  </span>
                  {' · '}
                  {timeAgo(n.created_at)}
                </div>
              </div>
            </Link>
          ))
        )}
      </div>
    </div>
  )
}
