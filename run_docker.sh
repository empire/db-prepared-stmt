#!/usr/bin/env bash

exitfn () {
  trap SIGINT

  exit
}

trap "exitfn" INT

docker run -it --rm --name=mysql-container -p 3306:3306 -e MYSQL_DATABASE=database -e MYSQL_ROOT_PASSWORD=password mariadb:latest
