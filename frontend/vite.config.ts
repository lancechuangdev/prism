import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/prism-api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/prism-api/, ''),
      },
    },
  },
  test: {
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    css: true,
    env: {
      VITE_PRISM_API_URL: 'http://localhost:8080',
      VITE_PRISM_CHAIN_ID: '31337',
      VITE_PRISM_CHAIN_NAME: 'Hardhat Local',
      VITE_PRISM_RPC_URL: 'http://127.0.0.1:8545',
      VITE_PRISM_POOL_ADDRESS: '0x0000000000000000000000000000000000000001',
      VITE_PRISM_MULTISIG_ADDRESS: '0x0000000000000000000000000000000000000002',
      VITE_PRISM_DEPLOYMENT_BLOCK: '0',
    },
  },
})
