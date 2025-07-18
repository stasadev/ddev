---
search:
  boost: 2 
---
# Database Management

DDEV provides lots of flexibility for managing your databases between your local, staging and production environments. You may commonly use the [`ddev import-db`](../usage/commands.md#import-db) and [`ddev export-db`](../usage/commands.md#export-db) commands, but there are plenty of other adaptable ways to work with your databases.

If your project does _not_ require a database, you can exclude it with the the [`omit_containers` configuration option](../configuration/config.md#omit_containers).

!!!tip
    You can run `ddev [command] --help` for more info on many of the topics below.

## Database Imports

Import a database with one command, from one of the following file formats:  
**`.sql`, `.sql.gz`, `.mysql`, `.mysql.gz`, `.tar`, `.tar.gz`, and `.zip`**.

Here’s an example of a database import using DDEV:

```bash
ddev import-db --file=dumpfile.sql.gz
```

You can also:

* Use [`ddev mysql`](../usage/commands.md#mysql) or `ddev psql` or the `mysql` and `psql` commands inside the `web` and `db` containers.
* Use a [database client](#database-clients) or [database GUI](#database-guis) to import and browse data.

## Database Backends and Defaults

You can use a [variety of different database types](../extend/database-types.md#database-server-types), including MariaDB (5.5–10.8, 11.4, 11.8), MySQL (5.5–8.0, 8.4), and PostgreSQL (9–17). If you want to _change_ database type, you need to export your database, run [`ddev delete`](../usage/commands.md#delete) to remove the project (and its existing database), change to a new database type, run [`ddev start`](../usage/commands.md#start) again, and [import your data](../usage/commands.md#import-db).

(For very old database types see [Using DDEV to spin up a legacy PHP application](https://ddev.com/blog/legacy-projects-with-unsupported-php-and-mysql-using-ddev/).)

DDEV creates a default database named `db` and default permissions for the `db` user with password `db`, and it’s on the (inside Docker) hostname `db`.

## Extra Databases

You can easily create and populate additional databases. For example, `ddev import-db --database=backend --file=backend.sql.gz` will create the database named `backend` with permissions for that same `db` user and import from the `backend.sql.gz` dumpfile.

You can export in the same way: `ddev export-db -f mysite.sql.gz` will export your default database (`db`). `ddev export-db --database=backend -f backend-export.sql.gz` will dump the database named `backend`.

## Snapshots

Snapshots let you easily save the entire status of all of your databases, which can be great when you’re working incrementally on migrations or updates and want to save state so you can start right back where you were.

Snapshots can be named for easier reference later on. For example, [`ddev snapshot --name=two-dbs`](../usage/commands.md#snapshot) would make a snapshot named `two-dbs` in the `.ddev/db_snapshots` directory. It includes the entire state of the db server, so in the case of our two databases above, both databases and the system level `mysql` or `postgres` database will all be snapshotted. Then if you want to delete everything with `ddev delete -O` (omitting the snapshot since we have one already), and then [`ddev start`](../usage/commands.md#start) again, we can `ddev snapshot restore two-dbs` and we’ll be right back where we were.

Use the [`ddev snapshot restore`](../usage/commands.md#snapshot-restore) command to interactively choose among snapshots, or append `--latest` to restore the most recent snapshot: `ddev snapshot restore --latest`.

Snapshots are stored as simple gzipped files in the project's `.ddev/db_snapshots` directory, and any or all snapshots can be removed with the `ddev snapshot --cleanup` command or by manually deleting the files when you want to save disk space or have no further use for them.

## Database Clients

The `ddev mysql` and `ddev psql` commands give you direct access to the `mysql` and `psql` clients in the database container, which can be useful for quickly running commands while you work. You might run `ddev mysql` to use interactive commands like `DROP DATABASE backend;` or `SHOW TABLES;`, or do things like `echo "SHOW TABLES;" | ddev mysql` or `ddev mysql -udb -pdb` to run with `db` user privileges.

The `web` and `db` containers are each ready with MySQL/PostgreSQL clients, so you can [`ddev ssh`](../usage/commands.md#ssh) or `ddev ssh -s db` and use `mysql` or `psql`.

## `mysqldump` and `pg_dump`

The `web` and `db` containers come with `mysqldump`. You could run [`ddev ssh`](../usage/commands.md#ssh) to enter the web container, for example, then `mkdir /var/www/html/.tarballs` and run `mysqldump db >/var/www/html/.tarballs/db.sql` or run `mysqldump db | gzip >/var/www/html/.tarballs/db.sql.gz` to create database dumps. Because `/var/www/html` is mounted into the container from your project root, the `.tarballs` directory will also show up in the root of the project on your host machine.

The PostgreSQL database container includes normal `pg` commands like `pg_dump`.

## Database GUIs

If you’d like to use a GUI database client, you’ll need the right connection details and there may even be a command to launch it for you:

* phpMyAdmin, formerly built into DDEV core, can be installed by running `ddev add-on get ddev/ddev-phpmyadmin`.
* Adminer can be installed with `ddev add-on get ddev/ddev-adminer`
* The [`ddev describe`](../usage/commands.md#describe) command displays the `Host:` details you’ll need to connect to the `db` container externally, for example if you're using an on-host database browser like SequelAce.
* macOS users can use [`ddev sequelace`](../usage/commands.md#sequelace) to launch the free [Sequel Ace](https://sequel-ace.com/) database browser, [`ddev tableplus`](../usage/commands.md#tableplus) to launch [TablePlus](https://tableplus.com), [`ddev querious`](../usage/commands.md#querious) to launch [Querious](https://www.araelium.com/querious), [`ddev dbeaver`](../usage/commands.md#dbeaver) to launch [DBeaver](https://dbeaver.io/), and the obsolete Sequel Pro is also supported with `ddev sequelpro`. (Each must be installed for the command to exist.)
* WSL2 and Linux users can use [`ddev dbeaver`](../usage/commands.md#dbeaver) to launch [DBeaver](https://dbeaver.io/). (Must be installed for the command to exist.)
* PhpStorm (and all JetBrains tools) have a nice database browser. (If you use the [DDEV Integration plugin](https://plugins.jetbrains.com/plugin/18813-ddev-integration) this is all done for you.)
    * Choose a static [`host_db_port`](../configuration/config.md#host_db_port) setting for your project. For example `host_db_port: 59002` (each project’s database port should be different if you’re running more than one project at a time). Use [`ddev start`](../usage/commands.md#start) for it to take effect.
    * Use the “database” tool to create a source from “localhost”, with the proper type `mysql` or `postgresql` and the port you chose, username `db` + password `db`.
    * Explore away!
* VS Code or any of its forks (Cursor, Windsurf, etc.) has [DevDb](https://marketplace.visualstudio.com/items?itemName=damms005.devdb), an extension that lets you access your database directly inside the editor.
* There’s a sample custom command that will run the free MySQL Workbench on macOS, Windows or Linux. To use it, run:
    * `cp ~/.ddev/commands/host/mysqlworkbench.example ~/.ddev/commands/host/mysqlworkbench`
    * `ddev mysqlworkbench`

## Database Query Examples

You can query, update, or alter databases in DDEV like you would do on a regular server, using `ddev mysql`,  `ddev mariadb`, or `ddev psql` or using those command-line tools inside the web or DB containers. Some examples are given below.

* Create a new empty database named `newdatabase`:

    ```bash
    ddev mysql -e 'CREATE DATABASE newdatabase; GRANT ALL ON newdatabase.* TO "db"@"%";'
    ```

* Show tables whose name begins with `node` :

    ```bash
    ddev mysql -e 'SHOW TABLES LIKE "node%";'
    ```

* Use `ddev mysql` or `ddev mariadb` and issue interactive queries:

    ```bash
    ddev mariadb
  
    Reading table information for completion of table and column names
    You can turn off this feature to get a quicker startup with -A
    
    Welcome to the MariaDB monitor.  Commands end with ; or \g.
    Your MariaDB connection id is 28
    Server version: 10.11.10-MariaDB-ubu2204-log mariadb.org binary distribution
    
    Copyright (c) 2000, 2018, Oracle, MariaDB Corporation Ab and others.
    
    Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
    
    MariaDB [db]> SELECT * FROM node WHERE type="article";
    +-----+------+---------+--------------------------------------+----------+
    | nid | vid  | type    | uuid                                 | langcode |
    +-----+------+---------+--------------------------------------+----------+
    |  11 |   56 | article | 8af917ea-b150-4006-aeb9-877b53ebf289 | en       |
    |  12 |   54 | article | 27a53763-9fd8-4813-853b-b476f0a73849 | en       |
    +-----+------+---------+--------------------------------------+----------+
    2 rows in set (0.007 sec) 
    ```
