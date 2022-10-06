source ~/.bash_profile

minikube start

minikube addons enable ingress

minikube addons enable dashboard

dapr init -k

linkerd install --crds | kubectl apply -f -

linkerd install --set proxyInit.runAsRoot=true | kubectl apply -f -

linkerd viz install | kubectl apply -f -
