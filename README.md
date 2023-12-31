# dataset-operator

Experiment for Dataset orchestration Operator.

It learns from https://github.com/kubevirt/containerized-data-importer
```
apiVersion: dataset-ops.my.domain/v1
kind: Dataset
metadata:
  name: "peter-dataset"
spec:
  source:
      http:
         url: "https://harbor-test2.cn-sh2.ufileos.com/docs/download/DCE5.0-intro.pdf" #S3 or GCS
         secretRef: "" # Optional
         certConfigMap: "" # Optional
  pvc:
    storageClassName: local-path
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: "64Mi"
```

Effect:

```
#######################################
# a dataset CR
#######################################

# kubectl get dataset
NAME            PHASE       PROGRESS   RESTARTS   AGE
peter-dataset   Completed      100%               8m

#######################################
###  temporary pod as a pre-loader/data-importer 
#######################################

# kubectl get po
NAME                           READY   STATUS             RESTARTS       AGE
importer-peter-dataset         0/1     Completed   6 (117s ago)   8m3s

#######################################
### check the underlying PV, yes, the data has been downloaded (so far, save as a.out for now)
#####################################
# kubectl get pvc
NAME                  STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
peter-dataset         Bound     pvc-e27757ee-a8c3-42cd-a01f-8f2455af6a02   64Mi       RWO            local-path     8m6s

## kubectl get pv pvc-e27757ee-a8c3-42cd-a01f-8f2455af6a02 -o yaml |grep local-path-provisioner
    path: /opt/local-path-provisioner/pvc-e27757ee-a8c3-42cd-a01f-8f2455af6a02_default_peter-dataset

# du -sh /opt/local-path-provisioner/pvc-e27757ee-a8c3-42cd-a01f-8f2455af6a02_default_peter-dataset/a.out
1.4M	/opt/local-path-provisioner/pvc-e27757ee-a8c3-42cd-a01f-8f2455af6a02_default_peter-dataset/a.out
```





## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/dataset-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified. 
And it is required to have access to pull the image from the working environment. 
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/dataset-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

