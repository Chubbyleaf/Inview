import axios from 'axios';
import { ElMessage } from 'element-plus';

// 创建axios实例
const http = axios.create({
    // 基础URL，所有请求都会带上/api前缀
    baseURL: '/api',
    timeout: 5000,
    headers: {
        'Content-Type': 'application/json; charset=utf-8'
    }
});

// 请求拦截器
http.interceptors.request.use((config) => {
    config.headers = config.headers || {};
    // 检查本地是否存储了用户ID
    const userId = sessionStorage.getItem("user_id");
    if (userId) {
        config.headers.user_id = userId;
    }
    return config;
});

// 响应拦截器
http.interceptors.response.use(
    (response) => {
        if (response && response.status === 200) {
            ElMessage({
                message: '成功！', 
                type: 'success',
                grouping: true,
            });
            return response;

        } else {
            ElMessage({
                message: response.data.message || 'An error occurred', 
                type: 'error',
            });
            return null;
        }
    },
    (error) => {
        ElMessage({
            message: error.response.data,
            type: 'error',
        });
        return null;
    }
);


/**
 * 响应处理
 * @param {Object} error 
 * @returns 
 */
// const ResponseProcessing = (error) => {
//     if (error.response) {
//         if (error.response.data.status !== 200) {
//             console.log(error.response)
//             ElMessage.warning(error.response.data.message); // 返回接口返回的错误信息
//         }
//     } else {
//         ElMessage.error("遇到跨域错误，请设置代理或者修改后端允许跨域访问！");
//     }
// };

export default http;
