import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'

import { tanstackStart } from '@tanstack/react-start/plugin/vite'

import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// The Go API (`cv_gen serve`) runs on :8080 by default; proxy /api to it in
// dev so the frontend can call same-origin /api/* without CORS. Override the
// target with CV_GEN_API when the server runs elsewhere.
const apiTarget = process.env.CV_GEN_API ?? 'http://localhost:8080'

const config = defineConfig({
  resolve: { tsconfigPaths: true },
  plugins: [devtools(), tailwindcss(), tanstackStart(), viteReact()],
  server: {
    proxy: {
      '/api': { target: apiTarget, changeOrigin: true },
    },
  },
})

export default config
