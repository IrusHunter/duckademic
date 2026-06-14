import { useMemo, useRef, useState, type FormEvent } from 'react'
import Avatar from '../Avatar/Avatar'
import css from './App.module.css'

function cx(...parts: Array<string | false | undefined | null>) {
  return parts.filter(Boolean).join(' ')
}

type ChatItem = {
  id: string
  name: string
  unread: number
  lastMessage: string
  participantsText: string
  timeText: string
}

type ChatMessage = {
  id: string
  chatId: string
  author: string
  text: string
  time: string
  own?: boolean
}

const chatsSeed: ChatItem[] = [
  {
    id: 'general',
    name: 'General Discussion',
    unread: 2,
    lastMessage: 'Agreed! I can handle the visualization components.',
    participantsText: '45 participants',
    timeText: '2 min ago',
  },
  {
    id: 'help',
    name: 'Programming Help',
    unread: 1,
    lastMessage: 'Can someone help with the loop?',
    participantsText: '32 participants',
    timeText: '15 min ago',
  },
  {
    id: 'study',
    name: 'Study Group',
    unread: 0,
    lastMessage: 'Meeting tomorrow at 3 PM',
    participantsText: '8 participants',
    timeText: '1 hour ago',
  },
]

const messagesSeed: ChatMessage[] = [
  { id: 'm1', chatId: 'general', author: 'Emily Johnson', text: "Hi everyone. Let's discuss the group assignment.", time: '9:20 AM' },
  { id: 'm2', chatId: 'general', author: 'Sarah Wilson', text: 'I think we should focus on the data analysis part first.', time: '9:22 AM' },
  { id: 'm3', chatId: 'general', author: 'You', text: 'Agreed! I can handle the visualization components.', time: '9:23 AM', own: true },
  { id: 'm4', chatId: 'help', author: 'Alex Kim', text: 'Can someone help with the loop?', time: '9:10 AM' },
  { id: 'm5', chatId: 'study', author: 'David Chen', text: 'Meeting tomorrow at 3 PM', time: '8:45 AM' },
]

export default function App() {
  const [query, setQuery] = useState('')
  const [activeChatId, setActiveChatId] = useState<string>(chatsSeed[0]?.id ?? 'general')
  const [chats, setChats] = useState<ChatItem[]>(chatsSeed)
  const [messages, setMessages] = useState<ChatMessage[]>(messagesSeed)
  const [draft, setDraft] = useState('')

  const listRef = useRef<HTMLDivElement | null>(null)

  const filteredChats = useMemo(() => {
    const q = query.trim().toLowerCase()
    if (!q) return chats
    return chats.filter(
      (c) =>
        c.name.toLowerCase().includes(q) ||
        c.lastMessage.toLowerCase().includes(q) ||
        c.participantsText.toLowerCase().includes(q),
    )
  }, [chats, query])

  const activeChat = useMemo(
    () => chats.find((c) => c.id === activeChatId) ?? chats[0],
    [chats, activeChatId],
  )

  const activeMessages = useMemo(
    () => messages.filter((m) => m.chatId === activeChatId),
    [messages, activeChatId],
  )

  function openChat(id: string) {
    setActiveChatId(id)
    setChats((prev) => prev.map((c) => (c.id === id ? { ...c, unread: 0 } : c)))
  }

  function onSubmit(e: FormEvent) {
    e.preventDefault()
    const text = draft.trim()
    if (!text) return

    const time = new Date().toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true,
    })

    const newMsg: ChatMessage = {
      id: `m-${crypto?.randomUUID?.() ?? Math.random().toString(16).slice(2)}`,
      chatId: activeChatId,
      author: 'You',
      text,
      time,
      own: true,
    }

    setMessages((prev) => [...prev, newMsg])
    setDraft('')
    setChats((prev) =>
      prev.map((c) =>
        c.id === activeChatId ? { ...c, lastMessage: text, timeText: 'just now' } : c,
      ),
    )

    requestAnimationFrame(() => {
      const el = listRef.current
      if (!el) return
      el.scrollTop = el.scrollHeight
    })
  }

  return (
    <main className={css.page}>
      <div className={css.messenger}>
        {/* ===== LEFT: chat list ===== */}
        <aside className={css.chatGroups}>
          <header className={css.chatGroupsHeader}>
            <h1 className={css.chatGroupsTitle}>DuckChat</h1>

            <div className={css.chatGroupsSearch}>
              <svg width="16" height="16" className={css.searchIcon} aria-hidden="true">
                <use href="/img/icons.svg#icon-search-1" />
              </svg>
              <input
                type="text"
                placeholder="Search conversations..."
                className={css.chatGroupsSearchInput}
                value={query}
                onChange={(e) => setQuery(e.target.value)}
              />
            </div>
          </header>

          <ul className={css.chatGroupsList}>
            {filteredChats.map((c) => (
              <li
                key={c.id}
                className={cx(css.chatGroupsItem, c.id === activeChatId && css.chatGroupsItemActive)}
                onClick={() => openChat(c.id)}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') openChat(c.id)
                }}
                aria-current={c.id === activeChatId}
              >
                <div className={css.chatGroupsItemMain}>
                  <div className={css.chatGroupsItemTop}>
                    <h2 className={css.chatGroupsName}>{c.name}</h2>
                    <span className={cx(css.chatGroupsUnread, c.unread === 0 && css.noNewMessages)}>
                      {c.unread}
                    </span>
                  </div>
                  <p className={css.chatGroupsLastMessage}>{c.lastMessage}</p>
                  <div className={css.chatGroupsMeta}>
                    <span className={css.chatGroupsParticipants}>{c.participantsText}</span>
                    <span className={css.chatGroupsTime}>{c.timeText}</span>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </aside>

        {/* ===== RIGHT: active chat ===== */}
        <section className={css.chatMessager}>
          <header className={css.messagerHeader}>
            <div className={css.messagerHeaderMain}>
              <h2 className={css.messagerTitle}>{activeChat?.name ?? ''}</h2>
              <p className={css.messagerParticipants}>{activeChat?.participantsText ?? ''}</p>
            </div>
            <div className={css.messagerHeaderActions}>
              <button type="button" className={css.messagerHeaderBtn} aria-label="Chat options">
                <svg width="16" height="16" aria-hidden="true">
                  <use href="/img/icons.svg#icon-SVG-12" />
                </svg>
              </button>
            </div>
          </header>

          <div className={css.chat}>
            <div className={css.chatMessages} ref={listRef}>
              {activeMessages.map((m) => (
                <div key={m.id} className={cx(css.message, m.own && css.messageOwn)}>
                  {!m.own && <div className={css.messageAvatarWrap}><Avatar name={m.author} size={32} /></div>}
                  <div className={css.messageBody}>
                    <span className={css.messageAuthor}>{m.author}</span>
                    <div className={css.messageBubble}>{m.text}</div>
                    <span className={css.messageTime}>{m.time}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <form className={css.messagerFooter} onSubmit={onSubmit}>
            <input
              type="text"
              placeholder="Type your message..."
              className={css.chatInputField}
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
            />
            <button type="button" className={css.chatInputBtn} aria-label="Emoji">
              <svg width="28" height="28" aria-hidden="true">
                <use href="/img/icons.svg#icon-happy-1" />
              </svg>
            </button>
            <button type="submit" className={`${css.chatInputBtn} ${css.chatInputBtnSend}`} aria-label="Send">
              <svg width="16" height="16" aria-hidden="true">
                <use href="/img/icons.svg#icon-SVG-13" />
              </svg>
            </button>
          </form>
        </section>
      </div>
    </main>
  )
}
