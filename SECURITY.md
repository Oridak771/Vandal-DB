# Security

This document describes the security best practices for using Vandal-DB.

## RBAC

Vandal-DB uses a least-privilege RBAC model. The `vandal-db-controller-manager` runs as a `ServiceAccount` with a `ClusterRole` that grants it only the permissions it needs to function.

## Secret Encryption

It is recommended to use a secret encryption solution such as [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets) or [Vault](https://www.vaultproject.io/) to encrypt the database credentials used by Vandal-DB.

## Network Policies

It is recommended to use network policies to restrict network traffic to and from the `vandal-db-controller-manager` and the cloned database pods. A default-deny policy is provided in `config/network/network_policy.yaml`.

## Pod Security Policies

It is recommended to use pod security policies to restrict the capabilities of the `vandal-db-controller-manager` and the cloned database pods.
