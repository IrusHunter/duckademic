import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
      name: 'gradesApp',
      filename: 'remoteEntry.js',
      exposes: {
        './GradesApp': './src/App.tsx'
      },
      shared: ['react', 'react-dom']
    })
  ],
  server: {
    port: 5004
  }
})