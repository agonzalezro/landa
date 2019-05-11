lambda-control-plane
====================

![](https://estaticos3.larazon.es/documents/10165/0/image_content_low_548139_20130510133304.jpg)

lambda-control-plane (aka landa) is a service that allows you to manage your landa functions. When you create a function it will be deployed on
Kubernetes and you will be able to call it.

Dev'ing
-------

We are using go modules, so a simple:

    cd cmd/landa
    go run main.go -f $HOME/.kube/config

should suffice.

Note that when you don't run this service inside Kubernetes you need to indicate how to reach the cluster config.

Using it
--------

    http localhost:9094/functions code=xxx
    http localhost:9094/functions/[id] # where id is the id returned by the previous call
    http localhost:9094/functions/[id]:call # NOT IMPLEMENTED YET

TODO
----

- [ ] Use lambda-engine instead the nginx image we are using as an example.
- [ ] Allow to call functions.
