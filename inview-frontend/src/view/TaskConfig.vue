<template>
    <el-row class="back-button-container">
        <el-button @click="$emit('back')">返回</el-button>
    </el-row>
    <el-row>
        <el-steps :active="active" align-center>
            <el-step title="截取图像" />
            <el-step title="区域绘制" />
            <el-step title="配置信息" />

        </el-steps>
    </el-row>
    <el-row class="step-content">
        <el-form v-if="active == 0" :model="currentCamera" :rules="rules" ref="currentTask" label-width="120px"
            class="add-task-form">
            <el-form-item label="摄像头Id" prop="deviceId">
                <el-input v-model="currentCamera.deviceId" disabled>
                </el-input>
            </el-form-item>
            <el-form-item label="摄像头名称" prop="name">
                <el-input v-model="currentCamera.name" disabled>
                </el-input>
            </el-form-item>
            <el-form-item label="摄像头备注" prop="remark">
                <el-input v-model="currentCamera.remark" disabled>
                </el-input>
            </el-form-item>
            <el-form-item prop="url" label="摄像头画面">
                <div class="camera-view-container">
                    <LivePlayer :videoUrl="configData.url" autoplay="true" muted="false" @error="handlePlayerError"
                        @snapOutside="snapOutside" ref="livePlayer" />
                </div>
                <div class="shot-button-container">
                    <el-button @click="captureScreenshot">截取</el-button>
                </div>
                <div class="screenshot-container">
                    <div v-if="!screenshotData">暂无截图</div>
                    <img v-else :src="screenshotData" alt="Screenshot" />
                </div>
            </el-form-item>
        </el-form>
        <el-row class="step-content" v-else-if="active == 1">
            <div style="width:100%;color: #939aec;font-weight: bold;">点击【绘制区域按钮】可在画面中绘制区域（再次点击可绘制多个区域）</div>
            <draw-area :screenshotData = "screenshotData" :coordinates="currentTask.coordinate" ref="drawArea" />
        </el-row>
        <el-form v-else-if="active == 2" :model="currentTask" :rules="rules" ref="currentTask" label-width="200px"
            class="add-task-form">
            <el-form-item label="模型" prop="model">
                <el-select v-model="currentTask.model" placeholder="请选择模型" :disabled="!isAddTask">
                    <el-option v-for="item in modelOptions" :key="item.value" :label="item.label" :value="item.value">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="细分模型" prop="algorithmType">
                <el-select v-model="currentTask.algorithmType" placeholder="请选择细分模型" :disabled="!isAddTask">
                    <el-option v-for="item in algorithmTypeOptions" :key="item.value" :label="item.label"
                        :value="item.value">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item label="Log周期" prop="logRotateDuration">
                <el-radio-group v-model="currentTask.logRotateDuration">
                    <el-radio :label="3">3天</el-radio>
                    <el-radio :label="7">7天</el-radio>
                    <el-radio :label="15">15天</el-radio>
                    <el-radio :label="30">30天</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="Log大小(MB)" prop="logSize">
                <el-input v-model.number="currentTask.logSize"></el-input>
            </el-form-item>
            <el-form-item label="本地保存周期" prop="dataRotation">
                <el-radio-group v-model="currentTask.dataRotation">
                    <el-radio :label="15">15天</el-radio>
                    <el-radio :label="30">30天</el-radio>
                    <el-radio :label="60">60天</el-radio>
                    <el-radio :label="90">90天</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item label="本地保存数据大小(GB)" prop="dataSize">
                <el-input v-model.number="currentTask.dataSize"></el-input>
            </el-form-item>
            <el-form-item label="工作时间段" prop="workingTime">
                <div v-for="(time, index) in formattedWorkingTime" :key="index" class="time-period">
                    <el-time-picker v-model="time.startTime" format="HH:mm:ss" value-format="HH:mm:ss"
                        @change="updateWorkingTime(index, 'startTime', time.startTime)"></el-time-picker>
                    <el-time-picker v-model="time.endTime" format="HH:mm:ss" value-format="HH:mm:ss"
                        @change="updateWorkingTime(index, 'endTime', time.endTime)"></el-time-picker>
                    <el-button @click="removeTimePeriod(index)">删除</el-button>
                </div>
                <el-button @click="addTimePeriod">添加时间段</el-button>
            </el-form-item>
            <el-form-item label="推送API地址" prop="targetApi">
                <el-input v-model="currentTask.targetApi"></el-input>
            </el-form-item>
            <el-form-item label="喇叭" prop="sound">
                <el-switch v-model="currentTask.sound"></el-switch>
            </el-form-item>
            <el-form-item label="短信" prop="sms">
                <el-switch v-model="currentTask.sms"></el-switch>
            </el-form-item>
            <el-form-item label="电话" prop="tel">
                <el-switch v-model="currentTask.tel"></el-switch>
            </el-form-item>
        </el-form>
    </el-row>
    <el-row class="step-content">

    </el-row>
    <el-row class="step-content">
        <div>
            <el-button v-if='active !== 0' @click="last" class="step-button">上一步</el-button>
            <el-button v-if='active !== 2' @click="next" class="step-button">下一步</el-button>
            <el-button v-if='active === 2' @click="saveTask()" class="step-button">完成配置</el-button>
        </div>

    </el-row>
</template>

<script>
import { getCameraList, addTask, editTask } from '../js/request';
import LivePlayer from '@liveqing/liveplayer-v3'
import DrawArea from './DrawArea.vue'

export default {
    name: 'AddTask',
    async mounted() {
        this.cameraList = await getCameraList();
        this.resetForm();
        this.initializeFormattedWorkingTime()
    },
    components: {
        LivePlayer,
        DrawArea
    },
    props: {
        configData: {
            type: Object,
        },
        isAddTask: {
            type: Boolean,
        }
    },
    methods: {
        // 下一步
        next() {
            if (this.active == 0 && !this.screenshotData) {
                this.$message.warning("请先截图！！")
            } else if (this.active == 1) {
                this.currentTask.coordinate = this.$refs.drawArea.savePoints();
                if (this.currentTask.coordinate.length == 0) {
                    this.$message.warning("请绘制区域!")
                } else { this.active++; }
            } else if (this.active++ > 2) {
                this.active = 0
            }
        },
        last() {
            if (this.active-- < 0) this.active = 0
        },
        //时间选择
        formatTimeWithColons(time) {
            if (!time) return '';
            const parts = time.match(/.{1,2}/g) || [];
            return parts.join(':');
        },
        removeColonsFromTime(time) {
            if (time) return time.replace(/:/g, '');
            else return "";
        },
        updateWorkingTime(index, type, value) {
            this.formattedWorkingTime[index][type] = value;
            this.currentTask.workingTime[index][type] = this.removeColonsFromTime(value);
        },
        addTimePeriod() {
            this.currentTask.workingTime.push({ startTime: '', endTime: '' });
            this.formattedWorkingTime.push({ startTime: '', endTime: '' });
        },
        removeTimePeriod(index) {
            this.currentTask.workingTime.splice(index, 1);
            this.formattedWorkingTime.splice(index, 1);
        },
        initializeFormattedWorkingTime() {
            this.formattedWorkingTime = this.currentTask.workingTime.map(period => {
                return {
                    startTime: this.formatTimeWithColons(period.startTime),
                    endTime: this.formatTimeWithColons(period.endTime)
                };
            });
        },
        //屏幕截图
        captureScreenshot() {
            const livePlayer = this.$refs.livePlayer;
            if (livePlayer && typeof livePlayer.snap === 'function') {
                livePlayer.snap();
            }
        },
        snapOutside(snapData) {
            if (snapData == 'data:,') {
                this.$message.warning("设备暂无画面,无法截图,请先检查设备！")
            } else {
                this.screenshotData = snapData
            }
        },
        resetForm() {
            this.currentCamera = this.configData
            this.currentTask = {
                deviceId: this.configData.deviceId,
                model: this.isAddTask ? '' : this.configData.model,
                algorithmType: this.isAddTask ? '' : this.configData.algorithmType,
                liveStreamInputUrl: this.configData.url,
                liveStreamOutputUrl: '',
                coordinate: this.isAddTask ? [] : this.configData.coordinate,
                logRotateDuration: this.isAddTask ? 3 : this.configData.logRotateDuration,///默认3天
                logSize: this.isAddTask ? 100 : this.configData.logSize,//默认100MB
                dataSize: this.isAddTask ? 50 : this.configData.dataSize,//默认50GB
                dataRotation: this.isAddTask ? 30 : this.configData.dataRotation,//默认30天
                workingTime: this.isAddTask ? [] : this.configData.workingTime,
                targetApi: this.isAddTask ? "" : this.configData.targetApi,
                gpu: 'NA',
                sound: this.isAddTask ? false : this.configData.sound,
                sms: this.isAddTask ? false : this.configData.sms,
                tel: this.isAddTask ? false : this.configData.tel,
            };
            if (!this.isAddTask) {
                this.currentTask._id = this.configData.taskId
            }
        },

        //校验时间段是否正确
        validateTimePeriods(rule, value, callback) {
            const workingTime = value;
            for (let i = 0; i < workingTime.length; i++) {
                const period = workingTime[i];
                if (!period.startTime || !period.endTime) {
                    return callback(new Error(`时间段 ${i + 1} 的开始时间或结束时间不能为空`));
                }
                if (parseInt(period.startTime) >= parseInt(period.endTime)) {
                    return callback(new Error(`时间段 ${i + 1} 的开始时间必须早于结束时间`));
                }
                for (let j = i + 1; j < workingTime.length; j++) {
                    const nextPeriod = workingTime[j];
                    if (
                        (parseInt(period.startTime) < parseInt(nextPeriod.endTime) && parseInt(period.endTime) > parseInt(nextPeriod.startTime)) ||
                        (parseInt(nextPeriod.startTime) < parseInt(period.endTime) && parseInt(nextPeriod.endTime) > parseInt(period.startTime))
                    ) {
                        return callback(new Error(`时间段 ${i + 1} 和时间段 ${j + 1} 有重叠`));
                    }
                }
            }
            callback();
        },

        async saveTask() {
            this.$refs.currentTask.validate(async (valid) => {
                if (valid) {
                    console.log("currentTask", this.currentTask)
                    let res;
                    if (this.isAddTask) {
                        res = await addTask(this.currentTask);
                    } else {
                        res = await editTask(this.currentTask)

                    }
                    if (res && res.data.status == 200) {
                        this.resetForm();
                        this.$emit('back')
                    }

                }
            });
        }
    },
    data() {
        return {
            screenshotData: null,//屏幕截图
            active: 0,
            cameraList: [],
            currentTask: {},
            currentCamera: {},
            formattedWorkingTime: [],
            canEdit: true,
            preDrawnPoints:[],//之前绘制的点
            modelOptions: [
                { value: 'safety', label: '安全' },
            ],
            algorithmTypeOptions: [
                { value: 'equipment', label: '安全帽/服' },
                { value: 'concreteSupport', label: '砼支撑' },
                { value: 'areaEdge', label: '临边' },
                { value: 'fireSmoke', label: '明火烟雾' }
            ],
            gpuOptions: [
                { value: '0', label: '0' },
                { value: '1', label: '1' },
                { value: 'NA', label: 'NA' },
            ],
            rules: {
                model: [
                    { required: true, message: '请选择算法大类', trigger: 'change' }
                ],
                algorithmType: [
                    { required: true, message: '请选择算法小类', trigger: 'change' }
                ],
                logRotateDuration: [
                    { required: true, message: '请选择Log周期', trigger: 'change' }
                ],
                logSize: [
                    { required: true, message: '请输入Log大小', trigger: 'blur' },
                    { type: 'number', message: 'Log大小必须是数字' },
                ],
                dataRotation: [
                    { required: true, message: '请选择本地保存周期', trigger: 'change' }
                ],
                dataSize: [
                    { required: true, message: '请输入保存数据大小', trigger: 'blur' },
                    { type: 'number', message: '保存数据大小必须是数字' }
                ],
                workingTime: [
                    { required: true, type: 'array', message: '请添加工作时间段', trigger: 'change' },
                    { required: true, type: 'array', validator: this.validateTimePeriods, trigger: 'change' }
                ],
                sound: [
                    { required: true, message: '请选择喇叭状态', trigger: 'change' }
                ],
                sms: [
                    { required: true, message: '请选择短信状态', trigger: 'change' }
                ],
                tel: [
                    { required: true, message: '请选择电话状态', trigger: 'change' }
                ],
                // gpu: [
                //     { required: true, message: '请选择GPU', trigger: 'change' }
                // ],
            },
        };
    }
};
</script>

<style src="../style/Task.css" scoped></style>