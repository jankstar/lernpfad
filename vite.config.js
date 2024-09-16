// vite.config.js
import { defineConfig } from 'vite'

export default defineConfig({
  base: "/lernpfad/",
  build: {
    rollupOptions: {
      external: [
        /^node:.*/,
      ]
    }
  }
})
