lambda-control-plane
====================

**WARN: this code make use of your .kube/config to deploy the example code. Be careful with your default context to avoid deploying this demo in a non desired cluster**

**Disclaimer: this code isn't production ready, it's just a small demo for Kubecon thought to be executed in a local environment.**

![](https://estaticos3.larazon.es/documents/10165/0/image_content_low_548139_20130510133304.jpg)

lambda-control-plane (aka landa) is a service that allows you to manage your landa functions. When you create a function it will be deployed on
Kubernetes and you will be able to call it.

Dev'ing
-------

Start minikube
```bash
minikube start
```

Enable tunneling (sudo pass usually required)

```bash
minikube tunnel 
```

Start control plane. We are using go modules, so a simple:
```bash
    cd cmd/landa
    go run main.go -f $HOME/.kube/config
```
should suffice.

Note that when you don't run this service inside Kubernetes you need to indicate how to reach the cluster config.

Using it
--------

Consider install HTTPie to use examples from https://httpie.org/ 
        
    http POST localhost:9094/functions code=xxx
    http localhost:9094/functions/[id] # where id is the id returned by the previous call
    http POST localhost:9094/functions/[id]:call

Hacking around
--------------
HTTPie is needed for next scripts. Install it from https://httpie.org/

JQ is needed for next scripts. Install it from https://stedolan.github.io/jq/

* There is a script called `clean.sh` that will clean up your Kubernetes environment.
* Another `deployFunction.sh` script is provided deploying a simple hello world function.