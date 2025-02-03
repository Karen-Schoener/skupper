# Test case: minikube, metallb; create nginx server in west, addressable by name

## Step: create namespaces: east, west
```
kubectl delete namespace east
kubectl delete namespace west

kubectl create namespace east
kubectl create namespace west
```

## Step: create nginx deployment in west 
```
kubectl apply -n west -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
EOF
```

## Step: create sevice for nginx deployment in west 
```
kubectl apply -n west -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
  type: LoadBalancer
EOF
```

```
kubectl -n west get service
```

## Step: Note the external-ip for the nginx-service
```
kubectl -n west get service
```

## Step: Assign DNS name for the external ip
```
echo "192.168.49.242   nginx-west.local" | sudo tee -a /etc/hosts
```

## Step: create test pod in east 
```
kubectl apply -n east -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: test-container
    image: curlimages/curl:latest
    command: ["sleep", "infinity"]
EOF
```

## Step: from test pod in east, run curl
```
kubectl exec -it -n east test-pod -- /bin/sh
curl 192.168.49.242
curl http://nginx-west.local
```

## Step: Unassign DNS name for the external ip
```
sudo sed -i '/192.168.49.242   nginx-west.local/d' /etc/hosts
```
