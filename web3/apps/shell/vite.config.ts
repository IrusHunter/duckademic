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
        adminApp: 'http://localhost:5010/assets/remoteEntry.js',
      },
      shared: {
        react: { singleton: true, requiredVersion: '^19.0.0' } as any,
        'react-dom': { singleton: true, requiredVersion: '^19.0.0' } as any,
        'react-router-dom': { singleton: true, requiredVersion: '^7.0.0' } as any,
      }
    })
  ],
  build: { target: 'esnext' },
  server: {
  port: 5000,
  proxy: {
    '/api/employee': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api\/employee/, '/employee'),
    },
    '/api/curriculum': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api\/curriculum/, '/curriculum'),
    },
    '/api/student-group': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api\/student-group/, '/student-group'),
    },
    '/api/student': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api\/student/, '/student'),
    },
    '/api/teacher-load': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api\/teacher-load/, '/teacher-load'),
    },
    '/api/asset': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api\/asset/, '/asset'),
    },
    '/api/schedule': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api\/schedule/, '/schedule'),
    },
  }
}
})