import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css';
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import './style/App.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const app = createApp(App)


app.use(router)
app.use(ElementPlus, {
    locale: zhCn,
  })
app.mount('#app')
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}