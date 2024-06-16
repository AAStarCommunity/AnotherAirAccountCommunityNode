# 使用官方 Node.js 作为基础镜像
FROM node:20.14.0-alpine3.20

# 安装 pnpm
RUN npm install -g pnpm

# 设置工作目录
WORKDIR /usr/src/app

# 复制 package.json 和 pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# 安装项目依赖
RUN pnpm install

# 复制项目文件到工作目录
COPY . .

# 应用运行在哪个端口
EXPOSE 3000

# 运行应用
CMD ["pnpm", "dev"]