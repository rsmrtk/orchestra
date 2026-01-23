# syntax=docker/dockerfile:1.4
# Production-grade Python Dockerfile (venv, non-root, slim)

FROM python:3.12-slim AS base
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1
WORKDIR /app

FROM base AS builder
RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y --no-install-recommends build-essential ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY requirements.txt ./
RUN --mount=type=cache,target=/root/.cache/pip \
    python -m venv /opt/venv && \
    /opt/venv/bin/pip install -r requirements.txt

FROM base AS runtime
RUN adduser --disabled-password --gecos "" --home /app appuser
ENV PATH="/opt/venv/bin:$PATH"

COPY --from=builder /opt/venv /opt/venv
COPY . .
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8000
ENV PORT=8000

# Replace with your actual entrypoint (gunicorn/uvicorn/flask/etc.)
CMD ["python", "-m", "app"]
