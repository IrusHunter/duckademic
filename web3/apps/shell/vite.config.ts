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
        homeApp: 'http://localhost:5006/assets/remoteEntry.js',
      },
      shared: {
        react: {
          singleton: true,
          requiredVersion: '^19.0.0'
        } as any,
        'react-dom': {
          singleton: true,
          requiredVersion: '^19.0.0'
        } as any
      }
    })
  ],
  build: { target: 'esnext' },
  server: { port: 5000 }
})