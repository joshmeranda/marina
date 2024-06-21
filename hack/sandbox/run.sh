#!/usr/bin/env bash

source "$(dirname $0)/.env"

values_file="$(dirname $0)/values.yaml"
generated_values_file="$(dirname $0)/generated-values.yaml"
chart_dir="$(dirname $0)/../../charts/marina"

port_ceil=$((32767 + 1))
port_floor=30000
node_port=$(((RANDOM % $(($port_ceil-$port_floor))) + $port_floor))

if k3d cluster list $K3D_CLUSTER_NAME &> /dev/null; then
	printf "Cluster '%s' already exists, doing nothing..." $K3D_CLUSTER_NAME
	exit 1
fi

k3d cluster create --image $K3S_IMAGE $K3D_CLUSTER_NAME \
	--api-port 6443 \
	--port "8081:$node_port@loadbalancer" --port '80:80@loadbalancer' \
	--k3s-arg '--disable=traefik@server:0'

until kubectl cluster-info &> /dev/null ; do
	echo "Waiting for cluster to be ready..."
	sleep 1
done

yq eval ".gateway.service.nodePort = $node_port" "$values_file" > "$generated_values_file"

helm upgrade --install --create-namespace --namespace marina-system marina "$chart_dir" \
	--values "$generated_values_file"