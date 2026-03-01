# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`todoist-proxy` is a Go service that proxies the [Todoist Sync API](https://developer.todoist.com/api/v1/#tag/Sync) (`https://api.todoist.com/api/v1/sync`) with project-level filtering. Clients authenticate with their own Todoist credentials, which are forwarded transparently. All non-Sync endpoints are passed through unmodified.

The primary deployment target is a GCP Cloud Function.

## Architecture

- **Sync endpoint** (`/api/v1/sync`): Intercepts requests, proxies to Todoist, then filters the response to only include allowed projects and their hierarchy.
- **All other endpoints**: Transparent reverse proxy to `https://api.todoist.com`.
- **Auth**: Passed through from client to Todoist as-is (no credential storage in the proxy).
