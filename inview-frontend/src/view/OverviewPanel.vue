<template>
    <div>
        <el-row class="overview-title">
            <div>资源管理</div>
        </el-row>

        <el-row :gutter="24" style="min-width: 1100px">
            <el-col :span="12">
                <div class="overview-box1" style="height: 430px;overflow-y: auto;">
                    <div class="overview-subtitle">主机信息</div>
                    <div v-if="systemInfo.hosts">
                        <el-row class="info-row">
                            <el-col :span="12">
                                <div class="info-data">{{ this.truncatedString(systemInfo.hosts.os) }}</div>
                                <div class="info-title">操作系统</div>
                            </el-col>
                            <el-col :span="12">
                                <div class="info-data">{{ this.truncatedString(systemInfo.hosts.procs) }}</div>
                                <div class="info-title">进程数</div>
                            </el-col>

                        </el-row>
                        <el-row class="info-row">
                            <el-col :span="12">
                                <div class="info-data">{{ this.truncatedString(systemInfo.hosts.platform) }}</div>
                                <div class="info-title">平台</div>
                            </el-col>
                            <el-col :span="12">
                                <div class="info-data">{{ this.truncatedString(systemInfo.hosts.platformVersion) }}
                                </div>
                                <div class="info-title">平台版本</div>
                            </el-col>
                        </el-row>
                        <el-row class="info-row">
                            <el-col :span="12">
                                <div class="info-data">{{ this.truncatedString(systemInfo.hosts.kernelArch) }}</div>
                                <div class="info-title">内核</div>
                            </el-col>
                            <el-col :span="12">
                                <div class="info-data">{{ this.truncatedString(systemInfo.hosts.kernelVersion) }}</div>
                                <div class="info-title">内核版本</div>
                            </el-col>
                        </el-row>
                        <el-row class="info-row">
                            <el-col :span="12">
                                <div class="info-data">{{ systemInfo.hosts.bootTime }}</div>
                                <div class="info-title">启动时间</div>
                            </el-col>
                            <el-col :span="12">
                                <div class="info-data">{{ formatUptime(elapsedUptime) }}</div>
                                <div class="info-title">运行时间</div>
                            </el-col>
                        </el-row>
                    </div>

                </div>
            </el-col>
            <el-col :span="12">
                <el-row>
                    <div class="overview-box1" style="height: 200px;">
                        <div class="overview-subtitle">内存信息</div>
                        <el-row v-if="systemInfo.mems">
                            <el-col :span="12">
                                <div id="mem-pieChart-container" ref='memoryChart'></div>
                            </el-col>
                            <el-col :span="12">
                                <el-descriptions :column="1" border style="width:100%;margin-top: 25px">
                                    <el-descriptions-item label="总内存">{{ `${systemInfo.mems.total}GB`
                                        }}</el-descriptions-item>
                                    <el-descriptions-item label="可用内存">{{ `${systemInfo.mems.available}GB`
                                        }}</el-descriptions-item>
                                    <el-descriptions-item label="已使用内存">
                                        {{ `${systemInfo.mems.used}GB` }}
                                    </el-descriptions-item>
                                </el-descriptions>
                            </el-col>

                        </el-row>
                    </div>
                </el-row>
                <el-row>
                    <div class="overview-box1" style="height: 200px;overflow-y: auto;margin-top: 15px;">
                        <div class="overview-subtitle">磁盘信息</div>
                        <el-row v-for="disk in systemInfo.disks" :key="disk.disk">
                            <el-descriptions :column="2" border style="width:100%;margin-top: 25px">
                                <el-descriptions-item label="磁盘名称">{{ `${disk.disk}`
                                    }}</el-descriptions-item>
                                <el-descriptions-item label="总空间">{{ `${disk.total}GB`
                                    }}</el-descriptions-item>
                                <el-descriptions-item label="空闲空间">{{ `${disk.free}GB`
                                    }}</el-descriptions-item>
                                <el-descriptions-item label="已用空间">{{ `${disk.used}GB`
                                    }}</el-descriptions-item>
                            </el-descriptions>
                        </el-row>
                    </div>
                </el-row>

            </el-col>

        </el-row>
        <el-row class="overview-title">
            <div>任务总览</div>
        </el-row>
        <el-row :gutter="24" style="min-width: 1000px">
            <el-col :span="12">
                <div class="overview-box2">
                    <div ref='cameraChart' class="chart-container"></div>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="overview-box2">
                    <div ref='taskChart' class="chart-container"></div>
                </div>
            </el-col>
        </el-row>
    </div>
</template>

<script>
import * as echarts from 'echarts';
import { getCameraList, getTaskList, getSystemInfo } from '../js/request';

export default {
    name: 'OverviewPanel',
    data() {
        return {
            cameraData: [],
            taskList: [],
            taskData: [],
            systemInfo: {},
            modelCategories: ['safety'],
            algorithmCategories: ['equipment', 'concreteSupport', 'areaEdge', 'fireSmoke'],
            elapsedUptime:"0"
        };
    },
    async mounted() {
        this.cameraData = await getCameraList()
        this.taskList = await getTaskList()
        this.systemInfo = await getSystemInfo()
        this.elapsedUptime = this.systemInfo.hosts.uptime;
        this.getTaskOverviewData();
        this.$nextTick(() => {
            this.drawCameraChart();
            this.drawTaskChart();
            this.drawMemoryChart();
        });
        this.updateUptime();

    },
    methods: {
        truncatedString(s) {
            if (s) {
                return s.length > 20 ? s.slice(0, 20) + '...' : s;

            } else {
                return "暂无数据"
            }
        },
        formatUptime(seconds) {
            const days = Math.floor(seconds / 86400);
            const hours = Math.floor((seconds % 86400) / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            const secs = seconds % 60;
            return `${days}天 ${hours}小时 ${minutes}分 ${secs}秒`;
        },
        updateUptime() {
            setInterval(() => {
                this.elapsedUptime += 1;
            }, 1000);
        },
        getTaskOverviewData() {
            const result = this.modelCategories.reduce((acc, algo) => {
                acc[algo] = this.algorithmCategories.map(() => 0);
                return acc;
            }, {});

            const algorithmToIndex = this.algorithmCategories.reduce((acc, model, idx) => {
                acc[model] = idx;
                return acc;
            }, {});

            this.taskList.forEach(task => {
                const model = task.model;
                const algorithm = task.algorithmType;
                if (model != '' && algorithm in algorithmToIndex && model in result) {
                    const index = algorithmToIndex[algorithm];
                    result[model][index] += 1;

                }
            });
            this.taskData = result

        },
        drawCameraChart() {
            const cameraChart = this.$refs.cameraChart;
            if (cameraChart) {
                const myChart = echarts.init(cameraChart);
                const option = this.getCameraChartOption();
                myChart.setOption(option);
                window.addEventListener('resize', myChart.resize);
            }

        },
        drawTaskChart() {
            const taskChart = this.$refs.taskChart;
            if (taskChart) {
                const myChart = echarts.init(taskChart);
                const option = this.getTaskChartOption();
                myChart.setOption(option);
                window.addEventListener('resize', myChart.resize);
            }
        },
        drawMemoryChart() {
            const memoryChart = this.$refs.memoryChart;
            if (memoryChart) {

                const myChart = echarts.init(memoryChart);
                const option = this.getMemoryChartOption();
                myChart.setOption(option);
                window.addEventListener('resize', myChart.resize);
            }
        },
        getCameraChartOption() {
            const data = this.cameraData;
            const typeCount = data.reduce((acc, camera) => {
                acc[camera.type] = (acc[camera.type] || 0) + 1;
                return acc;
            }, {});

            return {
                title: {
                    text: '摄像机统计',
                    left: 'center',
                    top: 'top',
                    textStyle: {
                        fontSize: 14,
                    }
                },
                tooltip: {},
                xAxis: {
                    type: 'value',
                    boundaryGap: [0, 0.01]
                },
                yAxis: {
                    type: 'category',
                    data: Object.keys(typeCount)
                },
                series: [
                    {
                        type: 'bar',
                        data: Object.values(typeCount),
                        itemStyle: {
                            color: '#939aec' // 更改柱状图颜色
                        }
                    }
                ]
            };
        },
        getTaskChartOption() {
            return {
                title: {
                    text: '任务统计',
                    left: 'center',
                    top: 'top',
                    textStyle: {
                        fontSize: 14,
                    }
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'shadow'
                    }
                },
                legend: {
                    data: ['安全'],
                    left: 'right',
                    top: 'top'
                },
                toolbox: {
                    show: true,
                    orient: 'vertical',
                    left: 'right',
                    top: 'center',
                    feature: {
                        mark: { show: true },
                        dataView: { show: false, readOnly: false },
                        magicType: { show: true, type: ['stack'] },
                        restore: { show: true },
                        saveAsImage: { show: false }
                    }
                },
                xAxis: [
                    {
                        type: 'category',
                        axisTick: { show: false },
                        data: ['安全帽/服', '砼支撑', '临边', '围栏', '明火烟雾']
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        minInterval: 1
                    }
                ],
                series: [
                    {
                        name: '安全',
                        type: 'bar',
                        barGap: 0,
                        emphasis: {
                            focus: 'series'
                        },
                        data: this.taskData.safety,
                        itemStyle: {
                            color: '#939aec' // 更改柱状图颜色
                        }

                    },
                ]
            };
        },
        getMemoryChartOption() {
            return {
                tooltip: {
                    trigger: 'item',
                    formatter: '{b}: {c}GB ({d}%)' // Format the tooltip to show the value with GB
                },
                legend: {
                    right: '10%',
                    orient: 'vertical',
                    top: 'center'
                },
                series: [
                    {
                        name: 'Access From',
                        type: 'pie',
                        radius: ['50%', '100%'],
                        avoidLabelOverlap: false,
                        itemStyle: {
                            borderRadius: 4,
                            borderColor: '#fff',
                            borderWidth: 2
                        },
                        center: ['30%', '50%'],
                        label: {
                            show: false,
                            position: 'center'
                        },
                        emphasis: {
                            label: {
                                show: true,
                                fontSize: 10,
                                fontWeight: 'bold'
                            }
                        },
                        labelLine: {
                            show: false
                        },
                        data: [
                            { value: this.systemInfo.mems.used, name: '已使用内存', itemStyle: { color: '#8287d6' } },
                            { value: this.systemInfo.mems.available, name: '可用内存', itemStyle: { color: '#c0c4ed' } }
                        ]
                    }
                ]
            };
        }
    }
};
</script>

<style src="../style/Overview.css"></style>