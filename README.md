# 🧠 WhatsNew Service – Org-Wide GitHub Activity Tracker

---

## 📊 High Level Design

```mermaid
graph TD
  A[GitHub Org Repos] -->|JWT Auth| B(WhatsNew Service)
  B --> C{Background Commit Cache}
  B --> D[Main API Server :10020]
  B --> E[Metrics Server :9200]
  D --> F[Swagger UI (/swagger/index.html)]
  E -->|Prometheus scrape| G[Grafana Alloy]
```

---

## 🚀 CX Pipeline Status

[![CI/CD Status](https://github.com/dasmlab/whatsnew-service/actions/workflows/main.yaml/badge.svg)](https://github.com/dasmlab/whatsnew-service/actions/workflows/main.yaml)
> _Isolated Docker Networks, Secure Builds, and FluxCD GitOps Integration_

---

## 🚀 Features

- 🔐 **GitHub App OAuth2 Auth**
  - JWT-based token exchange using a GitHub App created at the organization level
  - Access is limited to selected repositories (as per app permissions)
  - No personal access tokens (PATs) required
  - Org admin retains full control over scope and revocation
  - Currently running under GitHub Free Org tier

- ⚙️ **Fully Dockerized + CI/CD Ready**
  - Clean multi-stage Docker build
  - Integrated with GitHub Actions for automated build, test, push
  - Generates FluxCD-ready deployment manifests in `k8s_envelope/`

- 📈 **Out-of-band Prometheus Metrics**
  - Metrics served via second Gin server on port `9200`
  - Uses `ginprom` for Prometheus-compatible scrape format
  - Clean separation of business logic (API) and telemetry (metrics)

- 🔄 **Dynamic Repo Discovery**
  - GitHub App JWT token automatically pulls all repos accessible to the installation
  - Top 2 commits per repo fetched and cached periodically

- 🧰 **Tech Stack**
  - Go 1.21+
  - Gin Web Framework
  - Logrus for structured logs
  - Swaggo for Swagger auto-gen docs
  - Depado's `ginprom` middleware for metrics

---

## 📂 Local Development

### 🔧 Requirements

- Go 1.21 or newer
- Docker
- A GitHub App configured with:
  - Read-only access to your organization repos
  - A downloaded `.pem` file (for signing JWT)
  - The `APP_ID` and `INSTALLATION_ID` for the GitHub App

### 🔧 Required Environment Variables

| Variable           | Description                                       |
|--------------------|---------------------------------------------------|
| `APP_ID`           | GitHub App ID                                     |
| `INSTALLATION_ID`  | App Installation ID                               |
| `PEMFILE`          | Absolute path to downloaded GitHub App .pem file  |

---

### 🛠️ Build Locally

```bash
./buildme.sh
```

### ▶️ Run Locally (Dockerized)

```bash
./runme_local.sh
```

---

## 📦 CI/CD & GitOps Workflow

- GitHub Actions-based isolated runner pipeline
- Builds the image, runs local healthchecks
- Pushes image to GHCR under `ghcr.io/lmcdasm/whatsnew-service:<version>`
- Substitutes version into `k8s_envelope/` template and commits into FluxCD repo

---

## 📊 Metrics

- **Scrape endpoint:** `http://localhost:9200/metrics`
- **Scrape target:** `Grafana Alloy` or any standard Prometheus-compatible collector
- **Middleware:** `ginprom`

---

## 🔍 Swagger API Docs

- Access interactive API docs at:

```
http://localhost:10020/swagger/index.html
```

- Fully interactive with test-curl support inline

---

## 🔮 Example Usage (cURL)

```bash
curl http://localhost:10020/api/whatsnew
```

Returns a JSON list of the most recent commits (2 per repo) from the GitHub Org.

---

## 📎 Credits

MIT License © DASMLAB 2025

Built with:  
- [Gin](https://github.com/gin-gonic/gin)  
- [Logrus](https://github.com/sirupsen/logrus)  
- [Swaggo](https://github.com/swaggo/swag)  
- [ginprom](https://github.com/Depado/ginprom)  

---


