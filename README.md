## House Keeper
#### A tool to manage Escan
> go install -v github.com/EscanBE/house-keeper/cmd/hkd@v0.3.0

### Commands:

#### Listing files
> hkd files list --help

> hkd files list --working-directory '/tmp' --order-by date --contains '0' --skip 1

> hkd files list --working-directory '/tmp/backup-db' --order-by name --desc --contains '.dump' --skip 3 --silent --delete

#### Perform database backup:
> hkd db backup --help

> PGPASSWORD=1234567 hkd db backup --working-directory /mnt/md0/backup --dbname my_db_name --username my_user_name

> hkd db backup --working-directory /mnt/md0/backup --output-file db-backup-2023-01-02.dump --host localhost --port 5432 --dbname postgres --username postgres --schema public --password-file ~/password.txt

Notes:
- Current only support PostgreSQL
- Either environment variable PGPASSWORD or flag --password-file is required (priority flag)
- Rely on pg_dump command to perform backup action for PostgreSQL, it actually set environment variable PGPASSWORD and run the following command: pg_dump --host=(host) --port=(port) --schema=(schema) -Fc --username=(username) --file=(output file) (dbname)

###### This project uses Go Application Template v4.3 (by Escan)