import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '../../../packages/shared-styles.css'
import App from './components/App/App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)