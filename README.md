# 🧠 WhatsNew Service – Org-Wide GitHub Activity Tracker

---

## 🏗️ High Level Design

```mermaid
graph TD
  A[GitHub Org Repos] -->|JWT Auth| B(WhatsNew Service)
  B --> C{Background Commit Cache}
  B --> D[Main API Server (:10020)]
  B --> E[Metrics Server (:9200)]
  D --> F[Swagger UI (/swagger/index.html)]
  E -->|Prometheus scrape| G[Grafana Alloy]
```

---

## 🚀 Features

- 🔐 **GitHub App OAuth2 Auth**
  - JWT-based token exchange using a GitHub App created at the organization level
  - Access is limited to selected repositories
  - No personal access tokens (PATs) required
  - Scoping and revocation are controlled by the org admin

- ⚙️ **Fully Dockerized + CI/CD Ready**
  - Docker multi-stage build
  - GitHub Actions workflow with FluxCD-compatible GitOps output

- 📈 **Out-of-band Prometheus Metrics**
  - Dedicated metrics port (`:9200`) using `ginprom`
  - Avoids mixing telemetry and business endpoints

- 🔁 **Dynamic Repo Discovery**
  - GitHub App token fetches live org repo list each cycle

- 🧰 **Tech Stack**
  - Go 1.21+
  - Gin, Logrus, Swaggo, ginprom
  - GitHub App-based OAuth2 using private `.pem` key

---

## 🧪 Local Development

### 🔧 Requirements

- Go 1.21+
- Docker
- A GitHub App with:
  - Read-only access to selected org repos
  - Downloaded `.pem` file
  - Your `APP_ID` and `INSTALLATION_ID`

### 🔧 Required ENV Vars

| Variable           | Description                                       |
|--------------------|---------------------------------------------------|
| `APP_ID`           | GitHub App ID                                     |
| `INSTALLATION_ID`  | App Installation ID                               |
| `PEMFILE`          | Absolute path to downloaded GitHub App .pem file  |

---

### 🔨 Build Locally

```bash
./buildme.sh
```

### ▶️ Run Locally (Dockerized)

```bash
./runme_local.sh
```

---

## 📦 CI/CD & GitOps Workflow

- GitHub Actions based pipeline
- Builds, tests, and publishes to GHCR
- Generates version-tagged deployment YAML under `k8s_envelope/`
- Auto-syncs changes to FluxCD-monitored GitOps repo

---

## 📊 Metrics

- **Exposed at:** `http://localhost:9200/metrics`
- **Scrapable by:** Prometheus, Grafana Alloy, etc.
- **Exported via:** [ginprom](https://github.com/Depado/ginprom)

---

## 🔍 API and Swagger

Swagger docs are served at:

```
http://localhost:10020/swagger/index.html
```

Try your calls directly from the Swagger UI.

---

## 🧪 Example: cURL Call

```bash
curl http://localhost:10020/api/whatsnew
```

Returns the latest commits (2 per repo) from your GitHub org.

---

## 📎 Credits

MIT License © DASMLAB 2025

Built with:
- [Gin](https://github.com/gin-gonic/gin)
- [Logrus](https://github.com/sirupsen/logrus)
- [Swaggo](https://github.com/swaggo/swag)
- [ginprom](https://github.com/Depado/ginprom)

---

