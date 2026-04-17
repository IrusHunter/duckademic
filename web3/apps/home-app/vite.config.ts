import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
    name: 'homeApp',
    filename: 'remoteEntry.js',
    exposes: {
      './HomeApp': './src/App.tsx'
    }
  })
  ],
  server: {
    port: 5005
  }
})