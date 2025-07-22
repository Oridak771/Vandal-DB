# Phantom-DB

Phantom-DB is a Kubernetes-native tool for creating masked, anonymized, or synthesized clones of your production databases. It is designed to be used in CI/CD pipelines to provide realistic data for testing and development.

## Features

-   **Kubernetes Native:** Phantom-DB is built on top of Kubernetes and uses Custom Resource Definitions (CRDs) to manage database profiles and clones.
-   **Database Support:** Phantom-DB supports PostgreSQL, with plans to support other databases in the future.
-   **Data Masking:** Phantom-DB can mask, anonymize, or synthesize data using a variety of transformation rules.
-   **Storage Integration:** Phantom-DB integrates with various storage providers to create snapshots of your databases.
-   **Helm Chart:** Phantom-DB is packaged as a Helm chart for easy installation and management.

## Getting Started

To get started with Phantom-DB, you will need a Kubernetes cluster with a CSI driver that supports volume snapshots.

1.  **Install the CRDs:**
    ```
    kubectl apply -f config/crd
    ```
2.  **Install the Controller Manager:**
    ```
    kubectl apply -f config/manager
    ```
3.  **Create a `DataProfile`:**
    ```
    kubectl apply -f apis/v1alpha1/samples/dataprofile.yaml
    ```
4.  **Create a `DataClone`:**
    ```
    kubectl apply -f apis/v1alpha1/samples/dataclone.yaml
    ```

## Documentation

-   [API Reference](docs/api.md)
-   [Getting Started Guide](docs/getting-started.md)
-   [Troubleshooting Guide](docs/troubleshooting.md)

## Contributing

Contributions are welcome! Please see our [contributing guidelines](CONTRIBUTING.md) for more information.
