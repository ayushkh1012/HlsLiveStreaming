#!/bin/bash

# Create Kind cluster if it doesn't exist
if ! kind get clusters | grep -q "livehls"; then
    kind create cluster --name livehls --config kind-config.yaml
fi

# Build the Docker image
docker build -t livehls:latest .

# Load the image into Kind
kind load docker-image livehls:latest --name livehls

# Install/upgrade the Helm chart
helm upgrade --install livehls ./helm/livehls \
    --namespace livehls \
    --create-namespace \
    --wait

# Wait for pod to be ready
kubectl wait --for=condition=ready pod -l app=livehls -n livehls --timeout=60s

# Get the NodePort and cluster IP
NODE_PORT=$(kubectl get -n livehls -o jsonpath="{.spec.ports[0].nodePort}" services livehls-livehls)
CLUSTER_IP=$(kubectl get -n livehls -o jsonpath="{.spec.clusterIP}" services livehls-livehls)

echo "Application is available at http://localhost:$NODE_PORT"
echo "Cluster IP: $CLUSTER_IP"

# Verify the service
kubectl get pods,svc -n livehls



