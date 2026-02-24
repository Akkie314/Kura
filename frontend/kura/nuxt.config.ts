// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  css: ['~/assets/css/main.css'],

  nitro: {
    devProxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        secure: false,
      }
    }
  },

  modules: ["@nuxtjs/google-fonts"],
  
  runtimeConfig: {
    public: {
      baseUrl: process.env.NUXT_PUBLIC_BASE_URL || 'http://localhost:3000'
    }
  },

  googleFonts: {
    families: {
      'Noto Sans JP': [400, 700]
    }
  },
})