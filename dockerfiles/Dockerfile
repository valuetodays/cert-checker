FROM debian:12.6

# 设置非交互模式，避免卡住
ENV DEBIAN_FRONTEND=noninteractive

# 替换 Debian 12 APT 源为阿里云
RUN SOURCE_FILE="/etc/apt/sources.list.d/debian.sources" && \
    echo "Backup original sources..." && \
    cp "$SOURCE_FILE" "${SOURCE_FILE}.bak"

#RUN echo 'Writing Aliyun mirror sources...'
RUN cat > /etc/apt/sources.list.d/debian.sources <<'EOF'
Types: deb
URIs: http://mirrors.aliyun.com/debian
Suites: bookworm bookworm-updates
Components: main contrib non-free non-free-firmware
Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg

Types: deb
URIs: http://mirrors.aliyun.com/debian-security
Suites: bookworm-security
Components: main
Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg
EOF

# 安装CA证书、其它软件
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates dos2unix && \
    update-ca-certificates && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY ./cert-checker-linux /app/cert-checker-linux
COPY ./dockerfiles/cert-checker-config.yaml /app/
COPY ./dockerfiles/entrypoint.sh /app/entrypoint.sh

# 转换成 Unix 格式
RUN dos2unix /app/entrypoint.sh

RUN chmod +x /app/cert-checker-linux /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]
