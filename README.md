# cert-checker
CertChecker 是一款基于 Go 的轻量级、易配置的开源工具，用于监控 HTTPS 网站的 SSL/TLS 证书有效期。它能定期检查指定的域名，并在证书即将过期时发送告警通知，帮助开发者、运维人员和管理员避免因证书过期导致的服务中断问题。

## 特性

- 单一 Go 二进制文件，启动快、资源占用低
- 基于 Debian Slim，体积小
- 内置 CA 证书，支持 HTTPS / TLS 校验
- 通过配置文件定义检测规则
- 适合服务器、容器、定时任务场景


## 镜像地址

Docker Hub：https://hub.docker.com/r/193002818/cert-checker

示例docker-compose.yaml如下:

```
version: "3"

services:
  cert-checker:
    image: 193002818/cert-checker:v1.0.12-2026-01-10T18_27_36
    container_name: cert-checker
    restart: unless-stopped
    volumes:
      - "./cert-checker-config.yaml:/app/cert-checker-config.yaml:ro" 
      - /etc/localtime:/etc/localtime
      - /usr/share/zoneinfo:/usr/share/zoneinfo
    environment:
      - TZ=Asia/Shanghai
      - LANG=C.UTF-8
      - LC_ALL=C.UTF-8
```
