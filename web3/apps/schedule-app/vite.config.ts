import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
      name: 'scheduleApp',
      filename: 'remoteEntry.js',
      exposes: {
        './ScheduleApp': './src/App.tsx'
      },
      shared: ['react', 'react-dom']
    })
  ],
  build: { target: 'esnext' },
  server: { port: 5008 },
  preview: { port: 5008 }
})