## kuberetes manifest sync poc

sync local kubernetes manifests with kubernetes with ssa and diffing

* apply `go run main.go --file manifests.yaml apply`
* apply (dry run with debug logs) `go run main.go --file manifests.yaml apply --dry-run --log-level debug`
* diff `go run main.go --file manifests.yaml diff`
