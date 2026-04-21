import { Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { SERVICES } from '../../config/services'
import { Sidebar } from '../Sidebar/Sidebar'
import { ServiceHome } from '../ServiceHome/ServiceHome'
import { DynamicServiceRoute } from '../DynamicServiceRoute/DynamicServiceRoute'
import css from './App.module.css'

const queryClient = new QueryClient()

function AdminLayout() {
  return (
    <div style={{ display: 'flex'}}>
      <Sidebar />
      <div style={{ flex: 1 }} className={css.mainDiv}>
        <Routes>
          <Route path="/" element={<ServiceHome service={SERVICES[0]} />} />
          <Route path=":serviceKey" element={<DynamicServiceRoute />} />
          <Route path=":serviceKey/:tableKey" element={<DynamicServiceRoute />} />
          <Route path=":serviceKey/:tableKey/edit/:itemId" element={<DynamicServiceRoute />} />
        </Routes>
      </div>
    </div>
  )
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AdminLayout />
    </QueryClientProvider>
  )
}