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
        axios: { singleton: true } as any,
      }
    })
  ],
  build: { target: 'esnext' },
  server: {
  port: 5000,
  proxy: {
    '/api/auth': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''), // /api/auth/login → /auth/login
    },
    '/api/employee': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''), // /api/employee/teachers → /employee/teachers
    },
    '/api/curriculum': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
    '/api/student-group': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
    '/api/student': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
    '/api/teacher-load': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
    '/api/asset': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
    '/api/schedule': {
      target: 'http://localhost:10000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
  }
}
})