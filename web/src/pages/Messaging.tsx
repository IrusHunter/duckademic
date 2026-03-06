import { useMemo, useRef, useState, type FormEvent } from "react";

function cx(...parts: Array<string | false | undefined | null>) {
  return parts.filter(Boolean).join(" ");
}

type ChatItem = {
  id: string;
  name: string;
  unread: number;
  lastMessage: string;
  participantsText: string;
  timeText: string;
};

type ChatMessage = {
  id: string;
  chatId: string;
  author: string;
  avatar?: string;
  text: string;
  time: string;
  own?: boolean;
};

const chatsSeed: ChatItem[] = [
  {
    id: "general",
    name: "General Discussion",
    unread: 2,
    lastMessage: "Hi everyone. Let's discuss the group assignment.",
    participantsText: "45 participants",
    timeText: "2 min ago",
  },
  {
    id: "help",
    name: "Programming Help",
    unread: 1,
    lastMessage: "Can someone help with the loop?",
    participantsText: "32 participants",
    timeText: "15 min ago",
  },
  {
    id: "study",
    name: "Study Group",
    unread: 0,
    lastMessage: "Meeting tomorrow at 3 PM",
    participantsText: "8 participants",
    timeText: "1 hour ago",
  },
];

const messagesSeed: ChatMessage[] = [
  {
    id: "m1",
    chatId: "general",
    author: "Emily Johnson",
    avatar: "/img/profile_pic.png",
    text: "Hi everyone. Let's discuss the group assignment.",
    time: "9:20 AM",
  },
  {
    id: "m2",
    chatId: "general",
    author: "Sarah Wilson",
    avatar: "/img/profile_pic.png",
    text: "I think we should focus on the data analysis part first.",
    time: "9:22 AM",
  },
  {
    id: "m3",
    chatId: "general",
    author: "You",
    text: "Agreed! I can handle the visualization components.",
    time: "9:22 AM",
    own: true,
  },
];

export default function Messaging() {
  const [query, setQuery] = useState("");
  const [activeChatId, setActiveChatId] = useState<string>(chatsSeed[0]?.id ?? "general");
  const [chats, setChats] = useState<ChatItem[]>(chatsSeed);
  const [messages, setMessages] = useState<ChatMessage[]>(messagesSeed);
  const [draft, setDraft] = useState("");

  const listRef = useRef<HTMLDivElement | null>(null);

  const filteredChats = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return chats;

    return chats.filter((c) => {
      return (
        c.name.toLowerCase().includes(q) ||
        c.lastMessage.toLowerCase().includes(q) ||
        c.participantsText.toLowerCase().includes(q)
      );
    });
  }, [chats, query]);

  const activeChat = useMemo(() => chats.find((c) => c.id === activeChatId) ?? chats[0], [chats, activeChatId]);

  const activeMessages = useMemo(
    () => messages.filter((m) => m.chatId === activeChatId),
    [messages, activeChatId]
  );

  function openChat(id: string) {
    setActiveChatId(id);

    // скидаємо unread для відкритого чату
    setChats((prev) => prev.map((c) => (c.id === id ? { ...c, unread: 0 } : c)));
  }

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    const text = draft.trim();
    if (!text) return;

    const now = new Date();
    const time = now.toLocaleTimeString("en-US", {
      hour: "numeric",
      minute: "2-digit",
      hour12: true,
    });

    const newMsg: ChatMessage = {
      id: `m-${crypto?.randomUUID?.() ?? Math.random().toString(16).slice(2)}`,
      chatId: activeChatId,
      author: "You",
      text,
      time,
      own: true,
    };

    setMessages((prev) => [...prev, newMsg]);
    setDraft("");

    // оновлюємо lastMessage + time у списку чатів
    setChats((prev) =>
      prev.map((c) =>
        c.id === activeChatId
          ? { ...c, lastMessage: text, timeText: "just now" }
          : c
      )
    );

    // прокрутка вниз
    requestAnimationFrame(() => {
      const el = listRef.current;
      if (!el) return;
      el.scrollTop = el.scrollHeight;
    });
  }

  return (
    <main className="sidebar-space messaging-info">
      <div className="messenger">
        {/* ===== LEFT: chat groups ===== */}
        <aside className="chat-groups">
          <header className="chat-groups-header">
            <h1 className="chat-groups-title">DuckChat</h1>

            <div className="chat-groups-search">
              <svg width="16" height="16" className="icon" aria-hidden="true">
                <use href="/img/icons.svg#icon-search-1" />
              </svg>

              <input
                type="text"
                placeholder="Search conversations..."
                className="chat-groups-search-input"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
              />
            </div>
          </header>

          <ul className="chat-groups-list">
            {filteredChats.map((c) => (
              <li
                key={c.id}
                className={cx("chat-groups-item", c.id === activeChatId && "chat-groups-item--active")}
                onClick={() => openChat(c.id)}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => {
                  if (e.key === "Enter" || e.key === " ") openChat(c.id);
                }}
                aria-current={c.id === activeChatId}
              >
                <div className="chat-groups-item-main">
                  <div className="chat-groups-item-top">
                    <h2 className="chat-groups-name">{c.name}</h2>

                    <span className={cx("chat-groups-unread", c.unread === 0 && "no-new-messages")}>
                      {c.unread}
                    </span>
                  </div>

                  <p className="chat-groups-last-message">{c.lastMessage}</p>

                  <div className="chat-groups-meta">
                    <span className="chat-groups-participants">{c.participantsText}</span>
                    <span className="chat-groups-time">{c.timeText}</span>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </aside>

        {/* ===== RIGHT: active chat ===== */}
        <section className="chat-messager">
          <header className="messager-header">
            <div className="messager-header-main">
              <h2 className="messager-title">{activeChat?.name ?? ""}</h2>
              <p className="messager-participants">{activeChat?.participantsText ?? ""}</p>
            </div>

            <div className="messager-header-actions">
              <button type="button" className="messager-header-btn" aria-label="Chat options">
                <svg width="16" height="16" className="icon" aria-hidden="true">
                  <use href="/img/icons.svg#icon-SVG-12" />
                </svg>
              </button>
            </div>
          </header>

          <div className="chat">
            <div className="chat-messages" ref={listRef}>
              {activeMessages.map((m) => (
                <div key={m.id} className={cx("message", m.own && "message--own")}>
                  {!m.own && (
                    <div className="message-avatar">
                      <img src={m.avatar ?? "/img/profile_pic.png"} alt={`${m.author} avatar`} />
                    </div>
                  )}

                  <div className="message-body">
                    <span className="message-author">{m.author}</span>
                    <div className="message-bubble">{m.text}</div>
                    <span className="message-time">{m.time}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <form className="messager-footer chat-input" onSubmit={onSubmit}>
            <input
              type="text"
              placeholder="Type your message..."
              className="chat-input-field"
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
            />

            <button type="button" className="chat-input-btn" aria-label="Emoji">
              <svg width="28" height="28" className="icon icon-emoji" aria-hidden="true">
                <use href="/img/icons.svg#icon-happy-1" />
              </svg>
            </button>

            <button type="submit" className="chat-input-btn chat-input-btn--send" aria-label="Send">
              <svg width="16" height="16" className="icon icon-send" aria-hidden="true">
                <use href="/img/icons.svg#icon-SVG-13" />
              </svg>
            </button>
          </form>
        </section>
      </div>
    </main>
  );
}
