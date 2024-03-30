#!/usr/bin/env bash

docker exec -it proxysql mysql -h127.0.0.1 -P6032 -uroot -ppass --prompt "ProxySQL Admin> " -e '
  SELECT * FROM stats_mysql_users where username = "app";
  SELECT srv_host, ConnUsed, ConnFree, ConnOK, ConnERR, MaxConnUsed
    FROM stats_mysql_connection_pool;
  SELECT query, global_stmt_id, ref_count_client, ref_count_server, digest
    FROM stats_mysql_prepared_statements_info;'

docker exec -it proxysql mysql -h127.0.0.1 -P6032 -uroot -ppass --prompt "ProxySQL Admin> " -e 'show processlist;' | grep mydb | shuf -n 5 | sort -n

mysql --defaults-extra-file=config.cnf -P3307 mysql -e 'show full processlist;' | grep mydb | shuf -n 5 | sort -n

mysql --defaults-extra-file=config.cnf -P3308 mysql -e 'show full processlist;' | grep mydb | shuf -n 5 | sort -n
