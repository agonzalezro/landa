#!/bin/zsh

kubectl get -o=go-template="{{range .items}}{{.metadata.name}} {{end}}" deployments|tr "[:space:]" "\n"|xargs kubectl delete deployment
kubectl get -o=go-template="{{range .items}}{{.metadata.name}} {{end}}" service|tr "[:space:]" "\n"|grep -v kubernetes|xargs kubectl delete service

read -r -d '' FUNCTION_CODE << EndOfFunction
package chispas;

import org.linuxfoundation.events.kubecon.lambda.context.Context;

public class Chispas {

    public static class InData {
        private String foo;
        private String bar;

        public String getFoo() {
            return foo;
        }

        public void setFoo(String foo) {
            this.foo = foo;
        }

        public String getBar() {
            return bar;
        }

        public void setBar(String bar) {
            this.bar = bar;
        }

        public String toString() {
            return foo + " " + bar;
        }
    }

    public static class OutData {
        private String result;

        public String getResult() {
            return result;
        }

        public void setResult(String result) {
            this.result = result;
        }
    }

    public OutData doChispas(InData in, Context ctx) {

        System.out.println(in);

        OutData out = new OutData();

        out.setResult(in.toString());

        return out;
    }
}
EndOfFunction

ID=$(http localhost:8080/functions code="$FUNCTION_CODE"|jq .id|tr -d "\"")
echo localhost:8080/functions/$ID':call'
http localhost:8080/functions/$ID':call' foo=bar