## Multi-process Dockerfile: build Vite app, serve static on 3002, Nginx reverse proxy on 3000

# 1) Build stage: install deps, load production env, build with Vite
FROM node:20-alpine AS builder
WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

# Copy source and production env, then build
COPY . .
COPY deployments/production/.env .env
ENV NODE_ENV=production
RUN npm run build

# 2) Runtime stage: Nginx + static server
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

# Install nginx and a lightweight static file server (serve)
RUN apk add --no-cache nginx && \
	npm i -g serve

# Copy runtime files
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./package.json

# Configure Nginx by copying our static config into http.d (Alpine default)
RUN mkdir -p /etc/nginx/http.d && rm -f /etc/nginx/http.d/default.conf /etc/nginx/conf.d/default.conf
COPY nginx/default.conf /etc/nginx/http.d/default.conf
RUN mkdir -p /run/nginx /var/cache/nginx

# Add entrypoint to launch all processes
COPY docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 3000

CMD ["/entrypoint.sh"]
