package test

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	_ "github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/k8s"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestSamplePrometheusApp(t *testing.T) {
	// =============================================================
	kubectlOptions := k8s.NewKubectlOptions("", "", "istio-ingress")

	k8s.WaitUntilServiceAvailable(t, kubectlOptions, "istio-ingressgateway", 10, 30*time.Second)

	// =============================================================
	kubectlOptions = k8s.NewKubectlOptions("", "", "knative-serving")

	k8s.WaitUntilServiceAvailable(t, kubectlOptions, "controller", 10, 5*time.Second)

	// =============================================================
	kubectlOptions = k8s.NewKubectlOptions("", "", "prometheus")

	k8s.KubectlApply(t, kubectlOptions, "sample-prometheus-app.yaml")

	// =============================================================
	tunnel := k8s.NewTunnel(kubectlOptions, k8s.ResourceTypePod, "prometheus-kube-prometheus-stack-kube-prometheus-0", 0, 9090)

	defer tunnel.Close()

	tunnel.ForwardPort(t)

	// =============================================================
	endpoint := fmt.Sprintf("http://%s%s", tunnel.Endpoint(), "/api/v1/targets?state=active")

	tlsConfig := tls.Config{}

	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		endpoint,
		&tlsConfig,
		30,
		3*time.Second,
		verifySamplePrometheusTarget,
	)

	// =============================================================
	k8s.KubectlDelete(t, kubectlOptions, "sample-prometheus-app.yaml")
}

func verifySamplePrometheusTarget(statusCode int, body string) bool {
	if statusCode != 200 {
		return false
	}
	return strings.Contains(body, "podMonitor/prometheus/sample-prometheus-app")
}
