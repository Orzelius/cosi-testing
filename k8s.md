1. get a set of k8s manifest(s) as input (new or updated) `k8sResourceRequest`
2. read manifests and parse as k8s objects
3. create cosi resources containing said resources
4. compare to previous version(s) of the same `k8sResourceRequest` (this step is needed so we know if a resource needs to be deleted)
5. if there's a diff, reconcile to achieve desired state (update/create or delete)