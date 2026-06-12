import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './components/App/App.tsx'

// main.tsx — точка входу лише для ізольованого запуску модуля (npm run dev/serve).
// Усередині shell використовується експозиція './HomeApp' → components/App/App.tsx.
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
