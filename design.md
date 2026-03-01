todoist-proxy is a service that proxies the todoist API v1 in order to filter specific projects or a hierarchy of projects.

The main endpint that handles the filtering is the sync API (https://developer.todoist.com/api/v1/#tag/Sync)

or https://api.todoist.com/api/v1/sync


All other endpoints are passed through to the todoist API.

This tool is written in golang and deployed as a cloud function on GCP. Alternatively it can be deployed to another service like fly.io, render, railway, or cloudflare workers, and may need to be written in TS to support one of these.

Clients of the proxy send their authentication credentials to the proxy as if it is the actual todoist API, and these are proxied through to the actual endpints.