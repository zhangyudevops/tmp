#!/bin/bash
# this script delete k8s deployment and service
# monitor the pod status and wait for the pod to be deleted before creating new deployment and service

# delete deployment and service
kubectl delete -f yaml/nginx2-deployment.yml

# sleep 1 second
sleep 5

# monitor the pod status

num=$(kubectl get pods -n tiantong| grep  nginx2 | wc -l)
# if the pod is deleted, break the loop
if [ $num != 0 ]; then
    kubectl get pods -n tiantong| grep  nginx2 | awk '{print $1}' | xargs kubectl delete po -n tiantong --force &> /dev/null
fi


# create deployment and service
kubectl apply -f yaml/nginx2-deployment.yml

