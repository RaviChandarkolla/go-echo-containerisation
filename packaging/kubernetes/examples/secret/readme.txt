creating secrets

from literals

kubectl create secret generic db-secret \
  --from-literal=username=admin \
  --from-literal=password=secret123

==================================================================

from files

kubectl create secret generic tls-secret \
  --from-file=cert.pem \
  --from-file=key.pem


==================================================================
creating using yaml

apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
data:
  username: YWRtaW4=
  password: c2VjcmV0MTIz

==================================================================

Using secret in pods

env:
- name: DB_USER
  valueFrom:
    secretKeyRef:
      name: db-secret
      key: username

==================================================================

using secret as files

volumeMounts:
- name: secret-vol
  mountPath: /etc/secrets
volumes:
- name: secret-vol
  secret:
    secretName: db-secret

files

/etc/secrets/username
/etc/secrets/password

==================================================================

specific keys


secret:
  secretName: db-secret
  items:
  - key: username
    path: db-user


==================================================================

If using env vars → pod restart required

| Method       | Updates automatically? |
| ------------ | ---------------------- |
| Env vars     | ❌ No                   |
| Volume mount | ✅ Yes (with delay)     |

==================================================================


Secret for private docker registry

kubectl create secret docker-registry regcred \
  --docker-server=myregistry.io \
  --docker-username=myuser \
  --docker-password=mypassword \
  --docker-email=my@email.com


Usage in pod

spec:
  imagePullSecrets:
  - name: regcred

==================================================================

Create a TLS secret (HTTPS)

kubectl create secret tls tls-secret \
  --cert=server.crt \
  --key=server.key


YAML Equivalent

apiVersion: v1
kind: Secret
metadata:
  name: tls-secret
type: kubernetes.io/tls
data:
  tls.crt: <base64>
  tls.key: <base64>

==================================================================


Updating a secret

kubectl create secret generic db-secret \
  --from-literal=password=newpass \
  --dry-run=client -o yaml | kubectl apply -f -

==================================================================

RBAC - Allow Read-Only secrets


apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: secret-reader
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get"]


Bind it:
kind: RoleBinding

==================================================================





