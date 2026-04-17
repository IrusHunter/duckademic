import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
      name: 'authApp',
      filename: 'remoteEntry.js',
      exposes: {
        './AuthApp': './src/App.tsx'
      },
      shared: ['react', 'react-dom']
    })
  ],
  server: {
    port: 5001
  }
})