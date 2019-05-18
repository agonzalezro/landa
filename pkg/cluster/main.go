package cluster

import (
	"context"
	"fmt"

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
	functionLabelKey  = "function"
	lambdaServerPort  = 9443
	lambdaDockerImage = "hrodes/kubecon-barcelona-lambda-engine:0.1.0" // TODO: get this from conf or env var
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

	s, err := c.createService(ctx, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://localhost:%v", s.Spec.Ports[0].NodePort), nil // TODO: localhost is a hacky get to get this running quick&dirty on docker for mac
}

func (c *Cluster) buildEnvVars(ctx context.Context, code string) []corev1.EnvVar {
	kvs := map[string]string{
		"FUNCTION_CODE":             code,
		"FUNCTION_ENTRYPOINT":       "chispas.Chispas.doChispas", // TODO (for this and following envs): parametrize or hardcode in the engine
		"functionName":              "chispas.Chispas.doChispas", // TODO: uppercase
		"COMPILE_CLASSPATH":         "/lambda-server/*",
		"FUNCTION_SERVER_CLASSPATH": "/lambda-server/*",
		"BUILD_DIR":                 "/tmp",
		"MAIN_CLASS":                "org.linuxfoundation.events.kubecon.lambda.server.bootstrap.FunctionServerBootstraper",
	}

	var envVars []corev1.EnvVar
	for k, v := range kvs {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}
	return envVars
}

func (c *Cluster) createDeployment(ctx context.Context, id, code string) error {
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
							Name:    id,
							Image:   lambdaDockerImage,
							Command: []string{"starterd"},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: lambdaServerPort,
								},
							},
							Env: c.buildEnvVars(ctx, code),
						},
					},
				},
			},
		},
	}

	_, err := c.clientset.AppsV1().Deployments(apiv1.NamespaceDefault).Create(deployment)
	return err
}

func (c *Cluster) createService(_ context.Context, id string) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				functionLabelKey: id,
			},
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   apiv1.ProtocolTCP,
					Port:       lambdaServerPort, // TODO: this should be random if we want several functions
					TargetPort: intstr.FromInt(lambdaServerPort),
				},
			},
		},
	}

	return c.clientset.CoreV1().Services(apiv1.NamespaceDefault).Create(service)
}
