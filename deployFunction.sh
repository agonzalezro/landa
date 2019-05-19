#!/bin/zsh

read -r -d '' FUNCTION_CODE << EndOfFunction
package kubecondemo;

import org.linuxfoundation.events.kubecon.lambda.context.Context;

public class HelloWorld {

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

    public OutData sayHello(InData in, Context ctx) {

        System.out.println(in);

        OutData out = new OutData();

        out.setResult(in.toString());

        return out;
    }
}
EndOfFunction

ID=$(http POST localhost:8080/functions code="$FUNCTION_CODE" entryPoint="kubecondemo.HelloWorld.sayHello"|jq .id|tr -d "\"")
echo $ID
echo "To call your deployed function execute"
echo curl -X POST localhost:8080/functions/${ID}:call -d @functionCall.json