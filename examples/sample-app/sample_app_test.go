package test

import (
	"fmt"
	"os"
	"testing"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestExampleApp(t *testing.T) {
	// =============================================================
	kubectlOptions := k8s.NewKubectlOptions("", "", "istio-ingress")

	k8s.WaitUntilServiceAvailable(t, kubectlOptions, "istio-ingressgateway", 10, 30*time.Second)

	// =============================================================
	kubectlOptions = k8s.NewKubectlOptions("", "", "knative-serving")

	k8s.WaitUntilServiceAvailable(t, kubectlOptions, "controller", 10, 5*time.Second)

	// =============================================================
	kubectlOptions = k8s.NewKubectlOptions("", "", "default")

	k8s.KubectlApply(t, kubectlOptions, "sample-app.yaml")

	// =============================================================

	numRequests := 5

	for i:=0; i < numRequests; i++ {

		http_helper.HTTPDoWithRetry(
			t,
			"GET",
			fmt.Sprintf("http://%s", os.Getenv("INITIUM_LB_ENDPOINT")),
			[]byte("Hello world from initium-platform!"),
			map[string]string{"Host": "sample-app.default.example.com"},
			200,
			30,
			3*time.Second,
			nil,
		)
	}

	sampleAppLabel := "ksample"
	prefix := "sample-app" 

	//Wait for the pods to be created
	time.Sleep(5 * time.Second)

    podOptions := metav1.ListOptions{
        LabelSelector: fmt.Sprintf("service=%s", sampleAppLabel),
    }

	pods := k8s.ListPods(t, kubectlOptions, podOptions)

	assert.NotNil(t, pods)
	// Validate knative scale-out process was effective
	assert.Greater(t, len(pods), 1)

	// Wait for knative to scale-in
	time.Sleep(70 * time.Second)

    replicaSets, err := k8s.RunKubectlAndGetOutputE(t, kubectlOptions, "get", "rs")
    if err != nil {
        t.Fatal(err)
    }

	replicaSetNames := strings.Fields(replicaSets)

	var replicaSet string
	for _, rs := range replicaSetNames {
		if strings.HasPrefix(rs, prefix) {
			replicaSet = rs
		}
	}

	rs := k8s.GetReplicaSet(t, kubectlOptions, replicaSet)
	// Validate knative scale-in process was effective
    assert.Equal(t, int(rs.Status.Replicas), 0)

	// =============================================================
	k8s.KubectlDelete(t, kubectlOptions, "sample-app.yaml")
}
