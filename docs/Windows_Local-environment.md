# Windows SO

### Pre-requisites:

**Install WSL (using powershell)**
    - Installing WSL `wsl --install -d Ubuntu` * Debian version do not work well [WSL](https://learn.microsoft.com/en-us/windows/wsl/install) 
    - make ubuntu as default `wsl --set-default ubuntu`
    - update WSL to version 2 `wsl --set-version ubuntu 2`
    - update WSL2 to last packgae version `wsl --update`

**open WSL ubuntu that has been installed**
    - apt update
    - apt install curl git
    - install asdf
        - git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.8.1
        - echo -e "\n. $HOME/.asdf/asdf.sh" >> ~/.bashrc
        - . $HOME/.asdf/asdf.sh
        - asdf plugin-add nodejs https://github.com/asdf-vm/asdf-nodejs.git
        - asdf install nodejs 18.16.0
        - asdf global nodejs 18.16.0
    - apt install make
    - apt install unzip (could be into tools)
    - [Docker](https://docs.docker.com/engine/install/) ( cross-platform, paid solution )

3- Download the project
    - k8s-kurated-addons
        - git clone https://github.com/nearform/k8s-kurated-addons.git 
    - Get into k8s-kurated-addons run:
        - make asdf_install
        - make ci (***From here Docker desktop have to be running***)
        - make

### Preparing App

**Create and prepare your application locally:**
    - Create Dockerfile
    - Create manifest (yaml)
    - build image local
    - push image in your docker hub (have to be public)
    - push your app into your repo on github (have to be public)

**Creating/onfiguring your application to run into th K8s**
    - Copy the manifest to folder k8s-kurated-addons\examples\"nameyourapp\manifest.yaml"
        - have to use the template based on template [example](#\\wsl.localhost\Ubuntu\root\k8s-kurated-addons\examples\manifest\Example.yaml)
    - run `kubectl apply -f \examples\"nameyourapp\manifest.yaml`
    - kubectl get ksvc (Have to be with status running)
    - kubectl get pods

**Testing the app if is running on Cluster**
    - on Linux:
         ` export KKA_LB_ENDPOINT="$(kubectl get service -n istio-ingress istio-ingressgateway -o go-template='{{(index .status.loadBalancer.ingress 0).ip}}'):80" `
         ` export SAMPLE_APP_URL=$(kubectl get ksvc -o json | jq -r '.items[] | select(.metadata.name == "helloworld") | .status.url'  | sed 's#http://##') `



