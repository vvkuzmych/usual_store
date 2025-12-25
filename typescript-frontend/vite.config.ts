import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    port: 3001,
    host: true,
    proxy: {
      "/api": {
        target: "http://back-end:4001",
        changeOrigin: true,
      },
    },
  },
  preview: {
    port: 3001,
    host: true,
  },
});

