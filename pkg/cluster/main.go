package cluster

import (
	"context"

	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	functionLabelKey = "function"
	lambdaServerPort = 9443
)

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func New(configPath string) (*Cluster, error) {
	config, err := buildConfig(configPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Cluster{clientset: clientset}, nil
}

type Cluster struct {
	clientset *kubernetes.Clientset
}

func (c *Cluster) DeployFunction(ctx context.Context, id, code string) (string, error) {
	if err := c.createDeployment(ctx, id, code); err != nil {
		return "", err
	}
	return c.createService(ctx, id)
}

func (c *Cluster) createDeployment(_ context.Context, id, code string) error {
	var replicas int32 = 1

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					functionLabelKey: id,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						functionLabelKey: id,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  id,
							Image: "nginx:1.12", // TODO: change to the lambda server
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: lambdaServerPort,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := c.clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Create(deployment)
	return err
}

func (c *Cluster) createService(_ context.Context, id string) (string, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				functionLabelKey: id,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   apiv1.ProtocolTCP,
					Port:       lambdaServerPort,
					TargetPort: intstr.FromInt(lambdaServerPort),
				},
			},
		},
	}

	s, err := c.clientset.CoreV1().Services(apiv1.NamespaceDefault).Create(service)
	return s.Spec.ClusterIP, err
}
