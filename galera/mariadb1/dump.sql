GRANT ALL PRIVILEGES ON mydb.* To 'app'@'%' IDENTIFIED BY 'pass';

FLUSH PRIVILEGES;

use mydb;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255)
);

insert into users (id, name) values (1, "Alice"), (2, "Bob"), (3, "Charlie");
