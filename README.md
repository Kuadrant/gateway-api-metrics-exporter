
# Gateway API Metrics Exporter

## Overview

The Gateway API Metrics Exporter is a tool designed to provide observability for [Gateway API](https://gateway-api.sigs.k8s.io/) resources by exporting relevant metrics.

This project is a successor to [gateway-api-state-metrics project](https://github.com/Kuadrant/gateway-api-state-metrics/blob/main/METRICS.md).
There are [limitations](https://github.com/Kuadrant/gateway-api-state-metrics/issues/1) with what resource information can be exposed using the underlying kube-state-metrics project leveraged by gateway-api-state-metrics.
This project aims to allow additional stateful information about Gateway API resources to be made available via metrics.
For additional information on why this project exists, see https://github.com/Kuadrant/architecture/blob/main/rfcs/0010-gateway-api-metrics-exporter.md

## Features

- **Metrics Collection**: Gathers stateful metrics for Gateway API resources including GatewayClasses, Gateways and HTTPRoutes.
- **Prometheus Integration**: Exports metrics in a Prometheus-compatible format for easy integration.

## Installation

### Prerequisites

- Kubernetes cluster (v1.20+)

### Steps

#### Bringing up a local kind cluster

```shell
kind create cluster
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.1.0/standard-install.yaml
```

#### Running out of cluster

```shell
KUBECONFIG=/path/to/.kube/config go run main.go
```

#### Running in cluster

Details TBD

## Usage

Once deployed, the exporter will automatically start collecting and exporting metrics. You can view these metrics through your Prometheus instance or any compatible monitoring tool.
The metrics scrape endpoint is `/metrics`

If running the exporter locally, you can use curl to try it out:

```shell
curl http://localhost:8080/metrics
```

You can also create some sample resources (no need for a gateway-api provider to be installed, just the CRDs):

```shell
kubectl apply -f ./tests/manifests/
```

If the curl was successful, you should see some metrics like this:

```
curl -s http://localhost:8080/metrics | grep gatewayapi
# HELP gatewayapi_gatewayclass_info Information about a GatewayClass
# TYPE gatewayapi_gatewayclass_info gauge
gatewayapi_gatewayclass_info{name="testgatewayclass1"} 1
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the Apache 2.0 License.
