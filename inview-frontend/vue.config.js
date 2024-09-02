const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  publicPath:"./",
  parallel:false,
  assetsDir:"static",
  devServer: {
    port: 8088, 
    host: 'localhost',
    proxy: {
      "/api": {
        target: process.env.VUE_APP_TARGET_API, // 请求域名
        secure: false, // 请求是否为https
        changeOrigin: true, // 是否跨域
        pathRewrite: { "^/api": "/api" }
      }
    }
  },
})
