package test

import (
	"log"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func TestSecretReplication(t *testing.T) {

	// Helm deploy the kubernetes-replicator chart
	addonData := HelmAddonData{
		namespaceName:   "kubernetes-replicator",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "kubernetes-replicator",
		addonAlias:      "",
		chartPath:       "",
		hasCustomValues: true,
		manageNamespace: true,
	}
	kubectlOptions := k8s.NewKubectlOptions("", "", addonData.namespaceName)
	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	// Target Namespace for secret replication
	k8s.CreateNamespace(t, kubectlOptions, "initium-test1")
	k8s.CreateNamespace(t, kubectlOptions, "initium-test2")
	k8s.CreateNamespace(t, kubectlOptions, "initium")

	replicateToNamespaces := []string{"initium", "initium-test2", "initium-test1"}
	secretName := "secret-replicate"

	// Create a Secret with `replicator.v1.mittwald.de/replicate-to` annotation in default namespace
	secretData := map[string]string{
		"project": "initium",
	}
	byteSecretData := make(map[string][]byte)
	for k, v := range secretData {
		byteSecretData[k] = []byte(v)
	}

	// Define the Secret object with the replicator annotation
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: addonData.namespaceName,
			Annotations: map[string]string{
				"replicator.v1.mittwald.de/replicate-to": "initium,initium-.+,initium.+",
			},
		},
		Data: byteSecretData,
	}

	// Convert the Secret object to YAML
	secretYAML, err := yaml.Marshal(secret)
	if err != nil {
		t.Fatalf("Failed to convert secret to YAML: %v", err)
	}

	// Create the Secret
	k8s.KubectlApplyFromString(t, kubectlOptions, string(secretYAML))

	// Wait for 20 seconds to allow the operator to replicate the secret
	time.Sleep(20 * time.Second)

	// Check each namespace for the replicated secret
	for _, namespace := range replicateToNamespaces {
		// Update k8s options to the target namespace
		targetOptions := k8s.NewKubectlOptions("", "", namespace)

		// Try to get the replicated secret from the target namespace
		replicatedSecret, err := k8s.GetSecretE(t, targetOptions, secretName)

		// Check if the replicated Secret is not nil and has the correct data
		assert.NoError(t, err)
		assert.NotNil(t, replicatedSecret)
		assert.Equal(t, byteSecretData, replicatedSecret.Data)
	}

	// Cleanup
	error := k8s.KubectlDeleteFromStringE(t, kubectlOptions, string(secretYAML))
	assert.NoError(t, error)

	destroyHelmEnvironment(t, addonData, helmOptions)
	k8s.DeleteNamespace(t, kubectlOptions, "initium-test1")
	k8s.DeleteNamespace(t, kubectlOptions, "initium-test2")
	k8s.DeleteNamespace(t, kubectlOptions, "initium")
}
