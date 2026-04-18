import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
      name: 'messengerApp',
      filename: 'remoteEntry.js',
      exposes: {
        './MessengerApp': './src/App.tsx'
      },
      shared: ['react', 'react-dom']
    })
  ],
  build: { target: 'esnext' },
  server: { port: 5007 },
  preview: { port: 5007 }
})