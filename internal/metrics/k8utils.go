// SPDX-License-Identifier: Apache-2.0

package metrics

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	PodNameEnvKey      = "SAMBA_POD_NAME"
	PodNamespaceEnvKey = "SAMBA_POD_NAMESPACE"
)

type kclient struct {
	ClientSet *kubernetes.Clientset
	Config    *rest.Config
}

func newKClient() (*kclient, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return newExternalClient()
	}
	cset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return &kclient{}, err
	}

	return &kclient{
		ClientSet: cset,
		Config:    cfg,
	}, nil
}

func newExternalClient() (*kclient, error) {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		return &kclient{}, err
	}
	cset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &kclient{}, err
	}

	return &kclient{
		ClientSet: cset,
		Config:    config,
	}, nil
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube/config")
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func getRunningPod(ctx context.Context, clnt *kclient,
	nname types.NamespacedName) (*corev1.Pod, error) {
	pod, err := clnt.ClientSet.CoreV1().
		Pods(nname.Namespace).Get(ctx, nname.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if pod.Status.Phase != corev1.PodRunning {
		return pod, fmt.Errorf("pod %+v no running", nname)
	}
	return pod, nil
}

func GetSelfPod(ctx context.Context, clnt *kclient) (*corev1.Pod, error) {
	id := GetSelfPodID()
	if len(id.Name) == 0 || len(id.Namespace) == 0 {
		return nil, fmt.Errorf("failed to resolve self id %+v", id)
	}
	return getRunningPod(ctx, clnt, id)
}

func GetSelfPodID() types.NamespacedName {
	return types.NamespacedName{
		Namespace: os.Getenv(PodNamespaceEnvKey),
		Name:      os.Getenv(PodNameEnvKey),
	}
}
