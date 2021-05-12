# Golang RESTful API with Postgres database on Kubernetes

## Description
This is a very basic web api written in Go that allows you to create, delete, and list People using the POST, DELETE, and GET routes defined in `main.go`. I wanted to keep the API simple and focus more on the working parts of Kubernetes. I will explain how I set up a local environment using `Kind`, and how I deployed this same setup on Linode with data persistence with a Postgres database.

## The Code
Feel free to explore the code. I used the Gorm library for its simplicity. You can see check out the docs to learn more [here](https://gorm.io/docs/index.html). To run this code locally, you have a couple of options.

1. Use the following commands:
```bash
# Run a Postgres database in docker and map it to port 5432. Be sure to stop your local Postgres service if you have Postgres installed locally.
docker run --name pg -d --rm -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=golang -p 5432:5432 postgres

# Store the connection string in an environment variable however you prefer to do so. This is just one way.
export DSN="host=localhost user=user password=password dbname=golang port=5432"

# Run the Go web server
go run main.go

# Terminal output
[GIN-debug] GET    /people                   --> github.com/manedurphy/golang-start/api.GetPeople (3 handlers)
[GIN-debug] GET    /person/:id               --> github.com/manedurphy/golang-start/api.GetPerson (3 handlers)
[GIN-debug] POST   /person                   --> github.com/manedurphy/golang-start/api.CreatePerson (3 handlers)
[GIN-debug] DELETE /person/:id               --> github.com/manedurphy/golang-start/api.DeletePersion (3 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

Create a person --> POST http://localhost:8080/person
```json
{
	"name": "Dane Murphy",
	"age": 25,
	"hasDegree": true
}
```
Test the other routes  
Get person by ID --> GET http://localhost:8080/person/:id  
Get all people --> GET http://localhost:8080/people  
Delete person by ID --> DELETE http://localhost:8080/person/:id  

2. Using docker compose
```bash
# Start the containers
docker-compose up --build

# Delete the containers
docker-compose down
```

## Running in Kubernetes locally with Kind
I have created commands in the Makefile to ease the process of setting up the cluster. You need Kind installed locally for this to work. See instructions [here](https://kind.sigs.k8s.io/).
```bash
# Create a 3-node cluster
make cluster

# Ensure the nodes have a "Ready" status
watch kubectl get nodes

NAME                 STATUS   ROLES                  AGE   VERSION
kind-control-plane   Ready    control-plane,master   80s   v1.20.2
kind-worker          Ready    <none>                 50s   v1.20.2
kind-worker2         Ready    <none>                 50s   v1.20.2

# Load docker images into the Kind cluster
make load
```

## Start the application
```bash
# This puts the connection string into a secret using an imperative command
make secret

# This sets up all other objects declaratively with the YAML files in the kubernetes directory
make deploy

# Wait for the pods to reach a "Ready" state
watch kubectl get pods

# Port forward to the "goapp-service"
make forward

kubectl port-forward service/goapp-service 8080:8080
Forwarding from 127.0.0.1:8080 -> 8080
Forwarding from [::1]:8080 -> 8080
```

You should now be able to hit the routes just as before, and the data should persist. Try creating a Person, deleting the Postgres pod, and see if your data is still there.
```bash
# After creating a person on the POST route, delete the Postgres pod created by the StatefulSet
kubectl delete pod pg-statefulset-0

# Wait for a running state
kubectl get pods

NAME                                READY   STATUS    RESTARTS   AGE
goapp-deployment-5c4d6df7bd-k4bng   1/1     Running   0          6m50s
pg-statefulset-0                    1/1     Running   0          2m28s
```
Now try hitting the GET /people route and see if your person is still in the database!

# Deploy on Linode

## Run Linode CLI tool in docker
```bash
make linode
```

You will see in the Makefile that with this command I mount the working directory to the container. This will allow us to use `kubectl` in the container if we want to, as we will be adding our `kubeconfig.yaml` file to the `kubeconfig` directory.

## In the Linode CLI container that we just made, run the following:
```bash
# Get the setup prompt --> You will have three questions for default configurations on your account
# Note, you need and account on linode and a Personal Access Token
linode-cli # -> This will prompt you to paste your token

# Deploy a 3 node cluster
linode-cli lke cluster-create --label go-api --k8s_version 1.20 --node_pools.type g6-standard-1 --node_pools.count 3

┌───────┬────────┬─────────┐
│ id    │ label  │ region  │
├───────┼────────┼─────────┤
│ xxxxx │ go-api │ us-west │
└───────┴────────┴─────────┘
```

From the Linode dashboard, navigate to your cluster and look for the `Kubeconfig`. You can download it on your machine or you can open it in the browser and copy/paste it in your working directory. Create a `kubeconfig.yaml` file in the `kubeconfig` directory and paste that kubeconfig text from the Linode dashboard into your newly create `kubeconfig.yaml` file.

Once this is done, open another terminal window on your local machine and set the following environment variable
```bash
# This ensures that when we enter kubectl commands, that we are referring to this configuration file rather than the default one at ~/.kube/config
export KUBECONFIG=$(pwd)/kubeconfig/kubeconfig.yaml

# Confirm you can see the nodes in your cluster
kubectl get nodes

NAME                          STATUS   ROLES    AGE   VERSION
lke26467-37960-609b4fa060fb   Ready    <none>   13m   v1.20.6
lke26467-37960-609b4fa086f7   Ready    <none>   13m   v1.20.6
lke26467-37960-609b4fa0ab8f   Ready    <none>   13m   v1.20.6
```

You will need to make a few changes before following the steps that we took before with our local setup.  

In the `kubernetes/deployments/appdeployment.yaml` file, you will need to change the image that is used, as well as comment our the `imagePullPolicy` which is currently set to `never`. 

Now in `kubernetes/statefulsets/pgstatefulset.yaml`, you will need to change the storage class. Linode comes with storage classes preset in their clusters. We will be using the `linode-block-storage`, as this will provision an persistent volume for us automatically, and it will automatically delete the volume when we delete the persistent volume claim in our cluster. Read more about it [here](https://www.linode.com/docs/guides/deploy-volumes-with-the-linode-block-storage-csi-driver/).


And your are done! If you want to test the API, you can get the IP of any of the three nodes.
```bash
kubectl describe node <nodename> | grep -i externalip

# Now make a request to http://<ip>:32000/people
```

I have set the service of the application to be of type `NodePort`. If you change the type to `LoadBalancer`, Linode will automatically setup a loadbalancer for the 3 nodes in the cluster that allows you to send requests on port 80.

Be sure to delete your cluster when you are finished to avoid added costs.

## Delete your linode cluster
```bash
# Delete the resources in your cluster
make destroy

# Destory the PVC that was created when we deployed our Postgres StatefulSet to ensure the volume is deleted on Linode
kubectl delete pvc --all

# Get clusterId
linode-cli lke clusters-list

┌───────┬────────┬─────────┐
│ id    │ label  │ region  │
├───────┼────────┼─────────┤
│ xxxxx │ go-api │ us-west │
└───────┴────────┴─────────┘

# Delete cluster by its Id
linode-cli lke cluster-delete <clusterId>
```