# DistributedAccountManagementSystem
Containerized application with Google Kubernetes Engine


```
#gcloud config set project <projectId>
 export my_cluster=cluster-1
 export zone=europe-west1-c
 export PROJECT_ID=<projectId>
 export GOOGLE_APPLICATION_CREDENTIALS="./service-account-file.json"

 Create a new bucket:
 gsutil mb gs://<projectId>-bucket/
 Make the bucket public:
 gsutil iam ch allUsers:admin gs://<projectId>-bucket
 gsutil cp -n "" gs://<projectId>-bucket/accounts-data.txt  

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

kubectl apply -f PersistentVolume.yaml
kubectl apply -f PersistentVolumeClaim.yaml
sed "s/PROJECT_ID/${PROJECT_ID}/g" client-manager-deployment.yml | kubectl apply -f -
sed "s/PROJECT_ID/${PROJECT_ID}/g" accounts-monitor-deployment.yml | kubectl apply -f -

kubectl describe pod <podName> // view the complete details of the Pod

kubectl get pods
kubectl get deployments
kubectl get services  //view details about services in the cluster
kubectl logs <pod-name> -c <container-name>




```
