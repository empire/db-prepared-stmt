# Database Prepared Statement Simulator

This repository explores the inner workings of Go's SQL and MySQL drivers, uncovering crucial insights for optimizing database interactions and minimizing errors. These insights have been obtained through an in-depth study of the source code of Go's SQL standard library and the Go MySQL driver.

## Getting Started

To get started, clone the source code repository:

```
git clone git@github.com:empire/db-prepared-stmt.git
```

Ensure that you have the Go, K6, and Docker binaries installed on your system.

## Starting the MariaDB Cluster

Navigate to the `galera` directory and start the MariaDB cluster using Docker Compose:

```
cd galera
docker compose up --remove-orphans
```

## Executing the Service

Once the cluster is up and running, execute the service by running the following command:

```
go run .
```

## Simulating Sample Load

To simulate a sample load, you can either use the `k6` command or run the `run_watched.sh` script. 

## Monitoring MySQL Resources

For monitoring MySQL resources and process lists conveniently, a script called `monitor.sh` has been provided. Execute it using the following command:

```
watch -n 1 ./monitor.sh
```

## Customizing Configuration

If you need to customize the configuration, open the `galera/docker-compose.yml` and `galera/proxysql/proxysql.cnf` files. Make the necessary changes to the variables, and then restart the Docker Compose command.
