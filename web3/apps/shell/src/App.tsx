import { Suspense, lazy } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'

const AuthApp = lazy(() => import('authApp/AuthApp'))
const ClassroomApp = lazy(() => import('classroomApp/ClassroomApp'))
const HomeApp = lazy(() => import('homeApp/HomeApp'))

function App() {
  return (
    <BrowserRouter>
      <Suspense fallback={<div>Loading...</div>}>
        <Routes>
          <Route path="/" element={<HomeApp />} />
          <Route path="/login" element={<AuthApp />} />
          <Route path="/classroom/*" element={<ClassroomApp />} />
          {/* <Route path="/" element={<div>test</div>} /> */}
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}

export default App