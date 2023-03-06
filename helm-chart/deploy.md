- ## add nginx-ingress-controller
ps
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --namespace $namespace --wait
```
<br>

- ## Deploy Application

```ps
kubectl create ns logfire-local
cd microservices/helm-chart
################ Run operators ############
helm repo add strimzi https://strimzi.io/charts/
helm install my-strimzi-operator strimzi/strimzi-kafka-operator --namespace logfire-staging
helm install crds operators/crds/ -n logfire-local
helm install es-operator operators/es-operator/ -n logfire-local

###############End Operator#############

###############Run Secret and ConfigMap####
kubectl apply -f configmap/ -n logfire-local
kubectl apply -f secrets -n logfire-local
##############End Secet and ConfigMap######

############## Install Application ############

helm install postgres psql-chart/ -n logfire-local
helm install flink flink-chart/ -n logfire-local
helm install kafka kafka-chart/ -n logfire-local
helm install webapi webapp-chart/ -n logfire-local

############# Install application completed########

##############To Uninstall application#############
helm uninstall postgres psql-chart/ -n logfire-local
helm uninstall flink flink-chart/ -n logfire-local
helm uninstall kafka kafka-chart/ -n logfire-local
helm uninstall webapi webapp-chart/ -n logfire-local
##############Uninstall completed##################
```
<!-- kubectl apply -f deploy-local/operators -n $namespace --wait
kubectl apply -f deploy-local/deployment -n $namespace --wait -->

# Staging





# Production
