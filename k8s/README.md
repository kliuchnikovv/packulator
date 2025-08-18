# Kubernetes Deployment

This directory contains Kubernetes manifests for deploying Packulator to a Kubernetes cluster.

## Files

- `namespace.yaml` - Creates the packulator namespace
- `configmap.yaml` - Configuration values
- `secret.yaml` - Sensitive configuration (database password)
- `postgres.yaml` - PostgreSQL database deployment
- `deployment.yaml` - Main application deployment
- `service.yaml` - Service to expose the application
- `ingress.yaml` - Ingress for external access

## Deployment

1. Apply the namespace first:
```bash
kubectl apply -f namespace.yaml
```

2. Apply configuration and secrets:
```bash
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
```

3. Deploy PostgreSQL:
```bash
kubectl apply -f postgres.yaml
```

4. Wait for PostgreSQL to be ready, then deploy the application:
```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

5. (Optional) Apply ingress for external access:
```bash
kubectl apply -f ingress.yaml
```

## Verification

Check if pods are running:
```bash
kubectl get pods -n packulator
```

Check service endpoints:
```bash
kubectl get svc -n packulator
```

## Notes

- The secret contains a base64 encoded password. Change it before deploying to production.
- Update the ingress host in `ingress.yaml` to match your domain.
- Adjust resource requests/limits in `deployment.yaml` based on your requirements.
- For production, consider using external PostgreSQL service instead of running it in the cluster.