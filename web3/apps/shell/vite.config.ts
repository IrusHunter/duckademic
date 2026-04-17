import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
      name: 'shell',
      remotes: {
        authApp: 'http://localhost:5001/assets/remoteEntry.js',
        classroomApp: 'http://localhost:5002/assets/remoteEntry.js',
        dashboardApp: 'http://localhost:5003/assets/remoteEntry.js',
        gradesApp: 'http://localhost:5004/assets/remoteEntry.js',
        homeApp: 'http://localhost:5005/assets/remoteEntry.js',
        messengerApp: 'http://localhost:5006/assets/remoteEntry.js',
        scheduleApp: 'http://localhost:5007/assets/remoteEntry.js'
      },
      shared: ['react', 'react-dom']
    })
  ],
  server: {
    port: 5000
  }
})