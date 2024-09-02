<template>
    <div>
        <el-row>
            <!-- <el-date-picker v-model="filterDate" type="daterange" range-separator="至" start-placeholder="开始日期"
                end-placeholder="结束日期" @change="filterLogs">
            </el-date-picker>
            <el-select v-model="filterLevel" placeholder="选择日志等级（默认显示全部）" @change="filterLogs" clearable>
                <el-option v-for="item in logLevels" :key="item.value" :label="item.label" :value="item.value">
                </el-option>
            </el-select> -->
            <el-button type="primary" @click="exportLogs">导出日志</el-button>

        </el-row>
        <el-row>
            <div class="log-container">
                <pre>{{ filteredLogs }}</pre>
            </div>
        </el-row>
    </div>
</template>


<script>
import { getSystemLog } from '../js/request';

export default {
    name: 'SystemLog',
    data() {
        return {
            totalLength: 100,
            logs: "",
            filterDate: [],
            filterLevel: "",
            logLevels: [{
                value: "INFO",
                label: "信息",
            },
            {
                value: "WARN",
                label: "警告",
            },
            {
                value: "ERROR",
                label: "错误",
            },
            {
                value: "DEBUG",
                label: "调试",
            }],
            filteredLogs: ""
        }
    },
    async mounted() {
        this.logs = await getSystemLog();
        this.filteredLogs = this.logs;
    },
    methods: {
        formatDate(date) {
            if (!date) return '';
            const d = new Date(date);
            const year = d.getFullYear();
            const month = (`0${d.getMonth() + 1}`).slice(-2);
            const day = (`0${d.getDate()}`).slice(-2);
            return `${year}-${month}-${day}`;
        },
        filterLogs() {
            let startDate = '';
            let endDate = '';
            if (this.filterDate) {
                startDate = this.formatDate(this.filterDate[0]);
                endDate = this.formatDate(this.filterDate[1]);
            }
            const level = this.filterLevel;
            const logsArray = this.logs.split('\r\n');

            this.filteredLogs = logsArray.filter(log => {
                const logDateMatch = log.match(/\d{4}-\d{2}-\d{2}/); // 只匹配到年月日
                const logLevelMatch = log.match(/\b(INFO|WARN|ERROR|DEBUG)\b/);
                
                const logDate = logDateMatch ? logDateMatch[0] : null;
                const logLevel = logLevelMatch ? logLevelMatch[0] : null;

                const isDateInRange = (!startDate || !endDate) || (logDate && logDate >= startDate && logDate <= endDate);
                const isLevelMatch = !level || logLevel === level;

                return (!startDate || !endDate || isDateInRange) && (!level || isLevelMatch);
            }).join('\r\n');
        },
        exportLogs() {
            let fileName = '系统日志';
            // if (this.filterDate) {
            //     const startDate = this.formatDate(this.filterDate[0]);
            //     const endDate = this.formatDate(this.filterDate[1]);
            //     fileName += `_${startDate}_至_${endDate}`;
            // } else {
            //     fileName += '_全部时间';
            // }
            // if (this.filterLevel) {
            //     fileName += `_${this.filterLevel}`;
            // } else {
            //     fileName += '_全部等级';
            // }

            fileName += '.log';
            const blob = new Blob([this.filteredLogs], { type: 'text/plain;charset=utf-8' });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = fileName;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        }
    }
}
</script>

<style src="../style/Log.css"></style>