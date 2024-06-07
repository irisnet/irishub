FROM node:alpine AS builder
WORKDIR /app
COPY docs/package.json docs/yarn.lock ./
RUN yarn install
COPY docs/ .
RUN yarn run build

FROM nginx:alpine
RUN rm -rf /usr/share/nginx/html/*
COPY --from=builder /app/.vuepress/dist /usr/share/nginx/html/docs
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
