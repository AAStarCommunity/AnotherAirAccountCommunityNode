FROM node:20.14.0-alpine3.20

RUN npm install -g pnpm

COPY . .

RUN pnpm install \
    && pnpm build

EXPOSE 80

CMD ["pnpm", "start"]