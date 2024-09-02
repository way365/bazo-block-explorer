# Installation Guide
## &#9888; BAZO was moved to [Docker](https://github.com/way365/bazo-docker)
## Prerequisites
The explorer requires a working installation of Golang (version 1.9.2 was used to develop the application). The most recent version can be found on their website, along with installation instructions. A go folder has to be created in the user’s home directory after installation.
```sh
$ https://golang.org/dl/
```
Since the explorer requires a PSQL database, Postgres is also required.
```sh
$ sudo apt install postgresql postgresql-contrib
```
## Setting Up the Database
After installing Postgres, a database super user called ’postgres’ is automatically created. With that user, a new user has to be created, which will be used to manage the block explorer’s database. First, the psql shell needs to be opened using the postgres account.
```sh
$ sudo -i -u postgres
$ psql
```
The new database-user and his password has to be created. Additionally the database ’blockexplorerdb’ has to be created. Username and password of the new user can be chosen freely, however the database needs to be called ’blockexplorerdb’. The new user needs all privileges on the newly created database as well.
```sh
$ CREATE ROLE user1 WITH LOGIN PASSWORD 'userpasswd';
$ CREATE DATABASE blockexplorerdb;
$ GRANT ALL PRIVILEGES ON DATABASE blockexplorerdb TO user1;
$ \quit
```
If the terminal still displays the user ’postgres’ as the current user, it is now time to switch back to the computer’s regular user, as no further psql related actions are needed. If switching users requires a password for the ’postgres’ account, a restart of the terminal will also log in the regular user.
## Setting Up the Explorer
```sh
$ go get github.com/bazo-blockchain/bazo-block-explorer
```

From the explorer’s directory, it can now be started using the following arguments:
- DATA must either be "data" or "nodata", depending on whether the data retrieval mechanism should run. 
- PORT on which the block-explorer should be reachable on localhost. 
- USERNAME and PASSWORD are the values defined in the database setup step.
```sh
$ ./bazo-block-explorer data :8080 USERNAME PASSWORD
```
A running Bazo Miner application is required on the bootstrap server to load blocks from the blockchain and a running Bazo Client to send Config Transactions.
