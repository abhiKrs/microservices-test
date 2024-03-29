minikube delete
until read -r -p "Enter RAM for the Minikube in MB(example 4096 or 8192 or 16384): " mb && test "$mb" != ""; do
  continue
done
until read -r -p "Enter CPU core for the Minikube (example 4 or 8 or 16): " cpu && test "$cpu" != ""; do
  continue
done
minikube start --force --memory=$mb --cpus=$cpu
eval $(minikube docker-env)
until read -r -p "Enter namespace for logfire application: " namespace && test "$namespace" != ""; do
  continue
done
kubectl create namespace $namespace
until read -r -p "Enter namespace for daprsystem: " daprsystem && test "$daprsystem" != ""; do
  continue
done
kubectl create namespace $daprsystem
until read -r -p "Enter namespace for kafka namespace: " kafkanamespace && test "$kafkanamespace" != ""; do
  continue
done
kubectl create namespace $kafkanamespace
until read -r -p "Enter namespace for elastic namespace: " elasticnamespace && test "$elasticnamespace" != ""; do
  continue
done
kubectl create namespace $elasticnamespace
until read -r -p "Enter namespace for certmanager namespace: " certmanager && test "$certmanager" != ""; do
  continue
done
kubectl create namespace $certmanager

echo "created namespace $namespace"
echo "created namespace $daprsystem"
echo "created namespace $kafkanamespace"
echo "created namespace $elasticnamespace"
echo "created namespace $certmanager"
kubectl apply -f secrets.yaml -n $namespace
if [ $(awk '/^MemTotal:/ { print $2; }' /proc/meminfo) -lt 16252928  ]; then

#########Installing DAPR#####################
# Add the official Dapr Helm chart.
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm upgrade --install dapr dapr/dapr \
--version=1.9 \
--namespace $daprsystem \
--create-namespace \
--wait
kubectl get all -n $daprsystem
helm search repo dapr --devel --versions
#########Ending DAPR#########################


#########setup Ingress#######################
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --values ingress/dapr-annotations.yaml --namespace $namespace --wait
kubectl apply -f ingress/local-ingress.yaml
#########Ending local ingress################

#######Deploy postgres######################
kubectl apply -f postgres -n $namespace 
#######Ending postgres######################

########Deploy Application###################
docker build -t logfire/gowebapi:latest gowebapp/web-api
docker build -t logfire/notification:latest gowebapp/notification
docker build -t logfire/livetail:latest gowebapp/livetail
docker build -t logfire/filter-service:latest gowebapp/filter-service
#kubectl apply -f secrets -n $namespace
kubectl apply -f gowebapp -n $namespace
kubectl apply -f gowebapp/goweb-api.yaml -n $namespace
kubectl apply -f gowebapp/livetail.yaml -n $namespace
kubectl apply -f gowebapp/filter-service.yaml -n $namespace
#######Ending Applicaion deploy#############

#######Setup Kafka##########################
kubectl apply -f kafka/strimzi-operator.yaml -n $kafkanamespace --wait
kubectl apply -f kafka/kafka-persistent-single.yaml -n $kafkanamespace --wait
kubectl -n $kafkanamespace run kafka-producer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic
kubectl -n $kafkanamespace run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning
kubectl -n $kafkanamespace run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic filter_topic_21a67b7e-3057-45cc-8da1-cd4521fc7530 --from-beginning
#######Ending Kafka##########################

#######Setup Elasticsearch###################
kubectl apply -f elastic-search
PASSWORD=$(kubectl get secret -n $elasticnamespace my-elasticsearch-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')
curl -u "elastic:$PASSWORD" -k "https://localhost:9200"
#######Ending elasticsearch##################

########FlinkInstallation####################
kubectl apply -f flink -n $namespace --wait
#kubectl apply -f flink/flink-configuration-configmap.yaml -n $namespace --wait
#kubectl apply -f flink/jobmanager-session-deployment.yaml  -n $namespace --wait
#kubectl apply -f flink/jobmanager-service.yaml -n $namespace --wait
#kubectl apply -f flink/taskmanager-session-deployment.yaml -n $namespace --wait

#########Nginx Ingress#######################
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --values ingress/dapr-annotations.yaml --namespace $namespace --wait 
kubectl apply -f ingress/webapp-ingress.yaml
kubectl get ingress
##########Ending Nginx Ingress###############


##########Installing Deploy Apache Zookeeper#
#helm install zookeeper bitnami/zookeeper \ --set replicaCount=1 \  --set auth.enabled=false \  --set allowAnonymousLogin=true \   --namespace kafka \   --create-namespace \   --wait
##########Ending Apache Zookeeper############

###########cert-manager######################
#kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.1/cert-manager.yaml
#kubectl get ns
#kubectl -n $certmanager get all
#kubectl apply -f cert-issuer/staging-issuer.yaml
#kubectl apply -f ingress/staging-ingress.yaml
#kubectl describe ingress webapp-ingress
#sudo chmod 777 priv.key
#sudo chmod 777 crt.key
#sudo kubectl create secret tls apibeta-logfire-tls-secret --key priv.key --cert crt.key
kubectl port-forward -n ingress-nginx  --address 0.0.0.0 service/ingress-nginx-controller 80:80 
kubectl port-forward -n logfire-staging  --address 0.0.0.0 service/api-gateway-ingress-nginx-controller 80:80
###########Ending Cert Manager###############
#####################################################AVAILBALE RAM Is Higher Than 16GB#############################
else

#########Installing DAPR#####################
# Add the official Dapr Helm chart.
helm repo add dapr https://dapr.github.io/helm-charts/
helm repo update
helm upgrade --install dapr dapr/dapr \ --version=1.9 \ --namespace $daprsystem \ --create-namespace \ --wait
kubectl get all -n $daprsystem
helm search repo dapr --devel --versions
#########Ending DAPR#########################


#########setup Ingress#######################
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --values ingress/dapr-annotations.yaml --namespace $namespace --wait
kubectl apply -f ingress/local-ingress.yaml
#########Ending local ingress################

#######Deploy postgres######################
kubectl apply -f postgres -n $namespace 
#######Ending postgres######################

########Deploy Application###################
docker build -t logfire/gowebapi:latest gowebapp/web-api
docker build -t logfire/notification:latest gowebapp/notification
docker build -t logfire/livetail:latest gowebapp/livetail
kubectl apply -f secrets -n $namespace
kubectl apply -f gowebapp -n $namespace
kubectl apply -f gowebapp/goweb-api.yaml -n $namespace
kubectl apply -f gowebapp/livetail.yaml -n $namespace
kubectl apply -f gowebapp/filter-service.yaml -n $namespace
#######Ending Applicaion deploy#############

#######Setup Kafka##########################
kubectl apply -f kafka/strimzi-operator.yaml -n $kafkanamespace --wait
kubectl apply -f kafka/kafka-persistent-single.yaml -n $kafkanamespace --wait
kubectl -n $kafkanamespace run kafka-producer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic
kubectl -n $kafkanamespace run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning
kubectl -n $kafkanamespace run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic filter_topic_21a67b7e-3057-45cc-8da1-cd4521fc7530 --from-beginning
#######Ending Kafka##########################

#######Setup Elasticsearch###################
kubectl apply -f elastic-search
PASSWORD=$(kubectl get secret -n $elasticnamespace my-elasticsearch-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')
curl -u "elastic:$PASSWORD" -k "https://localhost:9200"
#######Ending elasticsearch##################

########FlinkInstallation####################
kubectl apply -f flink/flink-admin.yaml -n $namespace --wait
kubectl apply -f flink/flink-configuration-configmap.yaml -n $namespace --wait
kubectl apply -f flink/jobmanager-session-deployment.yaml  -n $namespace --wait
kubectl apply -f flink/jobmanager-service.yaml -n $namespace --wait
kubectl apply -f flink/taskmanager-session-deployment.yaml -n $namespace --replicas=2 --wait

#########Nginx Ingress#######################
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --values ingress/dapr-annotations.yaml --namespace $namespace --wait 
kubectl apply -f ingress/webapp-ingress.yaml
kubectl get ingress
##########Ending Nginx Ingress###############


##########Installing Deploy Apache Zookeeper#
helm install zookeeper bitnami/zookeeper \ --set replicaCount=3 \  --set auth.enabled=false \  --set allowAnonymousLogin=true \   --namespace kafka \   --create-namespace \   --wait
##########Ending Apache Zookeeper############

###########cert-manager######################
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.1/cert-manager.yaml
kubectl get ns
kubectl -n $certmanager get all
kubectl apply -f cert-issuer/staging-issuer.yaml
kubectl apply -f ingress/staging-ingress.yaml
kubectl describe ingress webapp-ingress
sudo chmod 777 priv.key
sudo chmod 777 crt.key
sudo kubectl create secret tls apibeta-logfire-tls-secret --key priv.key --cert crt.key
ubectl port-forward -n ingress-nginx  --address 0.0.0.0 service/ingress-nginx-controller 80:80 443:443
kubectl port-forward -n logfire-staging  --address 0.0.0.0 service/api-gateway-ingress-nginx-controller 80:80 443:443
fi
