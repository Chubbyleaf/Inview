<template>
    <div class="container">
        <el-container>
            <el-aside style="width: 180px;">
                <img src="../assets/logo/SMEDI.png" style="height: 100px;" />
                <el-menu :default-active="menuName" class="el-menu-vertical-demo" @select="handleClick">
                    <el-menu-item index="OverviewPanel">
                        <i class="el-icon-menu"></i>
                        <span>总览</span>
                    </el-menu-item>
                    <el-menu-item index="TaskManagement">
                        <i class="el-icon-document"></i>
                        <span>任务管理</span>
                    </el-menu-item>
                    <el-menu-item index="SystemLog">
                        <i class="el-icon-document"></i>
                        <span>系统日志</span>
                    </el-menu-item>
                    <el-menu-item index="SystemConfig">
                        <i class="el-icon-setting"></i>
                        <span>系统设置</span>
                    </el-menu-item>
                    <el-menu-item index="AccountPanel">
                        <i class="el-icon-s-custom"></i>
                        <span>账户管理</span>
                    </el-menu-item>
                </el-menu>
            </el-aside>
            <el-container>
                <el-header>
                    <span>{{ headerTitle }}</span>
                </el-header>
                <el-main>
                    <router-view :name="menuName"></router-view>
                </el-main>
            </el-container>
        </el-container>
    </div>
</template>

<script>
export default {
    name: 'DashboardPanel',
    methods: {
        //切换菜单
        handleClick(index) {
            this.menuName = index;
            this.headerTitle = this.getMenuTitle(index);
            this.$router.push({ name: index });
        },
        //登出
        logout() {
            this.$router.push({ name: 'Login' });
        },

        getMenuTitle(index) {
            const titles = {
                'OverviewPanel': '总览',
                'TaskManagement': '任务管理',
                'SystemLog': '系统日志',
                'SystemConfig': '系统设置',
                'AccountPanel': '账户管理',

            };
            return titles[index];
        },


    },
    mounted() {
        this.$router.push({ name: this.menuName });
        // this.$router.push({ name: this.$route.name });

    },
    data() {
        return {
            headerTitle: this.$route.name? this.getMenuTitle(this.$route.name):"总览",
            menuName: this.$route.name ? this.$route.name : 'OverviewPanel'
        };
    },
};
</script>

<style scoped></style>