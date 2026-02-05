creating configmap from CLI

kubectl create configmap app-config \
  --from-literal=APP_ENV=dev \
  --from-literal=APP_DEBUG=true


================================================

creating configmap from file

kubectl create configmap app-config \
  --from-file=application.properties

