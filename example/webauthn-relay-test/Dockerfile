FROM node:20.14.0-alpine3.20

RUN npm install -g pnpm

WORKDIR /usr/src/app

COPY package.json pnpm-lock.yaml ./

RUN pnpm install
RUN pnpm build

COPY . ./

EXPOSE 80

CMD ["pnpm", "start"]