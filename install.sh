#!/bin/bash

# /bin/bash

# 创建目录
mkdir -p /root/data/osapi
mkdir -p /root/data/osapi/log
mkdir -p /root/data/osapi/config

# 拉取主程序
VERSION=$(curl -s https://raw.githubusercontent.com/WJQSERVER/commandapi/main/VERSION)
wget -O /root/data/osapi/VERSION https://raw.githubusercontent.com/WJQSERVER/commandapi/main/VERSION
wget -O /root/data/osapi/osapi https://github.com/WJQSERVER/commandapi/releases/download/${VERSION}/osapi-linux-amd64
chmod +x /root/data/osapi/osapi

# 配置文件
if [ ! -f /root/data/commandapi/config/config.toml ]; then
    wget -O /root/data/commandapi/config/config.toml https://raw.githubusercontent.com/WJQSERVER/commandapi/main/config/config.toml
fi

# 拉取 systemd unit 文件
wget -O /etc/systemd/system/osapi.service https://raw.githubusercontent.com/WJQSERVER/commandapi/main/osapi.service

# 启动服务
systemctl daemon-reload
systemctl enable osapi.service
systemctl start osapi.service
systemctl restart osapi.service
