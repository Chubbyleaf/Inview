import http from './interceptor';
import { API_PATHS } from './apiPath'

export async function login(req) {
    const response = await http.post(API_PATHS.LOGIN, req);
    if (response && response.status === 200) {
        // 保存用户ID到localStorage
        localStorage.setItem('user_id', response.data.data._id);
        localStorage.setItem('username', response.data.username);
        return response;
    }

}

export async function register(req) {
    const response = await http.post(API_PATHS.REGISTER, req);
    return response;

}

export async function closeAccount() {
    const userId = localStorage.getItem('user_id');
    const response = await http.delete(`${API_PATHS.CLOSE_ACCOUNT}/${userId}`);
    return response;

}

export async function getUserInfo() {
    const userId = localStorage.getItem('user_id');
    const response = await http.get(`${API_PATHS.GET_USER_INFO}/${userId}`);
    return response;

}

export async function editPassword(req) {
    const response = await http.post(`${API_PATHS.EDIT_USER_PASSWORD}`, req);
    return response;
}


export async function editInfo(req) {
    const response = await http.post(`${API_PATHS.EDIT_USER_INFO}`, req);
    return response;
}

export async function getCameraList() {
    const response = await http.get(API_PATHS.GET_CAMERA_LIST);
    if (response && response.status === 200) {
        return response.data.data;
    } else {
        return [];
    }
}

export async function addCamera(req) {
    await http.post(API_PATHS.ADD_CAMERA, req);

}

export async function editCamera(req) {
    await http.post(API_PATHS.EDIT_CAMERA, req);

}

export async function deleteCamera(id) {
    await http.delete(`${API_PATHS.DELETE_CAMERA}/${id}`);
}

export async function getTaskResult(req) {
    const response = await http.post(API_PATHS.GET_TASK_RESULT, req);
    if (response && response.status === 200) {
        return response.data.data;
    } else {
        return {};
    }
}

export async function getTaskList() {
    const response = await http.get(API_PATHS.GET_TASK_LIST);
    if (response && response.status === 200) {
        return response.data.data;
    } else {
        return [];
    }
}

export async function getTaskInfo(taskId) {
    const response = await http.get(`${API_PATHS.GET_TASK_INFO}/${taskId}`);
    return response;
}

export async function addTask(req) {
    const response = await http.post(API_PATHS.ADD_TASK, req);
    return response;
}

export async function stopTask(taskId) {
    const response = await http.post(`${API_PATHS.STOP_TASK}/${taskId}`);
    return response;
}

export async function startTask(taskId) {
    const response = await http.post(`${API_PATHS.START_TASK}/${taskId}`);
    return response;
}
export async function deleteTaskById(taskId) {
    await http.delete(`${API_PATHS.DELETE_TASK}/${taskId}`);
}

export async function editTask(req) {
    const response = await http.post(API_PATHS.EDIT_TASK, req);
    return response;

}
export async function getSystemInfo() {
    const response = await http.get(`${API_PATHS.GET_SYSTEM_INFO}`);
    if (response && response.status === 200) {
        return response.data.data;
    } else {
        return {};
    }
}

export async function getSystemLog() {
    const response = await http.get(`${API_PATHS.GET_SYSTEM_LOG}`);
    if (response && response.status === 200) {
        return response.data.data;
    } else {
        return "";
    }
}

export async function getSystemIp() {
    const response = await http.get(`${API_PATHS.GET_SYSTEM_IP}`);
    if (response && response.status === 200) {
        return response.data.data;
    } else {
        return "";
    }
}
