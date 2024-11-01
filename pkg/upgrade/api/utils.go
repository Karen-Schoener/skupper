package api

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/skupperproject/skupper/internal/kube/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func GetKubeConfig() string {
	var kubeconfig string // TODO
	return kubeconfig
}

func ReadConfigMap(ctx context.Context, namespace string, name string, cli *client.KubeClient) (*corev1.ConfigMap, error) {
	kubeClient := cli.GetKubeClient()
	cm, err := kubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return cm, err
}
