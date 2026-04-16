# Cloudflare Wrangler Proxy

A lightweight Go-based orchestrator and Docker environment for **Cloudflare Wrangler**.  
This tool isolates the entire Cloudflare development stack, ensuring your host system remains clean of Node.js, npm, or global package bloat.

## Why this exists?

Cloudflare's Wrangler CLI often introduces breaking changes or requires specific Node.js versions. By wrapping it in a Docker container and providing a Go-based CLI proxy, we achieve:

* **Zero-Footprint:** No Node.js or npm required on the host machine.
* **Version Pinning:** Guaranteed consistency across different development environments.
* **Automated Updates:** A simple `-U` flag updates both the binary and the underlying Docker image.
* **Custom Extensibility:** Ability to intercept commands and add custom logic before passing them to Wrangler.

## Prerequisites

* **Docker** (installed and running)
* **Go** (only if building from source)

## Installation & Basic Usage

```bash
# Show Golang binary version
cf -V

# Self-updater
cf -U
```

## Commands

```bash
# Show Wrangler version
cf version
