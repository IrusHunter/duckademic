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
      shared: {
        react: { singleton: true, requiredVersion: '^19.0.0' } as any,
        'react-dom': { singleton: true, requiredVersion: '^19.0.0' } as any,
        axios: { singleton: true } as any,
      }
    })
  ],
  build: { target: 'esnext' },
  server: {
  port: 5001,
  proxy: {
    '/api/auth': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    }
  }
},
  preview: { port: 5001 }
})