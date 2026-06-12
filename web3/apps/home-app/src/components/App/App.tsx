import { useMemo } from 'react'
import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from '@tanstack/react-query'
import { getPersonalSchedule } from '../../services/scheduleService'
import { getUpcomingEvents } from '../../services/courseService'
import { addDays, startOfToday } from '../../utils/datetime'
import ScheduleList from '../ScheduleList/ScheduleList'
import UpcomingEvents from '../UpcomingEvents/UpcomingEvents'
import Loader from '../Loader/Loader'
import ErrorMessage from '../ErrorMessage/ErrorMessage'
import css from './App.module.css'

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 60_000, retry: 1 } },
})

function Home() {
  // Тиждень: від початку сьогоднішнього дня до +7 днів.
  const range = useMemo(() => {
    const start = startOfToday()
    return {
      startISO: start.toISOString(),
      endISO: addDays(start, 7).toISOString(),
    }
  }, [])

  const schedule = useQuery({
    queryKey: ['personal-schedule', range.startISO, range.endISO],
    queryFn: () =>
      getPersonalSchedule({ startTime: range.startISO, endTime: range.endISO }),
  })

  const events = useQuery({
    queryKey: ['upcoming-events', range.startISO],
    queryFn: () => getUpcomingEvents({ count: 5, startTime: range.startISO }),
  })

  return (
    <div className={css.page}>
      <div className={css.inner}>
        <header className={css.header}>
          <h1 className={css.heading}>Головна</h1>
          <p className={css.subheading}>Розклад на тиждень і найближчі дедлайни</p>
        </header>

        <div className={css.grid}>
          {/* Основна колонка — розклад */}
          <section className={css.main}>
            <h2 className={css.sectionTitle}>Мій розклад</h2>
            {schedule.isLoading && <Loader label="Завантаження розкладу…" />}
            {schedule.isError && (
              <ErrorMessage
                message="Не вдалося завантажити розклад."
                onRetry={() => schedule.refetch()}
              />
            )}
            {schedule.data && <ScheduleList lessons={schedule.data} />}
          </section>

          {/* Бічна колонка — дедлайни */}
          <aside className={css.aside}>
            <h2 className={css.sectionTitle}>Найближчі події</h2>
            {events.isLoading && <Loader label="Завантаження подій…" />}
            {events.isError && (
              <ErrorMessage
                message="Не вдалося завантажити події."
                onRetry={() => events.refetch()}
              />
            )}
            {events.data && <UpcomingEvents tasks={events.data} />}
          </aside>
        </div>
      </div>
    </div>
  )
}

// QueryClientProvider всередині App — щоб працювало і всередині shell,
// і при ізольованому запуску модуля (як в admin-app).
export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Home />
    </QueryClientProvider>
  )
}
