# InView使用说明书

## 1 部署

### 1.1 前端部署

```Plain
npm run build 
```

将产生的dist文件放在盒子内某个文件夹下面

例如：/home/nvidia/Service，与下面nginx的配置文件里的root目录文件相对应

盒子内安装nginx，可以在以下路径中安装

```Bash
#进入安装目录
cd /usr/local

#下载nginx
curl -O http://nginx.org/download/nginx-1.22.1.tar.gz

#解压
tar -zxvf nginx-1.22.1.tar.gz

#进入nginx-1.22.1目录
cd nginx-1.22.1

#执行配置
./configure
make
make install

#执行完后，在/usr/local目录下生成了一个nginx
#启动nginx
cd /usr/local/nginx/sbin
 ./nginx
 
 #修改配置文件
 vim nginx.conf
 
 #监听8080端口，设置默认的主页文件位置，配置静态文件路径
server {
    listen 8088;
    listen 80；
    server_name localhost;

    location / {
        root /home/nvidia/Service/dist; 
        index index.html index.htm;

        location /api/ {
            proxy_pass http://localhost:8079;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_cache_bypass $http_upgrade;
        }
    }

    location /statistic/ {
        alias /home/nvidia/statistic/;
    }
}

设置开启自动启动
systemctl enable nginx.service

#重新加载配置文件
/usr/local/nginx/sbin/nginx -s reload
重启 Nginx
/usr/local/nginx/sbin/nginx -s reopen
启动 Nginx
/usr/local/nginx/sbin/nginx 
```

打开ip:8088检查前端页面是否能显示

### 1.2 后端部署

交叉编译后端文件

```Plain
CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build main.go
```

将可执行文件拖到某个文件夹下面，执行下面的命令可以运行后端

```Plain
sudo chmod 777 编译文件的名字
nohup ./编译文件的名字 &
```

环境文件.env如下：

```Plain
NODE_ENV=development
LOG_PATH=/var/log/insense-0902.log
MERN_FRONTEND=8079
JSON_PATH=/home/nvidia/statistic
IMAGE_PATH=/home/nvidia/statistic
DB_HOST=127.0.0.1
DB_PORT=27017
DB_NAME=insense
DB_USER=insense
DB_PASS=insense
SERVER_PORT=:8079
CONTEXT_TIMEOUT=10
```

### 1.3 mongo

由于盒子系统原因，只支持mongo 4.0的镜像，为了快速安装，可以直接load打包好的镜像

```Plain
#加载镜像，这一步可以替换成pull
docker load -i mongo-image.tar

#查看镜像
docker images

#启动
docker run -itd --name mongo -p 27017:27017 mongo:4.0 --auth

#管理员账号
db.createUser({ user:'admin',pwd:'admin',roles:[ { role:'userAdminAnyDatabase', db: 'admin'},"readWriteAnyDatabase"]});

#数据库账号，与env中相匹配
db.createUser({user: "insense", pwd: "insense", roles: [{ role: "readWrite", db: "insense" }]})

#建四个collection
db.createCollection("cameras")
db.createCollection("tasks")
db.createCollection("results")
db.createCollection("users")
```

### 1.4 Supervisor开启自启与进程守护

安装supervisor支持开机自动启动前后端，并且支持进程守护

参考网址：https://gcore.com/learning/how-to-set-up-supervisor-on-ubuntu/

```Plain
#安装
sudo apt install supervisor -y

#查看状态
sudo systemctl status supervisor

#配置
sudo vim /etc/supervisor/conf.d/myscript.conf

#配置文件例子
[program:start-services]
command=bash -c "/home/nvidia/Service/start_service.sh"
directory=/home/nvidia/Service
autostart=true
autorestart=true
stderr_logfile=/var/log/services.err.log
stdout_logfile=/var/log/services.out.log
user=nvidia
startretries=3
startsecs=5
restartpause=10

#更新配置
sudo supervisorctl reread
sudo supervisorctl update

#启动与暂停
sudo supervisorctl start start-services
sudo supervisorctl stop start-services
```

start-services.sh脚本需要同时执行启动mongo和运行后端两个操作，下面是个例子（文件名字和日志前后都要相对应）：

```Plain
/usr/bin/docker start mongo
nohup /home/nvidia/Service/insense-0902 &>> /var/log/insense-0902.log &
```



docker启动命令行示例

```bash
docker run -d --runtime nvidia --network host -v /home/nvidia/statistic/1/66d52758f9af25e756aa1cda/log:/log -v /home/nvidia/statistic/1/66d52758f9af25e756aa1cda/image:/image -v /home/nvidia/statistic/1/66d52758f9af25e756aa1cda/video:/video -v /home/nvidia/statistic/1/66d52758f9af25e756aa1cda/config:/config safety_firesmoke:v1.0 --json /config/safety_fireSmoke.json
```



## 2 使用说明

登陆前端网址：ip:8088

### 2.1 登录和注册

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=YTg0OGZlNjYwNDBjZmIwYTIxN2FmNDBhYTFiOWRhMjZfcW9TaVRjcWtGWHRrZzBNOFhRVlliZ1VYTG9xaWlzemhfVG9rZW46UW1CY2Juc2JKb3RFQWV4RTRyVGNoYVdpbm9jXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

登录前端界面，首先注册账号，接着使用注册的账号登录

### 2.2 总览

总览界面可以查看机器的详细情况和任务详情

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=YmVmNTU3M2U4YzI1MjE5MDRjZThmOTlkY2ExMWQwMzhfMXNuY1JkZXF2Wmx4OER2ZE96U3lCZ0FrNlJudlpkaTNfVG9rZW46RUQ4WWI2dGlob3BYNmt4VmxUOGNBUFlNbmZmXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjE3ZDliYWM0ZGY5Mjg2MDM0MDVhZDMyMTg0OWY3YTRfbGFOeE9LVm01TFB3b1VxV0dEQnZ4d2h5YnlUcmdLNWpfVG9rZW46QjBSZGJsWG1Qb0I3eWx4ejRWNmN4RVJNbnpUXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

### 2.3 系统日志

查看边缘盒子日志，当盒子出现问题时，可以导出日志，发给运维人员查看详情

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=MTI3NGExZTg0YzUzNTA4NTk5MTk3OTFjYjRkMjM3NGZfRUlGaElUTXJVdmh5MFllTHNWQjhLQlk0YWJPQzlqY1hfVG9rZW46UnpPcGI1VU9ob0xxSmx4dzVIVGNzeWRFbjRlXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

### 2.4 账户管理

可以修改用户名，密码，手机号，邮箱，注册

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=OGMzZDlhZWE1ZjI2NTY4MTUzZmUzNDBmOTA5NTdjMDdfcTZSTUZTS2dQekNjZERRRDZJczBNeE13b1lyRWRlcTZfVG9rZW46Qm5vY2I4UkdIb1FvZHR4a3g2cGNDYVNVblNoXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

### 2.5 配置任务

算法任务的全流程分为三步：新建摄像头→为摄像头添加新的任务→查看任务结果

#### 2.5.1 新建、修改、删除摄像头

 在任务管理界面，点击添加摄像头按钮，填写摄像头信息

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=ZTkwNTcxMzI3MzNmOTI1M2ZmYmYwYTZlODg4MTY1MDZfVUhUazBmb1pvSm84VFRFN1NWZkM4THhjemtpY1FXMkVfVG9rZW46Q2tVWWJUblNQb3BYTVd4b1loYmNuZVp2bnVmXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=MmEwZmQ2YTliMGIwM2YwNDA4NTgyZTlkMjUzM2RhZWNfeDd5QjZOdEJwOU9odkViWW1NWjhUejhBTzZHcElJeGNfVG9rZW46UVpwbmJjZ201b3IyMEx4eDFOWWM5aGZIbktoXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

下面三个按钮的功能分别是：修改摄像头信息，删除摄像头，为摄像头添加新任务

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=NDFiYTViMGMwN2MxZmE2ZmM3Y2RmNTQ2ZDgyNjAzM2NfdkhaWDU4WTZCSzBGaUJ2eE5HSndWT3g5aWVpVUhSMG1fVG9rZW46TWlxaWJrYzlIb1ZqMUV4RFRjZWNEWUtnblhLXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=NmU0ZjM3YTkwNzdmZmM2NTgzNTdkZGQ2NGMyOTZlNTFfeWxLT1dhU20xTmhBNGxISVlHbW1LaHN5QklhNk9tTXVfVG9rZW46WTZjeWJZeHdab0NDbG54YTI2WGNxMmxTbmVmXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

**注意事项：**

- 修改摄像头信息不支持修改摄像头链接，只可以修改摄像头名称和备注
- 当摄像头有任务时，不支持删除摄像头
- 每个摄像头同类型的算法只能配置一个，可支持配置多个不同类的算法

#### 2.5.2 添加新的算法任务

点击加号按钮，为该摄像头配置新的算法识别任务

**①截图一张图片，点击下一步**

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=ZTdiMGNlOWM1YmYyNDNhZDMxZThiYTcyYjZhZDFjYzZfaGJQYkNnTkpMTE9OdHljTXJQRTdSTnpEeks5eVdPeWNfVG9rZW46U3k4c2JVN1ZObzloRkh4NUxKa2N3bFBHbkFoXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

**②点击绘制区域，绘制算法识别范围，绘制错误可点击清除清空画布**

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=MThjNWJmMTBiZDBkMTY5OWQwNzlhNjYwYmY1ZThkMzJfYlBWaFo1WXFsckIwZTV5OTBpencxTFFnYVoxV2F1NUFfVG9rZW46RklRcmJXNXpIb0NPVjl4YWVyY2NSNGowbmtmXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

**③配置算法详细信息**

包括：算法类型，logrotate相关信息，数据定时清理信息，工作时间段，推送地址，和与报警有关的喇叭、短信、功能

设置好之后，点击完成配置按钮

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=ZjA0ZjM1NDVhMDk4YzViNjNjZTM5MWE3MjAwZmJjMzZfOGtYa2tmWHpmWmdRWDFrcjRjTXl0ZWY1Sjh2ZUtOdXVfVG9rZW46SWI4d2IwcHg4b2xEWkZ4TWtMU2NWUkUzbkVnXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

#### 2.5.3 查看算法结果

①点击相对应的任务查看算法结果

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=MDI3OWZlNWRhOTllNTc2ZmQwNmNmYzUzYjJlNDgxYTZfMk1BdHAzbmtpSkdhZkNNaENxaGsyRFdOb3ZUTXF6RkRfVG9rZW46VWo4UWJ5cHhub0ZBNm54Q0VvYmNMakd5bjZiXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

②根据时间筛选

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=ODg5ZGQzZDU2NDZhMWE1Y2ZlMThiNmIzY2UzNTJlM2FfM0RyY0NrTnBZSVZCb1ZkakxsNHg5VkdzeFJPQVZJbDRfVG9rZW46VDRmNWI5dU40bzhYaTZ4c25SYWNKOVZLbnZnXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

③点击图片查看识别详情

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=OTRlZTUzZmQzYjhkNmZkODM3MTU5MmMxNmE0MmI4NDJfc3lXOXNydW5mdGJQUVZ1YVdrZEpNWGhhdFpTa2NjV3lfVG9rZW46UlQzUmJQWFRNb29tdXZ4WVhyVmNNNmxobktnXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

#### 2.5.4 算法启停

- 算法会在用户配置的时间范围内自动启动和停止
- 如果用户配置时间范围内算法运行错误，用户可以点击任务右侧对应的【启动算法】按钮，如果算法还是无法运行，请联系运维人员
- 在算法运行时间内，如果用户想要暂停现在的算法，用户可以点击任务右侧对应的【停止算法】钮，如需在启动需要点击任务右边相对应的【启动算法】按钮

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=ZWVhNGZlZTVlMTJiZDBlMDU0NDBhMmRkNWJkZmZkZjRfSzhXTEVoNldCWXlHdFlzYjdMcmFpeGdZbXI4aDViaXlfVG9rZW46SWkwWmI1T1FwbzBYU0d4dDBlTWNEQmdHblZiXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

#### 2.5.5 修改算法配置

点击任务右侧相对应的【修改算法配置】按钮

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=OTE3MTBlODZjYzg5MDkwMWU3ZmIzZTc3MGQwNjIzZDNfTEtlcmlpRzQwNG1pTkZpVWY4ZDYwWkZrTjQyeG1RN3JfVG9rZW46VEw2Z2JKYkoxb3dqZ1p4N2VqMGNYaU8ybnFkXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

①截取图片

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=YmNmNTUwNmUxYjA0N2I4YWI0OTBmMGYxNWIwZjYxMTZfRkI5SUNUSkFIUGNOWEJ5NVpxaUEyZ3Fjb3QwcklBc0pfVG9rZW46WGRMT2JnbnhGb21vVHF4UEtiS2NOcDVmbnNoXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

②修改之前的相关配置

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=MTExYTZlOTNmYzY3NWRhZTFkZGU1NGYxYjk4YjBmZjZfeTRUOXp0a0xUWVkycG5zbkhTNXJrREdjVjJoamtXNVpfVG9rZW46UlhDdGJzbTVQb2JHVXZ4VlBNTGNJRXVobkRnXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=ODY4Mjk0ZTRmNzhkYzg1MjE3ZmMwNTYzM2E2Nzc3OGJfc2Z5TnpzRUdOblM5SWs5Sm1LOEdoSUllU2ZxZFFib25fVG9rZW46VXE1Z2Jia3FUb1h1OHN4anB1Y2NWTGtxbjdiXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

#### 2.5.6 删除算法任务

点击任务右侧对应的【删除任务】按钮即可删除

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=MDU4NTBkZDQwMzMwMWU0MTBkYmJiOGQ4ZDRmNzI2NmZfWVZPTm5UMklwWU11YlFvMjRteTEwWVI1OW9oOU1yRmlfVG9rZW46RzhZVmIxNlNtb2U0b054WXJkcWNpUVFrbmdnXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)

#### 2.5.7 筛选算法任务

当用户配置了较多的额摄像头和算法任务的时候可以使用筛选框筛选

![img](https://ok8tt02f6l.feishu.cn/space/api/box/stream/download/asynccode/?code=MDVmYmM1OTIzMDRkMGZmOGJhZGViZGFmYWY5OTAxNWVfeTIxTWsxbmtsQ1RabTd0UjEzQXZRaU9kZlZ0aXZ0MThfVG9rZW46UWVwQ2JMMEVTb0Zpa1l4THNVQmNlbktVblJoXzE3MjUyNTU1NzE6MTcyNTI1OTE3MV9WNA)
