# Troubleshooting

This guide provides solutions to common problems you may encounter when using Vandal-DB.

## `DataClone` Stuck in `Pending` Phase

If a `DataClone` is stuck in the `Pending` phase, it may be due to one of the following reasons:

-   **No `VolumeSnapshot` available:** Ensure that the `DataProfile` has successfully created a `VolumeSnapshot`. You can check the status of the `DataProfile` and the `VolumeSnapshot` objects to verify this.
-   **Storage provider issues:** There may be an issue with the storage provider or the CSI driver. Check the logs of the `vandal-db-controller-manager` and the CSI driver pods for any errors.
-   **Insufficient resources:** The cluster may not have enough resources to create the `DataClone`. Check the cluster's resource usage and ensure that there are enough resources available.

## `DataProfile` Not Creating Snapshots

If a `DataProfile` is not creating snapshots, it may be due to one of the following reasons:

-   **Invalid cron schedule:** Ensure that the `schedule` field in the `DataProfile` spec is a valid cron expression.
-   **Storage provider issues:** There may be an issue with the storage provider or the CSI driver. Check the logs of the `phantom-db-controller-manager` and the CSI driver pods for any errors.
-   **RBAC permissions:** The `vandal-db-controller-manager` may not have the necessary RBAC permissions to create `VolumeSnapshot` objects. Ensure that the `ClusterRole` and `ClusterRoleBinding` are correctly configured.
