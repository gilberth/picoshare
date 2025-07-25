# PicoShare

[![CircleCI](https://circleci.com/gh/mtlynch/picoshare.svg?style=svg)](https://circleci.com/gh/mtlynch/picoshare)
[![Docker Version](https://img.shields.io/docker/v/mtlynch/picoshare?sort=semver&maxAge=86400)](https://hub.docker.com/r/mtlynch/picoshare/)
[![Docker Pulls](https://img.shields.io/docker/pulls/mtlynch/picoshare.svg?maxAge=604800)](https://hub.docker.com/r/mtlynch/picoshare/)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/m/mtlynch/picoshare)](https://github.com/mtlynch/picoshare/commits/master)
[![GitHub last commit](https://img.shields.io/github/last-commit/mtlynch/picoshare)](https://github.com/mtlynch/picoshare/commits/master)
[![Contributors](https://img.shields.io/github/contributors/mtlynch/picoshare)](https://github.com/mtlynch/picoshare/graphs/contributors)
[![License](http://img.shields.io/:license-agpl-blue.svg?style=flat-square)](LICENSE)

## Overview

PicoShare is a minimalist service that allows you to share files easily.

- [Live demo](https://demo.pico.rocks)

[![PicoShare demo](https://raw.githubusercontent.com/mtlynch/picoshare/master/docs/readme-assets/demo.gif)](https://raw.githubusercontent.com/mtlynch/picoshare/master/docs/readme-assets/demo-full.gif)

## Why PicoShare?

There are a million services for sharing files, but none of them are quite like PicoShare. Here are PicoShare's advantages:

- **Direct download links**: PicoShare gives you a direct download link you can share with anyone. They can view or download the file with no ads or signups.
- **No file restrictions**: Unlike sites like imgur, Vimeo, or SoundCloud that only allow you to share specific types of files, PicoShare lets you share any file of any size.
- **No resizing/re-encoding**: If you upload media like images, video, or audio, PicoShare never forces you to wait on re-encoding. You get a direct download link as soon as you upload the file, and PicoShare never resizes or re-encodes your file.

## Run PicoShare

### From source

```bash
PS_SHARED_SECRET=somesecretpass PORT=4001 \
  go run cmd/picoshare/main.go
```

### From Docker

To run PicoShare within a Docker container, mount a volume from your local system to store the PicoShare sqlite database.

#### With Shared Secret Authentication (default)

```bash
docker run \
  --env "PORT=4001" \
  --env "PS_SHARED_SECRET=somesecretpass" \
  --publish 4001:4001/tcp \
  --volume "${PWD}/data:/data" \
  --name picoshare \
  mtlynch/picoshare
```

#### With Authentik Authentication

```bash
docker run \
  --env "PORT=4001" \
  --env "PS_AUTH_PROVIDER=authentik" \
  --env "PS_AUTHENTIK_URL=https://auth.example.com" \
  --env "PS_AUTHENTIK_CLIENT_ID=your-client-id" \
  --env "PS_AUTHENTIK_CLIENT_SECRET=your-client-secret" \
  --env "PS_AUTHENTIK_REDIRECT_URI=http://yourserver:4001/auth/callback" \
  --publish 4001:4001/tcp \
  --volume "${PWD}/data:/data" \
  --name picoshare \
  mtlynch/picoshare
```

#### Using environment file

You can also use an environment file (`.env`) to configure PicoShare:

```bash
docker run \
  --env-file .env \
  --publish 4001:4001/tcp \
  --volume "${PWD}/data:/data" \
  --name picoshare \
  mtlynch/picoshare
```

### From Docker + cloud data replication

If you specify settings for a [Litestream](https://litestream.io/)-compatible cloud storage location, PicoShare will automatically replicate your data.

You can kill the container and start it later, and PicoShare will restore your data from the cloud storage location and continue as if there was no interruption.

```bash
PORT=4001
PS_SHARED_SECRET="somesecretpass"
LITESTREAM_BUCKET=YOUR-LITESTREAM-BUCKET
LITESTREAM_ENDPOINT=YOUR-LITESTREAM-ENDPOINT
LITESTREAM_ACCESS_KEY_ID=YOUR-ACCESS-ID
LITESTREAM_SECRET_ACCESS_KEY=YOUR-SECRET-ACCESS-KEY

docker run \
  --publish "${PORT}:${PORT}/tcp" \
  --env "PORT=${PORT}" \
  --env "PS_SHARED_SECRET=${PS_SHARED_SECRET}" \
  --env "LITESTREAM_ACCESS_KEY_ID=${LITESTREAM_ACCESS_KEY_ID}" \
  --env "LITESTREAM_SECRET_ACCESS_KEY=${LITESTREAM_SECRET_ACCESS_KEY}" \
  --env "LITESTREAM_BUCKET=${LITESTREAM_BUCKET}" \
  --env "LITESTREAM_ENDPOINT=${LITESTREAM_ENDPOINT}" \
  --name picoshare \
  mtlynch/picoshare
```

Notes:

- Only run one Docker container for each Litestream location.
  - PicoShare can't sync writes across multiple instances.

### Using Docker Compose

To run PicoShare under docker-compose, copy the following to a file called `docker-compose.yml` and then run `docker-compose up`.

```yaml
version: "3.2"
services:
  picoshare:
    image: mtlynch/picoshare
    environment:
      - PORT=4001
      # For shared secret authentication (default)
      - PS_SHARED_SECRET=dummypass # Change to any password
      
      # For Authentik authentication (uncomment and configure)
      # - PS_AUTH_PROVIDER=authentik
      # - PS_AUTHENTIK_URL=https://auth.example.com
      # - PS_AUTHENTIK_CLIENT_ID=your-client-id
      # - PS_AUTHENTIK_CLIENT_SECRET=your-client-secret
      # - PS_AUTHENTIK_REDIRECT_URI=http://yourserver:4001/auth/callback
    ports:
      - 4001:4001
    command: -db /data/store.db
    volumes:
      - ./data:/data
```

## Parameters

### Command-line flags

| Flag  | Meaning                 | Default Value     |
| ----- | ----------------------- | ----------------- |
| `-db` | Path to SQLite database | `"data/store.db"` |

### Environment variables

| Environment Variable | Meaning                                                                              |
| -------------------- | ------------------------------------------------------------------------------------ |
| `PORT`               | TCP port on which to listen for HTTP connections (defaults to 4001).                 |
| `PS_BEHIND_PROXY`    | Set to `"true"` for better logging when PicoShare is running behind a reverse proxy. |
| `PS_AUTH_PROVIDER`   | Authentication provider: `shared_secret` or `authentik` (defaults to shared_secret). |
| `PS_SHARED_SECRET`   | (required for shared_secret auth) Specifies a passphrase for the admin user.         |

### Authentik Authentication Variables

When using `PS_AUTH_PROVIDER=authentik`, configure these variables:

| Environment Variable | Meaning                                                                              |
| -------------------- | ------------------------------------------------------------------------------------ |
| `PS_AUTHENTIK_URL`   | Base URL of your Authentik instance (e.g., `https://auth.example.com`).              |
| `PS_AUTHENTIK_CLIENT_ID` | OAuth2 Client ID from your Authentik application.                                |
| `PS_AUTHENTIK_CLIENT_SECRET` | OAuth2 Client Secret from your Authentik application.                        |
| `PS_AUTHENTIK_REDIRECT_URI` | OAuth2 redirect URI (e.g., `http://yourserver:4001/auth/callback`).           |
| `PS_AUTHENTIK_REVERSE_PROXY` | Set to `"true"` to use reverse proxy authentication instead of OAuth2.        |
| `PS_AUTHENTIK_TRUSTED_PROXIES` | Comma-separated list of trusted proxy IPs for reverse proxy mode.           |

### Docker environment variables

You can adjust behavior of the Docker container by specifying these Docker-specific variables with `docker run -e`:

| Environment Variable           | Meaning                                                                                               |
| ------------------------------ | ----------------------------------------------------------------------------------------------------- |
| `LITESTREAM_BUCKET`            | Litestream-compatible cloud storage bucket where Litestream should replicate data.                    |
| `LITESTREAM_ENDPOINT`          | Litestream-compatible cloud storage endpoint where Litestream should replicate data.                  |
| `LITESTREAM_ACCESS_KEY_ID`     | Litestream-compatible cloud storage access key ID to the bucket where you want to replicate data.     |
| `LITESTREAM_SECRET_ACCESS_KEY` | Litestream-compatible cloud storage secret access key to the bucket where you want to replicate data. |
| `LITESTREAM_RETENTION`         | The amount of time Litestream snapshots & WAL files will be kept (defaults to 72h).                   |

### Docker build args

If you rebuild the Docker image from source, you can adjust the build behavior with `docker build --build-arg`:

| Build Arg            | Meaning                                                                     | Default Value |
| -------------------- | --------------------------------------------------------------------------- | ------------- |
| `litestream_version` | Version of [Litestream](https://litestream.io/) to use for data replication | `0.3.9`       |

## PicoShare's scope and future

PicoShare is maintained by Michael Lynch as a hobby project.

Due to time limitations, I keep PicoShare's scope limited to only the features that fit into my workflows. That unfortunately means that I sometimes reject proposals or contributions for perfectly good features. It's nothing against those features, but I only have bandwidth to maintain features that I use.

## Docker Build Optimizations

This fork includes optimizations to improve Docker build performance:

### Build Context Optimization

The `.dockerignore` file has been enhanced to exclude unnecessary files from the Docker build context:

- Development artifacts (`node_modules/`, `.git/`, `.serena/`, `.claude/`)
- Binary files (`picoshare-binary`, `*.log`)
- Documentation and test files (`e2e/`, `docs/`)

This reduces the build context from ~98MB to ~22MB, significantly improving build times.

### Build Performance

- Build context size: ~22MB (optimized from ~98MB)
- Typical build time: 2-5 minutes (depending on network and system resources)
- Docker cache is utilized effectively for faster rebuilds

### Troubleshooting Build Issues

If Docker builds timeout or fail:

1. Clean Docker system: `docker system prune -f`
2. Retry the build command (builds may require multiple attempts)
3. Ensure adequate disk space and memory for Docker operations

## Deployment

PicoShare is easy to deploy to cloud hosting platforms:

- [fly.io](docs/deployment/fly.io.md)

## Tips and tricks

### Reclaiming reserved database space

Some users find it surprising that when they delete files from PicoShare, they don't gain back free space on their filesystem.

When you delete files, PicoShare reserves the space for future uploads. If you'd like to reduce PicoShare's usage of your filesystem, you can manually force PicoShare to give up the space by performing the following steps:

1. Shut down PicoShare.
1. Run `sqlite3 data/store.db 'VACUUM'` where `data/store.db` is the path to your PicoShare database.

You should find that the `data/store.db` should shrink in file size, as it relinquishes the space dedicated to previously deleted files. If you start PicoShare again, the System Information screen will show the smaller size of PicoShare files.
