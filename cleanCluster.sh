#!/usr/bin/env bash
kubectl get -o=go-template="{{range .items}}{{.metadata.name}} {{end}}" deployments|tr "[:space:]" "\n"|xargs kubectl delete deployment
kubectl get -o=go-template="{{range .items}}{{.metadata.name}} {{end}}" service|tr "[:space:]" "\n"|grep -v kubernetes|xargs kubectl delete service