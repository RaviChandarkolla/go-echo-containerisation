creating configmap from CLI

kubectl create configmap app-config \
  --from-literal=APP_ENV=dev \
  --from-literal=APP_DEBUG=true


=================================================

creating configmap from file

kubectl create configmap app-config \
  --from-file=application.properties


=================================================

creating volume from configmap

Idea:
Kubernetes takes keys in a ConfigMap and projects them as files inside the container.

Each key → a file

Each value → file content

Mounted as a read-only volume

Updates propagate automatically (unlike env vars)


=================================================