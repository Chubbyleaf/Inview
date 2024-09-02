<template>
    <div class="container common-layout">
        <el-container>
            <el-main>
                <el-row style="width: 100%;height: 90vh;">
                    <el-col :span="8" class="login-container">
                        <img id="login-logo" src="../assets/logo/SMEDI.png" />
                    </el-col>
                    <el-col :span="16">
                        <div class="login-container">
                            <div class="login-box">
                                <div style="font-weight: bolder; font-size: larger; margin-bottom: 30px;">{{
                        isRegistering ?
                            '注册' : '登录' }}</div>
                                <el-form :model="form" label-width="80px">
                                    <el-form-item label="用户名">
                                        <el-input v-model="form.username"></el-input>
                                    </el-form-item>
                                    <el-form-item label="密码">
                                        <el-input type="password" v-model="form.password" show-password></el-input>
                                    </el-form-item>
                                    <el-form-item v-if="isRegistering" label="确认密码">
                                        <el-input type="password" v-model="form.confirmPassword"
                                            show-password></el-input>
                                    </el-form-item>
                                </el-form>
                                <el-row>
                                    <el-button @click="isRegistering ? handleRegister() : handleLogin()"
                                        id="login-btn1">
                                        {{ isRegistering ? '注册' : '登录' }}
                                    </el-button>
                                    <el-button class='border-button' id="login-btn2" @click="toggleRegistering">{{
                        isRegistering ? '已有账号？登录' : '没有账号，去注册' }}</el-button>
                                </el-row>
                            </div>
                        </div>
                    </el-col>
                </el-row>
            </el-main>
            <el-footer>Powered by Invix</el-footer>
        </el-container>
    </div>
</template>

<script>
import { login, register } from '../js/request';

export default {
    name: 'LoginPanel',
    data() {
        return {
            isRegistering: false,
            form: {
                username: '',
                password: '',
                confirmPassword: '', // 注册时使用
            },
        };
    },
    methods: {
        async handleLogin() {

            const { username, password } = this.form;
            const response = await login({ username, password });
            if (response && response.data.status === 200) {
                this.$router.push({ name: 'OverviewPanel' });
            }

        },
        async handleRegister() {
            if (this.form.password !== this.form.confirmPassword) {
                this.$message.error('两次输入的密码不一致');
                return;
            }
            const { username, password } = this.form;
            const response = await register({ username, password });
            if (response && response.data.status === 200) {
                this.isRegistering = false;
            }

        },
        toggleRegistering() {
            this.isRegistering = !this.isRegistering;
            this.form.username = '';
            this.form.password = '';
            this.form.confirmPassword = '';
        },
    },
}
</script>

<style src="../style/Login.css" scoped></style>
