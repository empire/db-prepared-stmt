(mysql --defaults-extra-file=config.cnf -P3307 mysql -e 'show full processlist;';
mysql --defaults-extra-file=config.cnf -P3308 mysql -e 'show full processlist;') |
  grep --line-buffered mydb | grep --line-buffered -i sleep ||
  (sleep 1; k6 -q run k6.js > /dev/null)
