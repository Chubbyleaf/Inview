<template>
    <div>
        <TaskConfig v-if="isTaskConfig" @back="backHandler" :configData="configData" :isAddTask="isAddTask"/>
        <TaskList v-else @show-config="showConfigHandler" />
    </div>
</template>

<script>
import TaskList from './TaskList.vue'
import TaskConfig from './TaskConfig.vue'


export default {
    name: 'TaskManagement',
    components: {
        TaskList,
        TaskConfig,
    },
    data() {
        return {
            isTaskConfig: false,
            isAddTask:false,
            configData: {},
        };
    },
    methods: {
        showConfigHandler(data,isAddTask) {
            this.configData = data; // 接收传递的数据
            this.isAddTask = isAddTask
            console.log(this.configData)
            this.isTaskConfig = true;
        },
        backHandler() {
            this.isTaskConfig = false;
            this.isAddTask = false;
            this.configData = {};
        }

    }
}
</script>
