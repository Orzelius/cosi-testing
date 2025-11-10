## kuberetes manifest sync poc

sync local kubernetes manifests with kubernetes with ssa and diffing

* apply `go run main.go --file manifests.yaml apply`
* apply (dry run with debug logs) `go run main.go --file manifests.yaml apply --dry-run --log-level debug`
* diff `go run main.go --file manifests.yaml diff`


### Takeaways

#### Apply

Both ssa and kubernetes libraries can do apply. But kubernetes apply logic has to following up-sides:
* Manages state with an inventory
  * Store what's been previously applied to enable pruning abilities
  * Keep track if previous applies were successfull to minimize reruns
  * Keep track if applied objects were reconciled successfully (Example: apply succeeded but Deployment never reaches Available -> `Actuation=Succeeded`, `Reconcile=Failed`.)
* Has status watching functionality

#### Diff

Although kubecli has diffing capabilities the logic is entirelly CLI oriented
and can not be easily used as a library.

The means that flucd's SSA prebuilt logic works better for us, as it doesn't require reverse
hacking of the kubernetes go code.
