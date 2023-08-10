# ===== Sanity check =====
for env in ['INITIUM_REPO_NAME', 'INITIUM_REPO_HOST_PATH', 'INITIUM_REPO_NODE_PATH', 'INITIUM_REPO_URI', 'INITIUM_REPO_BRANCH']:
    if os.getenv(env, '') == '': fail('Missing or empty {} env var. Did you run this project using the Makefile?'.format(env))

def parse_excluded_env_vars(prefix='INITIUM_AOA_EXCLUDE_'):
    mapping = {
        'knative': [
            'knative-operator',
            'knative-serving',
            'knative-eventing'
        ],
        'prometheusstack': [
            'kube-prometheus-stack'
        ]
    }

    excluded = [
        var.removeprefix(prefix).lower().replace('_', '-')
        for var in os.environ.keys()
        if var.startswith(prefix) and os.getenv(var) == 'true'
    ]

    values = []
    for e in excluded:
        values.extend(mapping.get(e, [e]))

    return ['apps.%s.excluded=true' % app for app in values ]

# ===== Internal variables =====
ARGOCD_EXTERNAL_PORT = os.getenv('INITIUM_ARGOCD_EXTERNAL_PORT', 8080)
ISTIO_HTTP_PORT = os.getenv('INITIUM_ISTIO_HTTP_PORT', 7080)
ISTIO_HTTPS_PORT = os.getenv('INITIUM_ISTIO_HTTPS_PORT', 7443)

# ===== Extensions =====
load('ext://namespace', 'namespace_yaml')

# ===== Kubernetes deployment =====

def bootstrap_app_values():
    VALUES_OVERRIDES='./manifests/bootstrap/overrides.local.yaml'
    valueFiles = [VALUES_OVERRIDES] if os.path.exists(VALUES_OVERRIDES) else []
    values = ['repoURL=%s' % os.getenv('INITIUM_REPO_URI')]
    values += parse_excluded_env_vars()
    return valueFiles, values


## MetalLB
local('helm dependency update ./utils/metallb')
k8s_yaml(namespace_yaml('metallb-system'), allow_duplicates=False)
k8s_yaml(helm(
    './utils/metallb',
    name='metallb',
    namespace='metallb-system',
    set=['cidrBlock="{}"'.format(os.getenv('INITIUM_METALLB_CIDR'))]
))

if os.getenv('INITIUM_DEPLOY_MINIMAL', 'false') == 'false':
    ## Git HTTP Backend
    docker_build('initium-platform/git-http-backend', './utils/git-http-backend/docker')
    k8s_yaml(namespace_yaml('argocd'), allow_duplicates=False)
    k8s_yaml(helm(
        './utils/git-http-backend/chart',
        name='git-http-backend',
        namespace='argocd',
        set=['volumes.git_volume.path={}'.format(os.getenv('INITIUM_REPO_NODE_PATH'))]
    ))

    ## k8s secret with TLS cert for wildcard.kube.local domains
    k8s_yaml(namespace_yaml('istio-ingress'), allow_duplicates=False)
    k8s_resource(
    objects=['istio-ingress:namespace'],
    new_name='istio-ingress-namespace',
    )
    local_resource(
        'wildcard.kube.local-tls-secret',
        cmd='kubectl create secret tls -n istio-ingress wildcard.kube.local-tls --save-config \
        --dry-run=client --cert=.ssl/cert-ingress-gateway.pem --key=.ssl/key-ingress-gateway.pem -o yaml | kubectl apply -f -',
        auto_init=True,
        resource_deps=['istio-ingress-namespace']
    )

    ## ArgoCD
    local('helm dependency update ./addons/argocd')
    k8s_yaml(helm(
        './addons/argocd',
        name='argocd',
        namespace='argocd',
    ))

    ## App-of-apps
    valueFiles, values = bootstrap_app_values()
    k8s_yaml(helm('./manifests/bootstrap', namespace="argocd", name="app-of-apps", values=valueFiles, set=values))

    k8s_resource(
        objects=['initium-platform:Application:argocd'],
        new_name='initium-platform',
        resource_deps=['argocd-redis', 'argocd-server', 'argocd-repo-server', 'metallb-metallb-source-controller', 'metallb-metallb-source-speaker', 'git-http-backend']
    )

    # ===== Tilt local resources =====

    ## ArgoCD admin password
    local_resource(
        'argocd-password',
        cmd='kubectl get secret -n argocd argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo',
        auto_init=False
    )

    ## ArgoCD HTTPS port
    local_resource(
        'argocd-portforward',
        serve_cmd='kubectl port-forward -n argocd svc/argocd-server {}:80'.format(ARGOCD_EXTERNAL_PORT),
        links=['http://localhost:{}'.format(ARGOCD_EXTERNAL_PORT)],
        readiness_probe=probe(
            initial_delay_secs = 20,
            timeout_secs = 5,
            period_secs = 10,
            success_threshold = 1,
            failure_threshold = 5,
            http_get=http_get_action(port=int(ARGOCD_EXTERNAL_PORT))),
        auto_init=False
    )

    ## Istio Ingress HTTPS port
    local_resource(
        'istio-ingress-portforward-https',
        serve_cmd='kubectl port-forward -n istio-ingress svc/istio-ingressgateway {}:443'.format(ISTIO_HTTPS_PORT),
        links=['https://localhost:{}'.format(ISTIO_HTTPS_PORT)],
        readiness_probe=probe(
            initial_delay_secs = 20,
            timeout_secs = 1,
            period_secs = 10,
            success_threshold = 1,
            failure_threshold = 5,
            http_get=http_get_action(port=int(ISTIO_HTTPS_PORT), scheme='https')
        ),
        auto_init=False
    )

    ## Istio Ingress HTTP port
    local_resource(
        'istio-ingress-portforward-http',
        serve_cmd='kubectl port-forward -n istio-ingress svc/istio-ingressgateway {}:80'.format(ISTIO_HTTP_PORT),
        links=['http://localhost:{}'.format(ISTIO_HTTP_PORT)],
        readiness_probe=probe(
            initial_delay_secs = 20,
            timeout_secs = 1,
            period_secs = 10,
            success_threshold = 1,
            failure_threshold = 5,
            http_get=http_get_action(port=int(ISTIO_HTTP_PORT), scheme='http')
        ),
        auto_init=False
    )
