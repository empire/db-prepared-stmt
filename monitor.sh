#!/usr/bin/env bash

set -x

docker exec -it proxysql mysql -h127.0.0.1 -P6032 -uroot -ppass --prompt "ProxySQL Admin> " -e 'SELECT * FROM stats_mysql_users;'
docker exec -it proxysql mysql -h127.0.0.1 -P6032 -uroot -ppass --prompt "ProxySQL Admin> " -e 'SELECT * FROM stats_mysql_connection_pool;'
docker exec -it proxysql mysql -h127.0.0.1 -P6032 -uroot -ppass --prompt "ProxySQL Admin> " -e 'SELECT * FROM stats_mysql_prepared_statements_info;'

set +x

mysql --defaults-extra-file=config.cnf -P3307 mysql -e 'show full processlist;' | grep mydb | wc
mysql --defaults-extra-file=config.cnf -P3307 mysql -e 'show full processlist;' | grep mydb | shuf -n 10 | sort -n

mysql --defaults-extra-file=config.cnf -P3308 mysql -e 'show full processlist;' | grep mydb | wc
mysql --defaults-extra-file=config.cnf -P3308 mysql -e 'show full processlist;' | grep mydb | shuf -n 10 | sort -n
