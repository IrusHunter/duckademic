import Icon from '../Icon/Icon'
import Avatar from '../Avatar/Avatar'
import css from '../App/App.module.css'

interface ProfileCardProps {
  name: string
  role: string
  courses: number | string
  assignments: number | string
  groups: number | string
  hideAssignments?: boolean
}

export default function ProfileCard({
  name, role, courses, assignments, groups, hideAssignments = false,
}: ProfileCardProps) {
  return (
    <article className={css.profileCard}>
      <span className={css.profileAvatarLink}>
        <Avatar name={name} size={80} />
      </span>

      <h3 className={css.name}>{name}</h3>
      <p className={css.role}>{role}</p>

      <div className={`${css.stats}${hideAssignments ? ` ${css.statsPair}` : ''}`}>
        <div>
          <Icon id="icon-SVG-1" size={20} />
          <h4>{courses}</h4>
          <p className={css.paragraf}>Courses</p>
        </div>
        {!hideAssignments && (
          <div>
            <Icon id="icon-SVG-3" size={20} />
            <h4>{assignments}</h4>
            <p className={css.paragraf}>Assignments Due</p>
          </div>
        )}
        <div>
          <Icon id="icon-SVG-9" size={20} />
          <h4>{groups}</h4>
          <p className={css.paragraf}>Study Groups</p>
        </div>
      </div>
    </article>
  )
}
