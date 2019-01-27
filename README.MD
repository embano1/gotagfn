# About
OpenFaaS function (Golang) using the [OpenFaaS vCenter Connector](https://github.com/openfaas-incubator/vcenter-connector) to add a tag to a VM when it is *powered on*.

# Usage
## Requirements
These steps require a running [OpenFaaS](https://docs.openfaas.com/deployment/) and [vSphere (vCenter)](https://docs.vmware.com/en/VMware-vSphere/index.html) environment. Also, the [OpenFaaS vCenter Connector](https://github.com/openfaas-incubator/vcenter-connector) must set up correctly before deploying this example.

**Note:** If you want to change the function handler, you need a working [Go](https://golang.org/) environment. `dep` is needed for the dependencies. Use `dep ensure -vendor-only` inside the `gotag-fn` folder.

## Deployment
### Supported `stack.yaml` parameters
Modify the function `stack.yaml` as per your environment:

- Either add the OpenFaaS `gateway: <URL|IP>` information under `provider` or `export` it as an environment variable used by `faas-cli`
- Modify environment variables:
  - VC        - URL to vCenter (e.g. `https://vcenter.local`)
  - VC_USER   - Username to connect to vCenter
  - VC_PASS   - Password to connect to vCenter
  - INSECURE  - Don't verify certificate of vCenter when connecting with TLS 
  - TAG_URN   - Unique ID for the tag to be attached (e.g. `urn:vmomi:InventoryServiceTag:604e66d9-97d4-414c-82f5-d25475320087:GLOBAL`)
- A key-value annotation defines which VM event should trigger the function
  - A list of VM events from vCenter can be found [here](https://code.vmware.com/doc/preview?id=4206#/doc/vim.event.VmEvent.html)
  - The `topic` key-value annotation uses `.` syntax, e.g. `VmPoweredOnEvent` maps to `vm.powered.on`

**Note:** OpenFaaS provides some useful [debugging](https://docs.openfaas.com/deployment/troubleshooting/) configured via environment variables `write_debug` and `read_debug`.
   
### How to retrieve the `TAG_URN`
The `TAG_URN` can be retrieved from the ID field in the object, e.g. with `govc`:

```bash
$ govc tags.info <tag>            # example uses "vmon" as the tag name
Name:           vmon
  ID:           urn:vmomi:InventoryServiceTag:4c57977e-27fa-4392-b344-0a90a502f524:GLOBAL
  Description:  VM powered on
  CategoryID:   urn:vmomi:InventoryServiceCategory:ee941e4d-56ce-4ebf-86ab-0f615828d585:GLOBAL
  UsedBy: []
```

Link to `govc`: https://github.com/vmware/govmomi/tree/master/govc

### Deploy the Function

```bash
$ faas-cli deploy -f stack.yml
Deploying: gotag-fn.

Deployed. 202 Accepted.
```

### Invoke the Function

If the event topic is `vm.powered.on` the function can be triggered via the vCenter Client UI or using `govc`:

```bash
$ govc vm.power -on '*'
```

### Logs/Debugging

If you have deployed OpenFaaS in Kubernetes you can query the function log with:

```bash
kubectl -n openfaas-fn logs gotag-fn-<POD_ID> -f
Forking - ./handler []
2019/01/22 17:27:52 Started logging stderr from function.
2019/01/22 17:27:52 Started logging stdout from function.
2019/01/22 17:27:52 OperationalMode: http
2019/01/22 17:27:52 Writing lock-file to: /tmp/.lock
2019/01/22 17:28:08 stderr: 2019/01/22 17:28:08 successfully tagged VM VirtualMachine:vm-51 with tag urn:vmomi:InventoryServiceTag:4c57977e-27fa-4392-b344-0a90a502f524:GLOBAL
```