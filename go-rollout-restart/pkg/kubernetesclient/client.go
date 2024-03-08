package kubernetesclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient struct {
	clientset *kubernetes.Clientset
}

func NewKubernetesClient() (*KubernetesClient, error) {

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error building clientset: %v", err)
	}

	return &KubernetesClient{clientset: clientset}, nil
}

func (c *KubernetesClient) ListDeployments(name string) ([]*appsv1.Deployment, error) {
	alldeployments, err := c.clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deploymentsList []*appsv1.Deployment
	for _, deployment := range alldeployments.Items {
		if strings.Contains(deployment.Name, name) {
			currentDeployment := deployment
			deploymentsList = append(deploymentsList, &currentDeployment)
		}
	}

	return deploymentsList, nil
}

func (c *KubernetesClient) RolloutRestartDeployment(deployment *appsv1.Deployment) error {

	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err := c.clientset.AppsV1().Deployments(deployment.Namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error updating deployment: %v", err)
	}

	return nil
}
