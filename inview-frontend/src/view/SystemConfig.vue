<template>
    <div>
        <el-descriptions :column=1>
            <el-descriptions-item label="设备ID">{{ systemInfo.machineID }}</el-descriptions-item>
            <el-descriptions-item label="设备类型">{{ systemInfo.machineType }}</el-descriptions-item>
            <el-descriptions-item label="设备名称">{{ systemInfo.machineName }}</el-descriptions-item>
            <el-descriptions-item label="软件版本">{{ systemInfo.softwareVersion }}</el-descriptions-item>
            <el-descriptions-item label="软核版本">{{ systemInfo.firmwareVersion }}</el-descriptions-item>
            <el-descriptions-item label="网络名称">{{ systemInfo.interfaceName }}</el-descriptions-item>
            <el-descriptions-item label="IP地址">{{ systemInfo.ip }}</el-descriptions-item>
            <el-descriptions-item label="MAC地址">{{ systemInfo.mac }}</el-descriptions-item>
            <el-descriptions-item label="恢复出厂设置">
                <el-button size="small" @click="showConfirmDialog">恢复</el-button>
            </el-descriptions-item>
        </el-descriptions>
        <!-- 确认恢复出厂设置对话框 -->
        <el-dialog title="确认" v-model="resetDialogVisible" width="30%">
            <span>是否确认恢复出厂设置？</span>
            <template v-slot:footer>
                <span class="dialog-footer">
                    <el-button class="text-button" @click="resetDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="confirmFactoryReset">确认</el-button>
                </span>
            </template>
        </el-dialog>
    </div>

</template>

<script>
import { getSystemIp } from '../js/request';

export default {
    name: 'SystemConfig',
    methods: {
        //恢复出厂设置
        showConfirmDialog() {
            this.resetDialogVisible = true;
        },
        confirmFactoryReset() {
            this.resetDialogVisible = false;
        },
        async initialInfo(){
            const res = await getSystemIp()
            if(res){
                this.systemInfo.machineID = res.machineID
                this.systemInfo.mac = res.mac
                this.systemInfo.ip =res.ip
                this.systemInfo.interfaceName = res.interfaceName
            }
        }
    },
    mounted(){
        this.initialInfo()
    },
    data() {
        return {
            resetDialogVisible: false,
            systemInfo: {
                machineID: "暂无",
                machineType: "边缘盒子",
                machineName: "创乐科技Insense-Inview",
                softwareVersion: "v1.0",
                firmwareVersion: "v1.0",
                interfaceName:"暂无信息",
                ip: "暂无信息",
                mac: "暂无信息"

            }
        }
    }
}
</script>

<style src="../style/Config.css"></style>