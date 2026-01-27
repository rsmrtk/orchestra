# syntax=docker/dockerfile:1.4
# Production-grade Node.js Dockerfile
# Purpose: install deps, build (if needed), and run with a minimal runtime image.
# How it works: deps stage -> build stage -> runtime stage (small + non-root).

# ===== Base image used by all stages =====
FROM node:20-bookworm-slim AS base
WORKDIR /app
ENV NODE_ENV=production

# ===== Stage 1: Dependencies (prod only) =====
FROM base AS deps

# Install minimal system certs for HTTPS.
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy lockfiles to get reproducible installs.
COPY package.json package-lock.json ./

# Install production dependencies only.
RUN --mount=type=cache,target=/root/.npm \
    npm ci --omit=dev

# ===== Stage 2: Build (optional) =====
FROM base AS build

# Install full deps for building (TS/webpack/etc.).
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci

# Copy source and build artifacts (if applicable).
COPY . .
# If your app has no build step, remove this line and use source in runtime.
RUN npm run build

# ===== Stage 3: Runtime =====
FROM base AS runtime

# Install tini to forward signals, and certs for HTTPS.
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y --no-install-recommends tini ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Use the pre-existing non-root 'node' user.
USER node

# Copy only what is needed to run.
COPY --from=deps --chown=node:node /app/node_modules /app/node_modules
COPY --from=build --chown=node:node /app/dist /app/dist
# If running from source, replace the line above with:
# COPY --from=build --chown=node:node /app /app

# Document the port. Use your real app port here.
EXPOSE 3000
ENV PORT=3000

# Tini handles SIGTERM/SIGINT properly in containers.
ENTRYPOINT ["tini", "--"]

# Default start command. Update to your app entrypoint.
CMD ["node", "dist/index.js"]
