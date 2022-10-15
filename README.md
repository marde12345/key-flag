# key-flag 

Key-flag is middleware for consul KV feature flags backed by ristretto and redis as cache and postgreSQL as Database

## Pre-requisite
1. Install soda https://gobuffalo.io/documentation/database/soda/

## First Time Development

```shell
make compose-up  # to enable all docker requirement
```


## Development

### Starting Development Mode
```shell
make compose-up
make run-dev
```

### Stopping Development Mode and Clear All Docker
```shell
make compose-down
```