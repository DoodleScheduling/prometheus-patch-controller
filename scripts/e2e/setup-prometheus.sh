#!/bin/bash
kubectl create ns prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update prometheus-community
helm install prometheus prometheus-community/prometheus \
  --version "15.6.0" \
  --namespace prometheus \
  --wait \
  --set kubeStateMetrics.enabled=false \
  --set nodeExporter.enabled=false \
  --set alertmanager.enabled=false \
  --set pushgateway.enabled=false
