# ===== Sanity check =====
for env in ['KKA_REPO_NAME', 'KKA_REPO_HOST_PATH', 'KKA_REPO_NODE_PATH', 'KKA_REPO_URI', 'KKA_REPO_BRANCH']:
    if os.getenv(env, '') == '': fail('Missing or empty {} env var. Did you run this project using the Makefile?'.format(env))

# ===== Internal variables =====
ARGOCD_EXTERNAL_PORT = os.getenv('KKA_ARGOCD_EXTERNAL_PORT', 8443)
ISTIO_HTTP_PORT = os.getenv('KKA_ISTIO_HTTP_PORT', 7080)
ISTIO_HTTPS_PORT = os.getenv('KKA_ISTIO_HTTPS_PORT', 7443)

# ===== Extensions =====
load('ext://namespace', 'namespace_yaml')

# ===== Kubernetes deployment =====

## MetalLB
local('helm dependency update ./utils/metallb')
k8s_yaml(namespace_yaml('metallb-system'), allow_duplicates=False)
k8s_yaml(helm(
    './utils/metallb',
    name='metallb',
    namespace='metallb-system',
    set=['cidrBlock="{}"'.format(os.getenv('KKA_METALLB_CIDR'))]
))

if os.getenv('KKA_DEPLOY_MINIMAL', 'false') == 'false':
    ## Git HTTP Backend
    docker_build('k8s-kurated-addons/git-http-backend', './utils/git-http-backend/docker')
    k8s_yaml(namespace_yaml('argocd'), allow_duplicates=False)
    k8s_yaml(helm(
        './utils/git-http-backend/chart',
        name='git-http-backend',
        namespace='argocd',
        set=['volumes.git_volume.path={}'.format(os.getenv('KKA_REPO_NODE_PATH'))]
    ))

    ## ArgoCD
    local('helm dependency update ./addons/argocd')
    k8s_yaml(helm(
        './addons/argocd',
        name='argocd',
        namespace='argocd',
    ))

    ## App-of-apps
    k8s_yaml(local('cat manifests/app-of-apps.yaml | envsubst'))
    k8s_resource(
    objects=['k8s-kurated-addons:Application:argocd'],
    new_name='k8s-kurated-addons',
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
        serve_cmd='kubectl port-forward -n argocd svc/argocd-server {}:443'.format(ARGOCD_EXTERNAL_PORT),
        links=['http://localhost:{}'.format(ARGOCD_EXTERNAL_PORT)],
        readiness_probe=probe(
            initial_delay_secs = 20,
            timeout_secs = 1,
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
