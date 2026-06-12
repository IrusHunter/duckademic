import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'

export default defineConfig({
  plugins: [
    react(),
    federation({
      name: 'homeApp',
      filename: 'remoteEntry.js',
      // Точка входу переїхала в components/App/App.tsx
      exposes: {
        './HomeApp': './src/components/App/App.tsx',
      },
      shared: {
        react: { singleton: true, requiredVersion: '^19.0.0' } as any,
        'react-dom': { singleton: true, requiredVersion: '^19.0.0' } as any,
        // axios як singleton — спільний із shell (один екземпляр на сторінку)
        axios: { singleton: true } as any,
      },
    }),
  ],
  build: { target: 'esnext' },
  server: { port: 5006 },
  preview: { port: 5006 },
})
