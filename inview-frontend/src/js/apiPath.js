
export const API_PATHS = {
    //用户相关接口
    LOGIN: `/user/login`,//用户登录
    REGISTER: `/user/register`,//用户注册
    CLOSE_ACCOUNT: `/user/delete`,//用户注销
    GET_USER_INFO: `/user/info`,//获取用户信息
    EDIT_USER_INFO: `/user/editUserInfo`,//修改账户信息
    EDIT_USER_PASSWORD: `/user/editpsw`,//修改密码
    UPDATE_USER_INFO: `/user/update`,//更新用户信息

    //设备相关接口
    GET_CAMERA_LIST: `/camera/getCameraList`,//获取相机列表
    ADD_CAMERA: `/camera/add`,//增加摄像机
    DELETE_CAMERA: `/camera/delete`,//删除摄像机
    EDIT_CAMERA: `/camera/editCameraInfo`,//修改摄像机

    //任务相关接口
    GET_TASK_RESULT: '/algorithm/getTaskResultList',//获取最新任务结果列表
    GET_TASK_LIST:'/algorithm/getTaskList',//获取任务列表
    EDIT_TASK:'/algorithm/editTask',//修改任务
    DELETE_TASK:'/algorithm/deleteTask',//删除任务
    ADD_TASK:'/algorithm/addTask',//增加任务
    STOP_TASK:'/algorithm/stopTask',//手动停止任务
    START_TASK:'/algorithm/startTask',//手动启动任务

    //获取系统信息
    GET_SYSTEM_INFO:'/system/getInfo',
    GET_SYSTEM_LOG:'/system/getLog',
    GET_SYSTEM_IP:'/system/getIp'

};
