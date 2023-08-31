package test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

type HelmAddonData struct {
	namespaceName   string
	releaseName     string
	dependencyRepo  string
	addonName       string
	addonAlias      string
	chartPath       string
	hasCustomValues bool
	manageNamespace bool
	overrideValues  map[string]string
}

func readYamlFile(filename string) (*map[string]interface{}, error) {

	var err error

	bufferedContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	mapped := make(map[string]interface{})

	err = yaml.Unmarshal(bufferedContent, &mapped)
	if err != nil {
		return nil, fmt.Errorf("values not found in file %q: %w", filename, err)
	}

	return &mapped, err
}

func getDependenciesFromYamlFile(data map[string]interface{}, dependencies []string) (map[string]interface{}, error) {

	dependenciesValues := make(map[string]interface{})

	for _, v := range data["dependencies"].([]interface{}) {
		for _, depValue := range dependencies {
			if k, value := v.(map[string]interface{})[depValue]; value {
				dependenciesValues[depValue] = k
			} else {
				return nil, fmt.Errorf("key %s not found in dependencies field", depValue)
			}
		}
	}

	return dependenciesValues, nil

}

func prepareHelmEnvironment(t *testing.T, addonData *HelmAddonData) (helm.Options, error) {
	helmOptions := helm.Options{}

	if addonData.namespaceName == "" {
		addonData.namespaceName = addonData.addonName
	}
	if addonData.releaseName == "" {
		addonData.releaseName = fmt.Sprintf("%s-test-%v", addonData.addonName, strings.ToLower(random.UniqueId()))
	}
	if addonData.dependencyRepo == "" {
		addonData.dependencyRepo = fmt.Sprintf("terratest-%s-%v", addonData.addonName, strings.ToLower(random.UniqueId()))
	}
	if addonData.chartPath == "" {
		addonData.chartPath = fmt.Sprintf("../addons/%v", addonData.addonName)
	}

	helmChartPath, err := filepath.Abs(addonData.chartPath)

	if err != nil {
		t.Errorf("Error processing %s error = %s", addonData.chartPath, err)
		return helmOptions, err
	}

	yamlContent, err := readYamlFile(fmt.Sprintf("%s/Chart.yaml", helmChartPath))

	if err != nil {
		t.Errorf("Error reading chart yaml file, error = %s", err)
		return helmOptions, err
	}

	dependencies, err := getDependenciesFromYamlFile(*yamlContent, []string{"alias", "repository"})

	addonData.addonAlias = dependencies["alias"].(string)

	if err != nil {
		t.Errorf("Error reading chart yaml file, error = %s", err)
		return helmOptions, err
	}

	kubectlOptions := k8s.NewKubectlOptions("", "", addonData.namespaceName)

	if addonData.manageNamespace {
		k8s.CreateNamespace(t, kubectlOptions, addonData.namespaceName)
	}

	valuesFiles := []string{}
	if addonData.hasCustomValues {
		valuesFiles = []string{fmt.Sprintf("%s/values.yaml", helmChartPath)}
	}
	helmOptions = helm.Options{
		KubectlOptions: kubectlOptions,
		ValuesFiles:    valuesFiles,
		SetValues:      addonData.overrideValues,
	}

	helm.AddRepo(t, &helmOptions, addonData.dependencyRepo, dependencies["repository"].(string))

	_, err = helm.RunHelmCommandAndGetOutputE(t, &helmOptions, "dependency", "build", helmChartPath)

	if err != nil {
		t.Errorf("Chart dependencies instalation failed, error = %s", err)
		return helmOptions, err
	}

	helm.Install(t, &helmOptions, helmChartPath, addonData.releaseName)

	return helmOptions, err
}

func destroyHelmEnvironment(t *testing.T, addonData HelmAddonData, helmOptions helm.Options) {
	helm.Delete(t, &helmOptions, addonData.releaseName, true)

	helm.RemoveRepo(t, &helmOptions, addonData.dependencyRepo)

	if addonData.manageNamespace {
		k8s.DeleteNamespace(t, helmOptions.KubectlOptions, addonData.namespaceName)
	}
}

func waitUntilHelmFormattedServicesAvailable(t *testing.T, addonData HelmAddonData, helmOptions helm.Options, services []string) {
	for _, v := range services {
		k8s.WaitUntilServiceAvailable(t, helmOptions.KubectlOptions, fmt.Sprintf("%s-%s", addonData.releaseName, v), 10, 30*time.Second)
	}
}

func waitUntilServicesAvailable(t *testing.T, kubectlOptions k8s.KubectlOptions, services []string) {
	for _, v := range services {
		k8s.WaitUntilServiceAvailable(t, &kubectlOptions, v, 10, 30*time.Second)
	}
}

type StatefulSetJsonStruct struct {
	Status struct {
		ReadyReplicas int `json:"readyReplicas"`
	} `json:"status"`
}

func waitUntilStatefulSetsAvailable(t *testing.T, kubectlOptions k8s.KubectlOptions, statefulSets []string) (success bool) {
	tries := 10
	readySS := 0
	for _, v := range statefulSets {
		currentTries := 0
		ready := false
		for currentTries < tries && !ready {
			ssstatus, err := k8s.RunKubectlAndGetOutputE(t, &kubectlOptions, "get", "statefulset", v, "-o=json")
			if err == nil {
				var ssstatusJson StatefulSetJsonStruct
				err = json.Unmarshal([]byte(ssstatus), &ssstatusJson)
				if err == nil {
					if ssstatusJson.Status.ReadyReplicas > 0 {
						ready = true
						readySS++
					}
				}
			}
			currentTries++
			time.Sleep(30 * time.Second)
		}
	}
	return readySS == len(statefulSets)
}

func waitUntilLoadBalancerAvailable(t *testing.T, kubectlOptions k8s.KubectlOptions) {
	for _, v := range k8s.ListServices(t, &kubectlOptions, v1.ListOptions{}) {
		if v.Spec.Type == "LoadBalancer" {
			k8s.WaitUntilServiceAvailable(t, &kubectlOptions, v.Name, 10, 30*time.Second)
		}
	}
}

func waitUntilDaemonSetAvailable(t *testing.T, kubectlOptions k8s.KubectlOptions, daemonSetName string) {
	retries := 10
	sleep := time.Second * 1
	for i := 1; i < retries; i++ {
		podsReady := k8s.GetDaemonSet(t, &kubectlOptions, daemonSetName).Status.NumberReady
		if podsReady > 0 {
			break
		}
		time.Sleep(sleep)
	}
	pods := k8s.ListPods(t, &kubectlOptions, v1.ListOptions{})
	require.Greater(t, len(pods), 0)
	pod := pods[0]

	k8s.WaitUntilPodAvailable(t, &kubectlOptions, pod.Name, 10, 30*time.Second)
}
