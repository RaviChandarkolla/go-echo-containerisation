Steps to run in minikube
1. kubectl port-forward service/system-service 8001:8001
2. execute this request - "postman request POST 'http://localhost:8001/dummyApi' \
  --header 'Content-Type: application/json' \
  --body '{
    "name": "admin",
    "description": "desc"
}'"