# ---- 构建阶段 ----
FROM node:22-alpine AS builder

RUN corepack enable && corepack prepare pnpm@latest --activate

# 使用国内镜像加速
RUN pnpm config set registry https://registry.npmmirror.com

WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/web-antd/package.json apps/web-antd/package.json

# 复制内部包的 package.json（monorepo 依赖）
COPY packages/ packages/

RUN pnpm install --frozen-lockfile 2>/dev/null || pnpm install

COPY . .

RUN pnpm run build:antd

# ---- 运行阶段 ----
FROM nginx:1.27-alpine

COPY --from=builder /app/apps/web-antd/dist /usr/share/nginx/html

# Nginx 配置：SPA 路由 + API 反向代理
RUN cat > /etc/nginx/conf.d/default.conf << 'EOF'
server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;

    # SPA 路由回退
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 反向代理到后端
    location /api/ {
        proxy_pass http://mb-backend:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
EOF

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
