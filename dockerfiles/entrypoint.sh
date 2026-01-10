#!/bin/bash
set -x
# 确保二进制具有可执行权限（挂载时可能丢失权限）
chmod +x /app/cert-checker-linux

# 打印一次启动信息
echo "[entrypoint] cert-checker 权限已确保，可执行"

# 启动程序
exec /app/cert-checker-linux -config=/app/cert-checker-config.yaml

