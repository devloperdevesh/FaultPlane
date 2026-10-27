# Kubernetes Sidecar Installer (Preview)

The installer is a page in the FaultPlane UI (`/installer`) that generates Kubernetes YAML and Helm
values from a form. It is a **preview manifest generator**:

- It does **not** connect to a cluster, deploy anything or read live cluster state.
- It does **not** show deployment status or metrics. The page states "not connected (preview)".
- The generated files are for you to review, copy or download, and apply yourself.

## What it generates

| Output | Use |
|---|---|
| Sidecar container | A container entry to paste under `spec.template.spec.containers` of an existing workload. |
| Deployment + Service | A standalone FaultPlane Deployment with a ClusterIP Service. |
| Helm values | A `values.yaml` for a chart you maintain. This repository does not ship a Helm chart yet. |

## Configuration

| Field | Notes |
|---|---|
| Namespace, workload name | RFC 1123 DNS labels (lowercase letters, digits, `-`, max 63). The namespace must already exist. |
| Replicas | 1-50. Used by the standalone Deployment and Helm values. |
| Container port | 1024-65535. Ports below 1024 are rejected because the container runs as non-root. Written to `FAULTPLANE_PORT`. |
| Image repository, tag | Editable. The default repository is a placeholder: replace it with the image you publish. The tag `latest` produces a warning. |
| Log level | `debug`, `info`, `warn`, `error`. Written to `FAULTPLANE_LOG_LEVEL`. |
| Primary / fallback upstream | Optional `host:port`. Written to `FAULTPLANE_PRIMARY_ADDR` and `FAULTPLANE_FALLBACK_ADDR`. A fallback requires a primary. |
| CPU / memory requests and limits | Kubernetes quantities such as `100m`, `0.5`, `128Mi`, `1Gi`. Requests must not exceed limits. |
| Probes | Readiness and liveness HTTP probes on `/health`. |
| Prometheus annotations | Adds `prometheus.io/*` pod annotations pointing at `/api/v1/metrics/prometheus`. |

Invalid input shows an error next to the field, and no manifest is generated until it is fixed.

## Security posture

Generated containers are unprivileged: `allowPrivilegeEscalation: false`, `runAsNonRoot: true` and all
capabilities dropped. The installer does not generate privileged mode, host networking or added
capabilities. Features that need privileged kernel (eBPF) access are therefore not configured by this
tool. See `docs/deployment.md` for the current deployment status.

## Known limitations

- `FAULTPLANE_HOST`, `FAULTPLANE_PORT` and `FAULTPLANE_LOG_LEVEL` are read by `internal/config`.
  `FAULTPLANE_PRIMARY_ADDR` and `FAULTPLANE_FALLBACK_ADDR` currently appear in the repository's
  docker-compose example; confirm that your build consumes them before relying on them.
- No published image tag is documented in this repository, so the image reference is yours to set.
- IPv6 literals are not accepted in upstream addresses.

## Applying a manifest

Run these yourself:

```bash
kubectl create namespace <namespace>                      # if it does not exist yet
kubectl apply --dry-run=client -f faultplane-deployment.yaml
kubectl apply -f faultplane-deployment.yaml
```

## Development

The generator is plain TypeScript in `ui/lib/installer/` (no runtime dependencies), and the form is
`ui/components/installer/ManifestInstaller.tsx`.

```bash
cd ui
npm install
npm test        # Vitest: validation rules and generated YAML
npm run lint
npm run build
npm run dev     # open http://localhost:3000/installer
```

The tests parse the generated YAML and assert on its structure, including that no privileged settings
are ever emitted.
