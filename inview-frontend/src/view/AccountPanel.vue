<template>
    <div>
        <el-row style="display: flex; justify-content: flex-end; margin-bottom: 20px; font-weight: bolder;">
            <el-button class="text-button" icon="el-icon-switch-button" @click="logout()">退出登录</el-button>
        </el-row>
        <el-descriptions :column="1" border>
            <el-descriptions-item>
                <template v-slot:label>
                    用户名
                </template>
                <template v-if="!isEditing.username">
                    {{ accountInfo.username }}
                    <el-button size="small" style='float:right' icon="User" @click="editField('username')">
                        修改
                    </el-button>
                </template>
                <template v-else>
                    <el-input style="width: 400px;" v-model="newInfo.username" @keyup.enter="confirmEdit"
                        @blur="cancelEdit('username')" size="small" ref="usernameInput"></el-input>
                    <div class="confirm-edit-text">
                        请按回车确认修改
                    </div>
                </template>
            </el-descriptions-item>
            <el-descriptions-item>
                <template v-slot:label>
                    手机号
                </template>
                <template v-if="!isEditing.phone">
                    {{ accountInfo.phone }}
                    <el-button size="small" style='float:right' icon="Iphone" @click="editField('phone')">
                        修改
                    </el-button>
                </template>
                <template v-else>
                    <el-input style="width: 400px;" v-model="newInfo.phone" @keyup.enter="confirmEdit"
                        @blur="cancelEdit('phone')" size="small" ref="phoneInput"></el-input>
                    <div class="confirm-edit-text">
                        请按回车确认修改
                    </div>
                </template>
            </el-descriptions-item>
            <el-descriptions-item>
                <template v-slot:label>
                    邮箱
                </template>
                <template v-if="!isEditing.email">
                    {{ accountInfo.email }}
                    <el-button size="small" style='float:right' icon="Message"
                        @click="editField('email')">修改</el-button>
                </template>
                <template v-else>
                    <el-input style="width: 400px;" v-model="newInfo.email" @keyup.enter="confirmEdit"
                        @blur="cancelEdit('email')" size="small" ref="emailInput"></el-input>
                    <div class="confirm-edit-text">
                        请按回车确认修改
                    </div>
                </template>
            </el-descriptions-item>
        </el-descriptions>
        <el-row class="button-row">
            <el-button class="red-button" style="width: 45%;" @click="deleteDialogVisible = true">注销</el-button>
            <el-button class="border-button" style="width: 45%;" @click="showPasswordDialog()">修改密码</el-button>
        </el-row>

        <!-- 确认修改的对话框 -->
        <el-dialog title="确认" v-model="editDialogVisible" width="30%">
            <span>是否确认更改{{ editingFieldLabel }}为 {{ newInfo[editingField] }}？</span>
            <template v-slot:footer>
                <span class="dialog-footer">
                    <el-button class="text-button" @click="editDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="updateInfo">确认</el-button>
                </span>
            </template>
        </el-dialog>

        <!-- 确认注销的对话框 -->
        <el-dialog title="确认" v-model="deleteDialogVisible" width="30%">
            <span>确认注销用户？</span>
            <template v-slot:footer>
                <span class="dialog-footer">
                    <el-button class="text-button" @click="deleteDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="deleteAccount">确认</el-button>
                </span>
            </template>
        </el-dialog>

        <!-- 修改密码的对话框 -->
        <el-dialog title="修改密码" v-model="passwordDialogVisible" width="30%">
            <el-form :model="passwordForm" ref="passwordForm" :rules="passwordRules" label-width="100px">
                <el-form-item label="旧密码" prop="oldPassword">
                    <el-input type="password" show-password v-model="passwordForm.oldPassword"></el-input>
                </el-form-item>
                <el-form-item label="新密码" prop="newPassword">
                    <el-input type="password" show-password v-model="passwordForm.newPassword"></el-input>
                </el-form-item>
                <el-form-item label="确认密码" prop="confirmPassword">
                    <el-input type="password" show-password v-model="passwordForm.confirmPassword"></el-input>
                </el-form-item>
            </el-form>
            <template v-slot:footer>
                <span class="dialog-footer">
                    <el-button @click="passwordDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="confirmPasswordChange">确认</el-button>
                </span>
            </template>
        </el-dialog>

    </div>
</template>

<script>
import { getUserInfo, editInfo, editPassword, closeAccount } from '../js/request';

export default {
    name: 'AccountInfo',
    async mounted() {
        this.showInfo();
    },
    data() {
        return {
            accountInfo: {
                _id: "",
                username: "获取失败",
                phone: "获取失败",
                email: "获取失败",
            },
            isEditing: {
                username: false,
                phone: false,
                email: false,
            },
            newInfo: {
                username: '',
                phone: '',
                email: ''
            },
            editingField: '',
            editingFieldLabel: '',
            editDialogVisible: false,
            passwordDialogVisible: false,
            deleteDialogVisible: false,
            showPassword: false,
            passwordForm: {
                oldPassword: '',
                newPassword: '',
                confirmPassword: '',
            },
            passwordRules: {
                oldPassword: [
                    { required: true, message: '请输入旧密码', trigger: 'blur' }
                ],
                newPassword: [
                    { required: true, message: '请输入新密码', trigger: 'blur' },
                    { min: 6, message: '新密码长度不能少于 6 个字符', trigger: 'blur' }
                ],
                confirmPassword: [
                    { required: true, message: '请确认新密码', trigger: 'blur' },
                    { validator: this.validateConfirmPassword, trigger: 'blur' }
                ]
            },
            phoneRules: [
                { required: true, message: '请输入手机号', trigger: 'blur' },
                { pattern: /^[1]([3-9])[0-9]{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
            ],
            emailRules: [
                { required: true, message: '请输入邮箱地址', trigger: 'blur' },
                { pattern: /^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$/, message: '请输入正确的邮箱地址', trigger: 'blur' }
            ]
        };
    },
    methods: {
        async showInfo() {
            const res = await getUserInfo();
            if (res && res.data.status === 200) {
                this.accountInfo = { ...res.data.data };
                this.newInfo = this.accountInfo
            }
        },
        editField(field) {
            this.isEditing[field] = true;
            this.editingField = field;
            this.newInfo[field] = this.accountInfo[field];
            this.editingFieldLabel = field === 'username' ? '用户名' : field === 'phone' ? '手机号' : '邮箱';
            this.$nextTick(() => {
                this.$refs[`${field}Input`].focus();
            });
        },
        confirmEdit() {
            if (this.editingField === 'phone' && !this.phoneRules[1].pattern.test(this.newInfo.phone)) {
                this.$message.error('请输入正确的手机号');
                return;
            }
            if (this.editingField === 'email' && !this.emailRules[1].pattern.test(this.newInfo.email)) {
                this.$message.error('请输入正确的邮箱地址');
                return;
            }
            this.editDialogVisible = true;
            this.isEditing[this.editingField] = false;
        },
        cancelEdit(field) {
            this.isEditing[field] = false;
        },
        async updateInfo() {
            await editInfo({
                _id: this.accountInfo._id,
                ...this.newInfo,
            });
            this.showInfo()
            this.editDialogVisible = false;
        },
        togglePasswordVisibility() {
            this.showPassword = !this.showPassword;
        },
        showPasswordDialog() {
            this.passwordDialogVisible = true;
        },
        validateConfirmPassword(rule, value, callback) {
            if (value !== this.passwordForm.newPassword) {
                callback(new Error('两次输入的密码不一致'));
            } else {
                callback();
            }
        },
        async confirmPasswordChange() {
            const valid = await new Promise((resolve) => {
                this.$refs.passwordForm.validate((valid) => {
                    resolve(valid);
                });
            });
            if (valid) {
                this.passwordDialogVisible = false;
                const response = await editPassword({
                    _id: localStorage.getItem('user_id'),
                    ...this.passwordForm,
                })
                if (response && response.status === 200) {
                    localStorage.clear('user_id')
                    this.$router.push('/');
                }
            }
        },
        logout() {
            localStorage.clear('user_id')
            this.$router.push('/');
        },
        async deleteAccount() {
            await closeAccount()
            localStorage.clear('user_id')
            this.$router.push('/');
        }
    }
};
</script>

<style src="../style/Config.css" scoped></style>
