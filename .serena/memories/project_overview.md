# PicoShare - Project Overview

## Purpose
PicoShare is a minimalist file sharing service written in Go that allows users to upload and share files with direct download links. Key features:
- Direct download links with no ads or signups
- No file restrictions (any file type, any size)
- No resizing/re-encoding of media files
- Simple web interface with dark/light theme support

## Tech Stack
- **Backend**: Go 1.24 with Gorilla Mux router
- **Database**: SQLite with optional Litestream replication
- **Frontend**: Vanilla JavaScript (ES modules), HTML templates, Bulma CSS framework
- **Icons**: FontAwesome 6
- **Image Processing**: Go image libraries for thumbnail generation
- **Build**: Docker containerization
- **Development**: modd for hot reload, development scripts in bash

## Key Architecture
- `/cmd/picoshare/main.go` - Application entry point
- `/handlers/` - HTTP handlers, templates, static assets
- `/store/` - Database layer
- `/picoshare/` - Core domain types
- `/space/` - File storage utilities
- `/garbagecollect/` - Cleanup utilities

## Special Features
- **Thumbnails**: System generates thumbnails for images (PNG, JPG, WEBP)
- **File Icons**: Custom icon system for different file types
- **Authentication**: Supports shared secret or Authentik OAuth2
- **Guest Links**: Allows temporary upload access
- **Themes**: Light/dark mode with CSS variables