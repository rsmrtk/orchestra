# syntax=docker/dockerfile:1.4
# Production-grade Node.js Dockerfile (multi-stage, non-root, tini)

FROM node:20-bookworm-slim AS base
WORKDIR /app
ENV NODE_ENV=production

FROM base AS deps
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --omit=dev

FROM base AS build
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci
COPY . .
# If you have a build step (TS/webpack), keep it; otherwise remove this line.
RUN npm run build

FROM base AS runtime
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y --no-install-recommends tini ca-certificates && \
    rm -rf /var/lib/apt/lists/*

USER node

COPY --from=deps --chown=node:node /app/node_modules /app/node_modules
COPY --from=build --chown=node:node /app/dist /app/dist
# If your app runs from source, replace the line above with:
# COPY --from=build --chown=node:node /app /app

EXPOSE 3000
ENV PORT=3000

ENTRYPOINT ["tini", "--"]
CMD ["node", "dist/index.js"]
