# Execute in VM Task

This task can execute a script, or a command in a Virtual Machine

## `execute-in-vm`

### Installation

Install the Task

```bash
kubectl apply -f https://raw.githubusercontent.com/suomiy/kubevirt-tekton-tasks/master/tasks/execute-in-vm/manifests/execute-in-vm.yaml
```

Install one of the following rbac permissions to the active namespace
  - Permissions for executing in VMs from active namespace
    ```bash
    kubectl apply -f https://raw.githubusercontent.com/suomiy/kubevirt-tekton-tasks/master/tasks/execute-in-vm/manifests/execute-in-vm-namespace-rbac.yaml
    ```
  - Permissions for executing in VMs from the cluster
    ```bash
    TARGET_NAMESPACE="$(kubectl config current-context | cut -d/ -f1)"
    wget -qO - https://raw.githubusercontent.com/suomiy/kubevirt-tekton-tasks/master/tasks/execute-in-vm/manifests/execute-in-vm-cluster-rbac.yaml | sed "s/TARGET_NAMESPACE/$TARGET_NAMESPACE/" | kubectl apply -f -
    ```

### Parameters

- **vmName**: Name of a VM to execute the action in.
- **vmNamespace**: Namespace of a VM to execute the action in (defaults to active namespace).
- **secretName**: Secret to use when connecting to a VM.
- **command**: Command to execute in a VM.
- **args**: Arguments of a command.
- **script**: Script to execute in a VM

### Secret format

The secret is used for storing credentials used in VM authentication. Following fields are recognized: 

- **type**: One of: ssh. Defaults to ssh when empty.
- **user**: Username (base64 encoded).
- **privateKey**: Private key to use for authentication (base64 encoded).
- **additionalOptions**: Additional arguments to pass to the SSH command (base64 encoded).

#### Example secret

```yaml

kind: Secret
apiVersion: v1
metadata:
  name: example-credentials
data:
  type: ssh
  # root
  user: Y205dmRBPT0=
  # -C -p 8022
  additionalSSHOptions: LUMgLXAgODAyMg==
  # -----BEGIN OPENSSH PRIVATE KEY-----
  # b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABlwAAAAdzc2gtcn
  # ...
  identity: >-
    LS0tLS1CRUdJTiBPUEVOU1NIIFBSSVZBVEUgS0VZLS0tLS0KYjNCbGJuTnphQzFyWlhrdGRqRUFBQUFBQkc1dmJtVUFBQUFFYm05dVpRQUFBQUFBQUFBQkFBQUJsd0FBQUFkemMyZ3RjbgpOaEFBQUFBd0VBQVFBQUFZRUExLzVVRTA4MFl3S1d3MlBhWVVISXZKQWYyZUJ4YVlMTzhYN3RTMHl4d1oxaGJYYzNSQTdjCnp3cEJ6UFZNMHVmSDdjUGRmU0ZQUmp3YlQrR0o1UERFTlZQc0ZRVHZ0YmpFYnJkQkRMWXRoRjVra3huVWMyVHZKSTJDWDUKOThRZ2FiYkxoY3k0Z2hUa1QwZzdPV3JKMWZENGRWcTArVEN0Q0tYbEQwK3FQZnBKZ051bnlSN3NTT2NHdXQzV3h2WnRSYgoxK2JmNlhFWEpHZGtHUlNZNEthQTFOOTQ1UXhmTHdsTmFaUE4yVGRacFFUdUNGL2YraDMxbGZzb0ZWbGFBQ2M0NmllUlpyCnEvM2VBM244Ti9wSmdybG4xMHFMeFR4UlpVMGVLa2QzY0o2dzBPZ0d5UHNMR0JEakVRcm5mRVJmS0hhRXdla0FHT041RHcKN3lGTmY2N04zTEFkM1NiRnE2b3dFYXl1a0hGRWQ2cHh5UGRQMll6T0JzbTM5TGV4TjJqOVp2Vk5abDhYSnVDS1grZ1FwZwpaUU9ubVhvSTgvdk9xZFFUTXh1V2hvODJzMVMrUjczeUxIcW9IQ2RIUTBGTjVTYjVYTUJNMjNhdXJ6STJPNkdXOWEzb2d4Ck9KZkd3UU9JVzlSNUlmK1N0dkN1N1YzcGlwV29lK1U4SndxQlM4bWZBQUFGZ0FqWDJxUUkxOXFrQUFBQUIzTnphQzF5YzIKRUFBQUdCQU5mK1ZCTlBOR01DbHNOajJtRkJ5THlRSDluZ2NXbUN6dkYrN1V0TXNjR2RZVzEzTjBRTzNNOEtRY3oxVE5Mbgp4KzNEM1gwaFQwWThHMC9oaWVUd3hEVlQ3QlVFNzdXNHhHNjNRUXkyTFlSZVpKTVoxSE5rN3lTTmdsK2ZmRUlHbTJ5NFhNCnVJSVU1RTlJT3pscXlkWHcrSFZhdFBrd3JRaWw1UTlQcWozNlNZRGJwOGtlN0VqbkJycmQxc2IyYlVXOWZtMytseEZ5Um4KWkJrVW1PQ21nTlRmZU9VTVh5OEpUV21UemRrM1dhVUU3Z2hmMy9vZDlaWDdLQlZaV2dBbk9Pb25rV2E2djkzZ041L0RmNgpTWUs1WjlkS2k4VThVV1ZOSGlwSGQzQ2VzTkRvQnNqN0N4Z1E0eEVLNTN4RVh5aDJoTUhwQUJqamVROE84aFRYK3V6ZHl3CkhkMG14YXVxTUJHc3JwQnhSSGVxY2NqM1Q5bU16Z2JKdC9TM3NUZG8vV2IxVFdaZkZ5YmdpbC9vRUtZR1VEcDVsNkNQUDcKenFuVUV6TWJsb2FQTnJOVXZrZTk4aXg2cUJ3blIwTkJUZVVtK1Z6QVROdDJycTh5Tmp1aGx2V3Q2SU1UaVh4c0VEaUZ2VQplU0gva3Jid3J1MWQ2WXFWcUh2bFBDY0tnVXZKbndBQUFBTUJBQUVBQUFHQkFJSjJVeEFxdWlqSHFvcXdhSHB3UXJCcW1XCkovcCtFa0twMUJubERxTFVBV01IL3ZEMG1qZEluMTF4TFdZNlhNYkQxWUhoTnpxQzl6Sk1nV0liVlFIQ05kNG00N3doSlEKWFEzRzBKcTkrelpGMFpmc2RSK1NLc1E0V1RRcTMrSDluYW1Jbmp3dHFIaGZQQlhFeHREU2FWVzIxKzhXckdsTzF4VHE3MQozbzE2L1NUQkNGSjV4Qnl3dHFidVVjbHhsVVVYbFNNWk1PUlAwdkdnM3dHL3N0ZHdtZFJhZUwxbjVScnhtRE5rWHFvZElECjdVYlBxamlvR2V1U2ZYNlByc1lDNmNoY1JjRWN6VDRtVUFxT3BRK2xsTjBuS2gwOUNZbWNTWXNsZEVYSW96MU14SlZKQ3IKMFZUbU9mdXVnSVIrZ25HWXRnclNpNzVsRlllVTlXdXVGcDdyYzZrL0xzZkNZL3BvQ05DMlBMZ01UbjFWQjROK0RBUVlaUQpVekVhQnphU3RyUkwvUThBWTREREs4d0xtdjlGa1EwTGtDempTYThvU0dpZTZsand5R0pGVjVYRkJBV3dVbHRaUExpL2l3Ci9NeXVhcjNlMUdWc1RPMjZpS3hjN2NCV1Q3WERacGhQdzJib25tUVloVTlXMWhYZ2ZaYnluRFo2WFBFU3VJMzBTVzZRQUEKQU1CQk4rSG0zckFBejk2dkZGWC83Y09va2w3emdpbThvbkRwYWpDRkxXOU16alhvdU1xSmdBZzZRSXJ2amdtWlJsdkV6MgozOGU3MHFOY3NOVk54TTR5UFFwL3c5SnZuMDZqeS8rTS8zVTVlSXIwQncwYi9jRkVPU2FxWHh2SFM1cWdYdEtMVTFWRTFTCjM2bUI1RTMrdVU4M3dSZ1NFdDNUSGVzQ3h4RitpbHBoQy9rTzVnUzdNak5wN1RySCt6QnBxNWQyN3l1RStXZVBWbVRTU3QKZTdhLy9XM2xLczVqSERhN3ovdHFreWtId1JwRURaaDBOd2dEenVhdndSTUsyZVRhTUFBQURCQVBadEROQ2hGblJUSE9MVApiZmJhb2J2VVBGUThsNlIvRk95akc1cEt4d09KRDRxWHdVNUQ3em5QSFRaRDNMZC9lQ3ZCb2lDMWNMZ3FzS3FlK3RselprCm9LOFpPQXQyOVEya2FWNzZMbm1aVEZqQ2lSdElpaEwyMnpZbE1wNno0WE9tbXJ5a3pVbnorNmMxZXd1RUFSbVVIQlZkaWYKQ05vamdHRnZkR0JhTEZtSXRyYXlHV0JyTDdKVHQzVTUxTDhMaCtWanp3SWoxZVZ3R1pHcHNwQk5aZkRFOFE2UzVOYmkrTQpYM2FyUm5rWmpmQnQ0dk4xNkMyYi9NV1pyMzdWWENkUUFBQU1FQTRHS1kxUksvQTBIdC85L09pUSt6SEFZVjk4RGNKNWJmCm5mWWVFSlkrZEpVVGlDSGhnY0xEblNDQno5WlFKd2NSUTVvQ2xlWW5sMWt5T28vNjNUSjFzaWhQeHA5YTcrMjN3VzRRbncKd2NscERJMVVTNElqaDZxMnA4ei94eS9aN1hXUVpmK0hjYjdVbkZrOTdiRXIxNktGSHdpZGFaMnZ3eWtMcG12M2U3ZkRXQgp6WDB2UFJzekZic0xhdDFaMWdQZk5LWm50Y21hRFhXZkYxY3ZINGFzZFYzSmdEaUJ1WVkyV2F0S0RRbHVyS0dHL2R4S3kxCnE0RVBLRnM1NWI2N0ZEQUFBQUNXRnVjM2xBWVc1emVRRT0KLS0tLS1FTkQgT1BFTlNTSCBQUklWQVRFIEtFWS0tLS0t
type: Opaque

```

### Usage

Please see [examples](examples)
