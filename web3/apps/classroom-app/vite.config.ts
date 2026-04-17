import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
    name: 'classroomApp',
    filename: 'remoteEntry.js',
    exposes: {
      './ClassroomApp': './src/App.tsx'
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
  server: {
    port: 5002
  },
  preview: { port: 5002 }
})