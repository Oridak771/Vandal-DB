# Getting Started

This guide will walk you through the process of setting up Vandal-DB and creating your first database clone.

## Prerequisites

-   A Kubernetes cluster (v1.20+)
-   A CSI driver that supports volume snapshots
-   `kubectl` installed and configured

## Installation

1.  **Install the CRDs:**
    ```
    kubectl apply -f https://raw.githubusercontent.com/vandal-db/vandal-db/main/config/crd/vandal.db.io_dataprofiles.yaml
    kubectl apply -f https://raw.githubusercontent.com/vandal-db/vandal-db/main/config/crd/vandal.db.io_dataclones.yaml
    ```
2.  **Install the Controller Manager:**
    ```
    kubectl apply -f https://raw.githubusercontent.com/vandal-db/vandal-db/main/config/manager/manager.yaml
    ```

## Creating a Clone

1.  **Create a `DataProfile`:**
    Create a `dataprofile.yaml` file with the following content:
    ```yaml
    apiVersion: vandal.db.io/v1alpha1
    kind: DataProfile
    metadata:
      name: postgres-profile-example
    spec:
      schedule: "0 0 * * *"
      retentionPolicy:
        count: 3
      target:
        secretName: postgres-credentials
      masking:
        rules:
          - table: users
            column: email
            transformation: hash
    ```
    Apply the `DataProfile`:
    ```
    kubectl apply -f dataprofile.yaml
    ```
2.  **Create a `DataClone`:**
    Create a `dataclone.yaml` file with the following content:
    ```yaml
    apiVersion: vandal.db.io/v1alpha1
    kind: DataClone
    metadata:
      name: postgres-clone-example
    spec:
      sourceProfile: postgres-profile-example
      ttl: "1h"
    ```
    Apply the `DataClone`:
    ```
    kubectl apply -f dataclone.yaml
    ```

## Accessing the Clone

Once the `DataClone` is ready, you can access the cloned database using the connection information in the generated secret:
```
kubectl get secret postgres-clone-example -o jsonpath='{.data}'
