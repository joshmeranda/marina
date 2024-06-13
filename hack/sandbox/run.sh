#!/usr/bin/env bash

source "$(dirname $0)/.env"

values_file="$(dirname $0)/values.yaml"
chart_dir="$(dirname $0)/../../charts/marina"

port_ceil=$((32767 + 1))
port_floor=30000
node_port=$(((RANDOM % $(($port_ceil-$port_floor))) + $port_floor))

if k3d cluster list $K3D_CLUSTER_NAME &> /dev/null; then
	printf "Cluster '%s' already exists, doing nothing..." $K3D_CLUSTER_NAME
	exit 1
fi

k3d cluster create --image $K3S_IMAGE $K3D_CLUSTER_NAME --port "8081:$node_port@loadbalancer"

until kubectl cluster-info &> /dev/null ; do
	echo "Waiting for cluster to be ready..."
	sleep 1
done

helm upgrade --install --create-namespace --namespace marina-system marina "$chart_dir" \
	--values "$values_file" \
	--set gateway.service.nodePort=$node_port

until go run "$(dirname $0)/../../cmd/marina/main.go" --address localhost:8081 health; do
	echo "Waiting for marina to be ready..."
	sleep 1
done