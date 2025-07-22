# API Reference

This document provides a reference for the Phantom-DB API.

## DataProfile

A `DataProfile` defines a template for creating database clones. It specifies the source database, the snapshot schedule, the retention policy, and the masking rules.

### Spec

| Field | Type | Description |
|---|---|---|
| `schedule` | string | The cron schedule for automated snapshots. |
| `retentionPolicy` | object | The policy for retaining snapshots. |
| `target` | object | The database to be profiled. |
| `masking` | object | The data masking configuration. |

## DataClone

A `DataClone` represents a clone of a database created from a `DataProfile`.

### Spec

| Field | Type | Description |
|---|---|---|
| `sourceProfile` | string | The name of the `DataProfile` to clone from. |
| `snapshotName` | string | The name of the specific snapshot to use. |
| `ttl` | string | The time-to-live for the clone. |
