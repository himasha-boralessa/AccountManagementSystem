# DistributedAccountManagementSystem
Containerized application with Google Kubernetes Engine


```
 export my_cluster=cluster-1
 export zone=europe-west1-c
 export PROJECT_ID=<projectId>

 Push docker image to Google Container Registry (GCR)
 docker build -t gcr.io/$PROJECT_ID/account-manager:latest .
 docker push gcr.io/$PROJECT_ID/account-manager:latest
 docker build -t gcr.io/$PROJECT_ID/client:latest .
 docker push gcr.io/$PROJECT_ID/client:latest
 docker build -t gcr.io/$PROJECT_ID/accounts-monitor:latest .
 docker push gcr.io/$PROJECT_ID/accounts-monitor:latest

 gcloud container clusters create $my_cluster --num-nodes 3 --zone $zone --enable-ip-alias //add cluster
 gcloud container clusters create $my_cluster --zone $zone   //add cluster
 gcloud container clusters resize $my_cluster --zone $zone --num-nodes=4 //Modify GKE cluster
 gcloud container clusters get-credentials $my_cluster --zone $zone  //to allow authentication
 kubectl cluster-info  //the cluster information

sed "s/PROJECT_ID/${PROJECT_ID}/g" deployment.yml | kubectl apply -f -

kubectl describe pod <podName> // view the complete details of the Pod

kubectl get pods
kubectl get deployments
kubectl get services  //view details about services in the cluster
kubectl logs <pod-name> -c account-manager
curl http://localhost:8080/transaction?amount=84


```
