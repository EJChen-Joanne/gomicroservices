## k8s commands:
* [minikube](https://minikube.sigs.k8s.io/docs/start/)
```
minikube start --nodes [number of nodes]
minikube status 
minikube dashboard
minikube stop
minikube tunnel
```
* Deployment - [kubectl](https://kubernetes.io/docs/tasks/tools/)
```
kubectl apply -f [folder or file.yml]
kubectl get pods
kubectl get svc
kubectl get deployments
kubectl delete deployments [deployments name]
kubectl delete svc [deployments name]
```
* hit the service
```
kubectl expose deployment broker-service --type=LoadBalancer --port=80 --target-port=80
minikube tunnel
```
