# Flight Monitor

## Prerequisites
### Roster Database

The Roster server requires a MariaDB database in order to store the management status of nodes and provide an audit log of previous states.

Install MariaDB per upsteam documentation. The commands below are for a system running Rocky 9.3:
```
cat << EOF > /etc/yum.repos.d/MariaDB.repo
[MariaDB-11.4]
name=MariaDB 11.4 repo
baseurl=https://mirror.mariadb.org/yum/11.4/rocky9-amd64
gpgcheck=0
EOF
yum makecache
yum install -y MariaDB-server MariaDB-client
systemctl enable --now mariadb
mariadb -u root -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';"
mariadb -u root -p -e "DELETE FROM mysql.user WHERE User = '';"
mariadb -u root -p -e "DELETE FROM mysql.user WHERE User = 'root' AND Host NOT IN ('localhost','127.0.0.1','::1');"
mariadb -u root -p -e "DROP DATABASE IF EXISTS test;"
mariadb -u root -p -e "FLUSH PRIVILEGES;"
```

Once installed, create a new database and user for Flight Roster.
```
mariadb -u root -p -e "CREATE DATABASE flightroster;"
mariadb -u root -p -e "CREATE USER 'flightroster'@'localhost' IDENTIFIED BY 'password';"
mariadb -u root -p -e "GRANT ALL ON flightroster.* TO 'flightroster'@'localhost';"
```

Import the database structure.
```
mariadb -u flightroster -p flightroster < database.sql
```

### Go
Ensure Go is installed.
```
yum install -y golang
```
