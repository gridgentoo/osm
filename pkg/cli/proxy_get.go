package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/k8s"
)

// GetEnvoyProxyConfig returns the sidecar envoy proxy config of a pod
func GetEnvoyProxyConfig(clientSet kubernetes.Interface, config *rest.Config, namespace string, podName string, localPort uint16, query string) ([]byte, error) {
	// Check if the pod belongs to a mesh
	pod, err := clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Errorf("Could not find pod %s in namespace %s", podName, namespace)
	}
	if !proxyLabelExists(*pod) {
		return nil, errors.Errorf("Pod %s in namespace %s is not a part of a mesh", podName, namespace)
	}
	if pod.Status.Phase != corev1.PodRunning {
		return nil, errors.Errorf("Pod %s in namespace %s is not running", podName, namespace)
	}

	dialer, err := k8s.DialerToPod(config, clientSet, podName, namespace)
	if err != nil {
		return nil, err
	}

	portForwarder, err := k8s.NewPortForwarder(dialer, fmt.Sprintf("%d:%d", localPort, constants.EnvoyAdminPort))
	if err != nil {
		return nil, errors.Errorf("Error setting up port forwarding: %s", err)
	}

	var envoyProxyConfig []byte
	err = portForwarder.Start(func(pf *k8s.PortForwarder) error {
		defer pf.Stop()
		url := fmt.Sprintf("http://localhost:%d/%s", localPort, query)

		// #nosec G107: Potential HTTP request made with variable url
		resp, err := http.Get(url)
		if err != nil {
			return errors.Errorf("Error fetching url %s: %s", url, err)
		}

		envoyProxyConfig, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Errorf("Error rendering HTTP response: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Errorf("Error retrieving proxy config for pod %s in namespace %s: %s", podName, namespace, err)
	}

	return envoyProxyConfig, nil
}

// proxyLabelExists returns a boolean indicating if the pod is part of a mesh
func proxyLabelExists(pod corev1.Pod) bool {
	// osm-controller adds a unique label to each pod that belongs to a mesh
	proxyUUID, proxyLabelSet := pod.Labels[constants.EnvoyUniqueIDLabelName]
	return proxyLabelSet && isValidUUID(proxyUUID)
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}