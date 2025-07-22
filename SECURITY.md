# Security

This document describes the security best practices for using Phantom-DB.

## RBAC

Phantom-DB uses a least-privilege RBAC model. The `phantom-db-controller-manager` runs as a `ServiceAccount` with a `ClusterRole` that grants it only the permissions it needs to function.

## Secret Encryption

It is recommended to use a secret encryption solution such as [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets) or [Vault](https://www.vaultproject.io/) to encrypt the database credentials used by Phantom-DB.

## Network Policies

It is recommended to use network policies to restrict network traffic to and from the `phantom-db-controller-manager` and the cloned database pods. A default-deny policy is provided in `config/network/network_policy.yaml`.

## Pod Security Policies

It is recommended to use pod security policies to restrict the capabilities of the `phantom-db-controller-manager` and the cloned database pods.
