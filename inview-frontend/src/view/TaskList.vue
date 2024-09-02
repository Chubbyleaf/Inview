<template>
    <div>
        <el-row id="header-row">
            <el-button @click="addCameraDialogVisible = true">
                <el-icon>
                    <Plus />
                </el-icon>
                添加摄像头
            </el-button>
            <el-select v-model="filter.deviceId" placeholder="请选择摄像头" class="select-style" clearable
                @change="handlePagination()">
                <el-option v-for="item in cameraList" :key="item.deviceId" :label="item.deviceId"
                    :value="item.deviceId">
                    <span style="float: left">{{ item.deviceId }}</span>
                    <span style="float: right;color: var(--el-text-color-secondary);font-size: 13px;">
                        {{ item.remark }}
                    </span>
                </el-option>
            </el-select>
            <el-select v-model="filter.model" placeholder="模型" class="select-style" clearable
                @change="handlePagination()">
                <el-option label="安全" value="safety"></el-option>
            </el-select>
            <el-select v-model="filter.algorithmType" placeholder="子模型" class="select-style" clearable
                @change="handlePagination()">
                <el-option label="安全帽/服" value="equipment"></el-option>
                <el-option label="砼支撑" value="concreteSupport"></el-option>
                <el-option label="临边" value="areaEdge"></el-option>
                <el-option label="明火烟雾" value="fireSmoke"></el-option>
            </el-select>
            <el-button @click="handlePagination()">筛选</el-button>
            <el-button @click="resetFilter" class="border-button">重置</el-button>
        </el-row>

        <el-row>
            <el-table :data="paginatedData" :span-method="objectSpanMethod" style="width: 100%" border
                :cell-style="{ textAlign: 'center' }" :header-cell-style="{ 'text-align': 'center' }">
                <el-table-column fixed="left" prop="deviceId" label="摄像头" width="200">
                    <template v-slot="scope">
                        <el-tooltip effect="dark" :content='scope.row.remark' placement="top">
                            {{ `${scope.row.deviceId}-${scope.row.name}` }}
                        </el-tooltip>
                        <div v-if="scope.$index === getFirstRowIndex(scope.row.deviceId)" style="margin-top: 5px">
                            <el-tooltip effect="dark" content="修改摄像头" placement="top">
                                <el-button size="small" class="border-button" @click="handleEditCamera(scope.row)">
                                    <el-icon>
                                        <Edit />
                                    </el-icon>
                                </el-button>
                            </el-tooltip>
                            <el-popconfirm width="220" confirm-button-text="确认删除" confirm-button-type="danger"
                                cancel-button-text="不用了，谢谢" @confirm="submitDeleteCamera(scope.row.deviceId)"
                                icon-color="#626AEF" title="您确定要删除此摄像头？">
                                <template #reference>
                                    <el-button size="small" class="border-button">
                                        <el-icon>
                                            <Delete />
                                        </el-icon>
                                    </el-button>
                                </template>
                            </el-popconfirm>
                            <el-tooltip effect="dark" content="为本摄像头添加新任务" placement="top">
                                <el-button size="small" class="border-button" @click="showConfig(scope.row, true)">
                                    <el-icon>
                                        <Plus />
                                    </el-icon>
                                </el-button>
                            </el-tooltip>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column prop="model" label="模型类别">
                    <template v-slot="scope">
                        {{ scope.row.model ? modelMap[scope.row.model] : '-' }}
                    </template>
                </el-table-column>
                <el-table-column prop="algorithmType" label="子模型">
                    <template v-slot="scope">
                        {{ scope.row.algorithmType ? algorithmTypeMap[scope.row.algorithmType] : '-' }}
                    </template>
                </el-table-column>
                <el-table-column prop="status" label="运行状态" width="120">
                    <template v-slot="scope">
                         <!-- -1表示运行失败 0表示空  1运行成功 2用户手动停止 3表示不在运行时间范围内 -->
                        <el-icon v-if="scope.row.status == 1" color="green">
                            <CircleCheck />
                        </el-icon>
                        <el-icon v-else-if="scope.row.status == -1" color="red">
                            <CircleClose />
                        </el-icon>
                        <span v-else-if="scope.row.status == 0">-</span>
                        <span v-else-if="scope.row.status == 3">不运行时间内</span>
                        <span v-else-if="scope.row.status == 2 ">已手动停止</span>
                    </template>
                </el-table-column>
                <el-table-column prop="url" label="摄像头画面" width="120">
                    <template v-slot="scope">
                        <span v-if="scope.row.url">
                            <img id='play-video-icon' src="../assets/PlayVideo.png" @click="handleView(scope.row)" />
                        </span>
                        <span v-else>-</span>
                    </template>
                </el-table-column>
                <el-table-column fixed="right" label="操作" width="400">
                    <template v-slot="scope">
                        <el-button class="text-button" size="small" @click="showConfig(scope.row, false)"
                            :disabled='!scope.row.model'>修改算法配置</el-button>
                        <!-- 用户手动停止了需要启动按钮 -->
                        <el-button class="text-button" size="small" v-if="scope.row.status != 1"
                            :disabled='!scope.row.model' @click="startTask(scope.row.taskId)">启动算法</el-button>
                        <el-button class="text-button" size="small" v-else :disabled='!scope.row.model'
                            @click="stopTask(scope.row.taskId)">停止算法</el-button>
                        <el-button class="text-button" size="small" @click="getResultInfo(scope.row)"
                            :disabled='!scope.row.model'>查看结果</el-button>
                        <el-popconfirm width="220" confirm-button-text="确认删除" confirm-button-type="danger"
                            cancel-button-text="不用了，谢谢" @confirm="deleteTask(scope.row.taskId)" icon-color="#626AEF"
                            title="您确定要删除此任务？">
                            <template #reference>
                                <el-button class="text-button" size="small"
                                    :disabled='!scope.row.model'>删除任务</el-button>
                            </template>
                        </el-popconfirm>
                    </template>
                </el-table-column>
            </el-table>
        </el-row>
        <el-row>
            <el-pagination background small @current-change="handleCurrentChange" :current-page="currentPage"
                :page-size="pageSize" layout="total, prev, pager, next" :total="taskTotalData.length">
            </el-pagination>
        </el-row>

        <!-- 添加摄像头表单 -->
        <el-dialog title="添加摄像头" v-model="addCameraDialogVisible">
            <el-form :model="form" :rules="cameraRules" ref="form" label-width="120px">
                <el-form-item label="摄像头名称" prop="name">
                    <el-input v-model="form.name" placeholder="请输入摄像头名称"></el-input>
                </el-form-item>
                <el-form-item label="备注">
                    <el-input v-model="form.remark" placeholder="选填，字符不建议超过100个字"></el-input>
                </el-form-item>
                <el-form-item label="摄像头类型" prop="type">
                    <el-radio-group v-model="form.type">
                        <el-radio value="RTSP">网络摄像头(RTSP)</el-radio>
                        <el-radio value="USB">本地摄像头(USB)</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="选择摄像头" prop="url">
                    <div class="camera-row">
                        <el-input v-model="form.url" placeholder="请填写摄像头链接" style="width:80%"></el-input>
                        <el-button style="margin-left: 10px;" @click="checkCameraUrl">检验</el-button>
                    </div>
                </el-form-item>
                <el-form-item label="摄像头查看">
                    <div class="camera-view-container">
                        <LivePlayer :videoUrl="newCameraUrl" autoplay="true" muted="false" ref="checkCameraPlayer" />
                    </div>
                </el-form-item>
            </el-form>

            <template v-slot:footer>
                <div class="dialog-footer">
                    <el-button class="text-button" @click="addCameraDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="submitAddCamera">确认</el-button>
                </div>
            </template>
        </el-dialog>


        <!-- 修改摄像头表单 -->
        <el-dialog title="修改摄像头" v-model="editCameraDialogVisible">
            <el-form :model="editCameraForm" :rules="cameraRules" ref="editCameraForm" label-width="120px">
                <el-form-item label="摄像头名称" prop="name">
                    <el-input v-model="editCameraForm.name" placeholder="请输入摄像头名称"></el-input>
                </el-form-item>
                <el-form-item label="备注">
                    <el-input v-model="editCameraForm.remark" placeholder="选填，字符不建议超过100个字"></el-input>
                </el-form-item>
                <el-form-item label="摄像头类型" prop="type">
                    <el-radio-group v-model="editCameraForm.type" disabled>
                        <el-radio value="RTSP">网络摄像头(RTSP)</el-radio>
                        <el-radio value="USB">本地摄像头(USB)</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="摄像头链接" prop="url">
                    <el-input v-model="editCameraForm.url" style="width:80%" disabled></el-input>
                </el-form-item>
            </el-form>
            <template v-slot:footer>
                <div class="dialog-footer">
                    <el-button class="text-button" @click="editCameraDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="submitEditCamera">确认</el-button>
                </div>
            </template>
        </el-dialog>

        <!-- 查看摄像头 -->
        <el-dialog v-model="viewCameraDialogVisible" width="70%" @open="handlePlayerOpen">
            <LivePlayer ref="livePlayer" :videoUrl="currentCameraUrl" autoplay="true" muted="false"
                @error="handlePlayerError" />
        </el-dialog>

        <!-- 查看结果 -->
        <el-dialog v-model="resultDialogVisible" width="80%" @close="timeRange = []">
            <el-row v-if="currentResult" id="result-operations-container">
                <div>
                    <span class="margin-right-20">请选择时间范围</span>
                    <el-date-picker v-model="timeRange" type="daterange" unlink-panels range-separator="To"
                        start-placeholder="开始时间" end-placeholder="结束时间" class="margin-right-20"
                        @change="filterTaskResult" />
                    <el-button class="margin-right-20" @click="getResultInfo()">筛选</el-button>
                </div>

            </el-row>
            <div v-if="resultList">
                <div id='image-gallery'>
                    <el-row>
                        <el-col v-for="(image, index) in paginatedResultList" :key="index" :span="8"
                            style="text-align: center;">
                            <img :src="`${currentResultPath}/${image.imgName}`" alt="result image" class="result-image"
                                @click="zoomIn(image)" />
                        </el-col>
                    </el-row>
                </div>
                <el-pagination background small @current-change="handleImagePageChange" :current-page="currentImagePage"
                    :page-size="imagesPerPage" layout="total, prev, pager, next" :total="resultList.length">
                </el-pagination>
            </div>

            <div v-else>
                <div id="no-fund-text">暂无结果</div>
                <img src="../assets/NoResult.png" alt='no-result' />
            </div>

        </el-dialog>

        <!-- 查看放大结果 -->
        <el-dialog v-model="zoomInDialogVisible" width="60%" @open="drawAnnotations">
            <div id='annotation-canvas' v-if="currentResult">
                <canvas ref="annotationCanvas" style="border: 1px solid #000;"></canvas>
            </div>

        </el-dialog>

    </div>
</template>

<script>
import LivePlayer from '@liveqing/liveplayer-v3'
import { ElMessage } from 'element-plus';

import {
    getTaskList, deleteTaskById,
    addCamera, deleteCamera, editCamera, getCameraList,
    stopTask, startTask,
    getTaskResult
} from '../js/request';

export default {
    name: 'TaskList',
    async mounted() {
        await this.loadData()
    },
    components: {
        LivePlayer
    },
    methods: {
        //重置筛选条件
        resetFilter() {
            this.filter = {
                deviceId: '',
                model: '',
                algorithmType: ''
            };
            this.handlePagination();
        },

        // 增加摄像头
        async submitAddCamera() {
            this.$refs.form.validate(async (valid) => {
                if (valid) {
                    await addCamera(this.form);
                    this.resetForm();
                    this.addCameraDialogVisible = false;
                    await this.loadData()
                }
            });
        },

        //重置表单
        resetForm() {
            this.form = {
                name: '',
                type: '',
                remark: '',
                url: '',
            };
        },

        // 修改摄像头信息
        handleEditCamera(item) {
            this.editCameraForm = {
                deviceId: item.deviceId,
                name: item.name,
                remark: item.remark,
                type: item.type,
                url: item.url
            };
            this.editCameraDialogVisible = true;
        },

        async submitEditCamera() {
            this.$refs.editCameraForm.validate(async (valid) => {
                if (valid) {
                    await editCamera(this.editCameraForm);
                    this.editCameraDialogVisible = false;
                    await this.loadData()
                }
            });
        },

        // 删除摄像机
        async submitDeleteCamera(cameraId) {
            await deleteCamera(cameraId);
            this.deleteDialogVisible = false;
            await this.loadData()
        },


        // 查看摄像头
        handleView(camera) {
            this.currentCameraUrl = camera.url;
            this.viewCameraDialogVisible = true;
        },

        handlePlayerOpen() {
            this.$nextTick(() => {
                if (this.$refs.livePlayer) {
                    const videoElement = this.$refs.livePlayer.$el.querySelector('video');
                    videoElement.onloadeddata = () => {
                        videoElement.play().catch(error => {
                            console.error('Video playback failed', error);
                        });
                    };
                }
            });
        },

        //检查摄像头画面
        checkCameraUrl() {
            this.newCameraUrl = this.form.url;
            this.$nextTick(() => {
                if (this.$refs.checkCameraPlayer) {
                    const videoElement = this.$refs.checkCameraPlayer.$el.querySelector('video');
                    videoElement.onloadeddata = () => {
                        videoElement.play().catch(error => {
                            console.error('Video playback error:', error);
                            ElMessage({
                                message: '视频播放失败，请检查摄像头链接',
                                type: 'error'
                            });
                        });
                    };

                    videoElement.onerror = () => {
                        console.error('Video error event triggered');
                        ElMessage({
                            message: '无法加载视频，请检查摄像头链接',
                            type: 'error'
                        });
                    };
                }
            });
        },

        async deleteTask(taskId) {
            await deleteTaskById(taskId);
            await this.loadData();
        },

        //显示画框区域
        showArea(status) {
            this.showAreaBoundary = status;
            this.drawAnnotations();
        },

        //显示算法配置
        showConfig(row, isAddTask) {
            this.$emit('show-config', row, isAddTask); // 发出事件并传递当前行数据
        },

        //停止运行算法
        async stopTask(taskId) {
            await stopTask(taskId)
            this.taskTotalData = await getTaskList();
            await this.loadData()
        },

        //启动算法
        async startTask(taskId) {
            await startTask(taskId)
            await this.loadData()
        },

        async loadData() {
            this.taskTotalData = await getTaskList();
            this.cameraList = await getCameraList();
            this.handlePagination();
        },

        handlePagination() {
            const start = (this.currentPage - 1) * this.pageSize;
            const end = this.currentPage * this.pageSize;
            let filtered = this.taskTotalData.filter(item => {
                return (this.filter.deviceId ? item.deviceId == this.filter.deviceId : true)
                    && (this.filter.model ? item.model === this.filter.model : true)
                    && (this.filter.algorithmType ? item.algorithmType === this.filter.algorithmType : true);
            });
            this.paginatedData = filtered.slice(start, end);
        },

        handleCurrentChange(page) {
            this.currentPage = page;
            this.handlePagination();
        },

        getFirstRowIndex(deviceId) {
            for (let i = 0; i < this.paginatedData.length; i++) {
                if (this.paginatedData[i].deviceId === deviceId) {
                    return i;
                }
            }
            return -1;
        },

        //表头合并规则
        objectSpanMethod({ row, rowIndex, columnIndex }) {
            if (columnIndex === 1 || columnIndex === 0) {
                if (rowIndex === 0 || this.paginatedData[rowIndex].deviceId !== this.paginatedData[rowIndex - 1].deviceId) {
                    let rowspan = 1;
                    for (let i = rowIndex + 1; i < this.paginatedData.length; i++) {
                        if (this.paginatedData[i].deviceId === row.deviceId) {
                            rowspan++;
                        } else {
                            break;
                        }
                    }
                    return {
                        rowspan,
                        colspan: 1
                    };
                } else {
                    return {
                        rowspan: 0,
                        colspan: 0
                    };
                }
            }
        },

        // 自定义函数将日期转换为 YYYYMMDD 格式
        formatDateToYYYYMMDD(date) {
            const year = date.getFullYear();
            const month = (date.getMonth() + 1).toString().padStart(2, '0'); // 月份从 0 开始，所以加 1
            const day = date.getDate().toString().padStart(2, '0');
            return `${year}${month}${day}`;
        },

        //查看算法结果
        async getResultInfo(row) {
            if (row) {
                this.currentTask = row
            }
            let startDate = "";
            let endDate = "";

            if (this.timeRange && this.timeRange.length == 2) {
                startDate = this.formatDateToYYYYMMDD(new Date(this.timeRange[0].toLocaleString('en-US', { timeZone: 'Asia/Shanghai' })));
                endDate = this.formatDateToYYYYMMDD(new Date(this.timeRange[1].toLocaleString('en-US', { timeZone: 'Asia/Shanghai' })));
            }
            this.resultList = []

            const req = {
                deviceId: this.currentTask.deviceId,
                model: this.currentTask.model,
                algorithmType: this.currentTask.algorithmType,
                startTime: startDate,
                endTime: endDate
            }
            this.resultList = await getTaskResult(req)
            if (this.resultList) {
                if (process.env.VUE_APP_IMAGE_URL) {
                    this.currentResultPath = `${process.env.VUE_APP_IMAGE_URL}/${this.currentTask.deviceId}/${this.currentTask.taskId}/image`
                } else {
                    this.currentResultPath = `http://${window.location.hostname}/statistic/${this.currentTask.deviceId}/${this.currentTask.taskId}/image`
                }
                this.resultDialogVisible = true;
                this.handleImagePageChange(1)
            }
        },

        //算法结果分页
        handleImagePageChange(page) {
            this.currentImagePage = page;
            const start = (this.currentImagePage - 1) * this.imagesPerPage;
            const end = start + this.imagesPerPage;
            this.paginatedResultList = this.resultList.slice(start, end);
        },

        //放大图片
        zoomIn(result) {
            this.zoomInDialogVisible = true
            this.currentResult = result
        },

        drawAnnotations() {
            const canvas = this.$refs.annotationCanvas;
            if (canvas) {
                const ctx = canvas.getContext('2d');
                ctx.clearRect(0, 0, canvas.width, canvas.height);
                const img = new Image();
                img.src = `${this.currentResultPath}/${this.currentResult.imgName}`;
                img.onerror = () => {
                    this.$message.warning("未存储图片！")
                };

                img.onload = () => {
                    // 获取页面视高的60%和视宽的70%
                    const maxWidth = window.innerWidth * 0.7;
                    const maxHeight = window.innerHeight * 0.6;

                    // 计算图像的缩放比例
                    const scaleWidth = maxWidth / img.width;
                    const scaleHeight = maxHeight / img.height;
                    const scale = Math.min(scaleWidth, scaleHeight);

                    // 计算缩放后的图像尺寸
                    const scaledWidth = img.width * scale;
                    const scaledHeight = img.height * scale;

                    canvas.width = scaledWidth;
                    canvas.height = scaledHeight;

                    // 绘制缩放后的图像
                    ctx.drawImage(img, 0, 0, scaledWidth, scaledHeight);

                    this.currentResult.predResult.forEach((result) => {
                        const [classIndex, centerX, centerY, width, height] = result;
                        const className = this.currentResult.className[classIndex];

                        // 计算缩放后的标注位置和大小
                        const x = centerX * scaledWidth - (width * scaledWidth) / 2;
                        const y = centerY * scaledHeight - (height * scaledHeight) / 2;
                        const rectWidth = width * scaledWidth;
                        const rectHeight = height * scaledHeight;

                        // 绘制矩形框
                        ctx.strokeStyle = 'red';
                        ctx.lineWidth = 2;
                        ctx.strokeRect(x, y, rectWidth, rectHeight);

                        // 计算文字宽度和高度
                        ctx.font = 'bold 18px Arial'; // 设置为粗体
                        const textWidth = ctx.measureText(className).width;
                        const textHeight = 25; // 字体大小
                        const padding = 5;
                        const margin = 5; // 与红色矩形框的距离

                        // 绘制背景矩形
                        const backgroundX = x - padding;
                        const backgroundY = y - textHeight - padding - margin;
                        const backgroundWidth = textWidth + padding * 2;
                        const backgroundHeight = textHeight + padding;
                        ctx.fillStyle = 'red';
                        ctx.fillRect(backgroundX, backgroundY, backgroundWidth, backgroundHeight);

                        // 绘制文字
                        ctx.fillStyle = 'white';
                        const textX = x;
                        const textY = y - 10 - margin;
                        ctx.fillText(className, textX, textY);
                    });
                };
            }
        }


    },
    data() {
        return {
            currentResultPath: "",
            addCameraDialogVisible: false,//增加摄像头对话框
            editCameraDialogVisible: false,//修改摄像头对话框
            viewCameraDialogVisible: false,//查看摄像头对话框
            resultDialogVisible: false,//查看结果对话框
            zoomInDialogVisible: false,//放大显示结果的dialog
            filter: {
                deviceId: '',
                model: '',
                algorithmType: ''
            },
            form: {
                name: '',
                type: '',
                remark: '',
                url: ''
            },//增加摄像头表单
            editCameraForm: {},//修改摄像机信息表单
            //显示对饮关系
            modelMap: {
                safety: "安全",
            },
            algorithmTypeMap: {
                equipment: "安全帽/服",
                concreteSupport: "砼支撑",
                areaEdge: "临边",
                fireSmoke: "明火烟雾",
            },
            cameraRules: {
                name: [
                    { required: true, message: '请输入摄像头名称', trigger: 'blur' }
                ],
                type: [
                    { required: true, message: '请选择摄像头类型', trigger: 'change' }
                ],
                url: [
                    { required: true, message: '请输入摄像头链接', trigger: 'blur' }
                ]
            },
            paginatedData: [],
            taskTotalData: [],
            dialogVisible: false,
            currentTask: {},
            currentPage: 1,
            pageSize: 10,
            cameraList: [],
            currentCameraUrl: '',
            currentResult: {},
            resultList: [],//全部的任务结果
            newCameraUrl: '',//添加新摄像头时的url
            timeRange: [],//查看结果的时间范围
            paginatedResultList: [], // 存储结果图片数组
            currentImagePage: 1, // 当前分页
            imagesPerPage: 9, // 每页显示的图片数量
        }
    }
};
</script>

<style src="../style/Task.css" scoped></style>
