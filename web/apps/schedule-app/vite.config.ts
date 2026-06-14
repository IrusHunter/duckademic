import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
      name: 'scheduleApp',
      filename: 'remoteEntry.js',
      exposes: {
        './ScheduleApp': './src/components/App/App.tsx',
      },
      shared: {
        react: { singleton: true, requiredVersion: '^19.0.0' } as any,
        'react-dom': { singleton: true, requiredVersion: '^19.0.0' } as any,
        '@tanstack/react-query': { singleton: true } as any,
        axios: { singleton: true } as any,
      },
    }),
  ],
  build: { target: 'esnext' },
  server: { port: 5008 },
  preview: { port: 5008 },
})
