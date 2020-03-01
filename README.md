This is a proposal of using a config map in order to map resources
from  ovirt to kubevirt.

Using ConfigMap with an expected data format will provide the ability to
extend the mapping to contain additional information to complete a
proper resource creation on kubevirt side.

The user will have to provide this resource priorly to initiating the VM
import process, in which the VM will be referred to this map.

The patch contains an example of a mapping in
`ovirt-mapping-example.yaml`:

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: ovirt-mapping-example
  namespace: default
data:
  mappings: |-
    networkMapping:
      - source: red # maps of ovirt logic network to network attachment definition
        target: xyz
        type: bridge
      - source: ovirtmgmt
        target: pod
        type: pod
        kind: ConfigMap
    storageMapping:
      - source: ovirt_storage_domain_1 # maps ovirt storage domains to storage class
        target: storage_class_1
    affinityMapping: # affinity mappings will be limited at first for 'pinned-to-host', since affinity/anti-affinity requires more than a single vm import
      - source: ovirt_node_1
        target: k8s_node_1
        policy: pinning
```

An example of the parsed output can be shown by running:
`go run main.go`