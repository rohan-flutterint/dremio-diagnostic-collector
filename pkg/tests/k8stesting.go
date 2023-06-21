package tests

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var cmd *exec.Cmd

// update getClientset with the new function
func GetClientset() (*kubernetes.Clientset, error) {
	c, err := getConfig()
	if err != nil {
		return nil, err
	}
	context := c.Kube.Context
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: context}

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		configOverrides,
	)

	config, err := kubeconfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// Define a struct to hold the configuration data
type K8sTestConfig struct {
	Kube struct {
		Context string `json:"context"`
	} `json:"kube"`
}

func getRootProjectDir() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	rootDir := strings.TrimSpace(string(output))
	return rootDir, nil
}

func getConfig() (K8sTestConfig, error) {
	var config K8sTestConfig

	rootDir, err := getRootProjectDir()
	if err != nil {
		return config, err
	}

	// Read the file
	file, err := os.ReadFile(filepath.Join(rootDir, "test.json"))
	if err != nil {
		return config, err
	}

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

const dremioTestPort = 12831

func SetupDremioK8STestEnv() (string, error) {
	// Create a random namespace name
	namespace := fmt.Sprintf("test-ns-%d", rand.Int())
	clientset, err := getClientset()
	if err != nil {
		return "", err
	}

	// Check if that name is already present on k8s
	_, err = clientset.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err == nil {
		return "", fmt.Errorf("namespace %s already exists", namespace)
	}

	// Create the namespace
	ns, err := clientset.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	})
	if err != nil {
		return "", err
	}

	// Apply a k8s yaml to the namespace
	// This part would vary depending on what you're applying. You might use the client-go library or shell out to kubectl.

	// Wait for all the resources to be up

	podName := "dremio-master-0"
	podPort := 9047
	localPort := dremioTestPort
	cmd = exec.Command("kubectl", "port-forward", podName, fmt.Sprintf("%s:%s", localPort, podPort), "--context", context)

	return ns.Name, nil
}

func DremioTestPort() int {
	return dremioTestPort
}
func TeardownDremioK8sTestEnv(namespace) error {
	clientset, err := getClientset()
	if err != nil {
		return err
	}

	// Stop port-forwarding
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
	}

	return clientset.CoreV1().Namespaces().Delete(namespace, nil)
}
