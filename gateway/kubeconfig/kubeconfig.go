package kubeconfig

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
)

type kubeNode struct {
	ClusterName string
	Server      string
	Cert        string
	User        string
}

type data struct {
	ClusterName     string
	Host            string
	ClusterID       string
	Cert            string
	User            string
	Username        string
	Password        string
	Token           string
	EndpointEnabled bool
	Nodes           []kubeNode
}

func caCertString() string {
	certFilePath := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	data, err := os.ReadFile(certFilePath)
	if err != nil {
		// todo: less than ideal silent failure
		return ""
	}

	base64Data := base64.StdEncoding.EncodeToString(data)

	return base64Data
}

func getDefaultNode(clusterName, host string) kubeNode {
	return kubeNode{
		Server:      fmt.Sprintf("https://%s", host),
		Cert:        caCertString(),
		ClusterName: clusterName,
		User:        clusterName,
	}
}

func ForTokenBased(clusterName, host, token string) (string, error) {
	data := &data{
		ClusterName:     clusterName,
		Host:            host,
		Cert:            caCertString(),
		User:            clusterName,
		Token:           token,
		Nodes:           []kubeNode{getDefaultNode(clusterName, host)},
		EndpointEnabled: false,
	}

	if data.ClusterName == "" {
		data.ClusterName = data.ClusterID
	}

	buf := &bytes.Buffer{}
	err := tokenTemplate.Execute(buf, data)

	return buf.String(), err
}
