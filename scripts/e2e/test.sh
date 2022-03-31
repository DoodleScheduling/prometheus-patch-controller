#!/bin/bash
kubectl -n default apply -f config/testdata/rule.yaml
kubectl -n default wait prometheuspatchrule/annotate-namespace --for=condition=Active --timeout=1m
kubectl -n default wait prometheuspatchrule/annotate-namespace --for=condition=PatchApplied --timeout=15s
test $(kubectl get ns default -o=jsonpath='{.metadata.annotations.foo}') == "bar"
