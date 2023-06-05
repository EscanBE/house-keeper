## House Keeper
#### A tool to manage Escan
> go install -v github.com/EscanBE/house-keeper/cmd/hkd@v0.6.0

### Commands:

#### Listing files
> hkd files list --help

> hkd files list --working-directory '/tmp' --order-by date --contains '0' --skip 1

> hkd files list --working-directory '/tmp/backup-db' --order-by name --desc --contains '.dump' --skip 3 --silent --delete

#### Sync files:
> hkd files rsync --help

> RSYNC_PASSWORD=1234567 hkd files rsync /var/log/nginx/access.log backup@192.168.0.2:/mnt/md0/backup/nginx-logs --local-to-remote

> hkd files rsync /var/log/nginx/access.log backup@192.168.0.2:/mnt/md0/backup/nginx-logs --local-to-remote --password-file ~/password.txt

> SSHPASS=1234567 hkd files rsync /var/log/nginx/access.log backup-server:/mnt/md0/backup/nginx-logs --local-to-remote --passphrase

Notes:
- This use rsync
- When either source or destination is remote machine:
  - Either environment variable RSYNC_PASSWORD or ENV_SSHPASS or flag --password-file is required (priority flag)
  - Environment variables RSYNC_PASSWORD and ENV_SSHPASS are treated similar thus either needed. If both provided, must be identical
  - You must connect to that remote server at least one time before to perform host key verification (one time action) because the transfer will be performed via ssh.

#### File checksum:
> hkd files checksum --help

> hkd files checksum /tmp/test.txt

#### Perform PostgreSQL DB backup:
> hkd db pg_dump --help

> PGPASSWORD=1234567 hkd db pg_dump --working-directory /mnt/md0/backup --dbname my_db_name --username my_user_name

> hkd db pg_dump --working-directory /mnt/md0/backup --output-file db-2023-01-02.dump --host localhost --port 5432 --dbname postgres --username postgres --schema public --password-file ~/password.txt

Notes:
- Either environment variable PGPASSWORD or flag --password-file is required (priority flag)
- Rely on pg_dump command to perform backup action for PostgreSQL, it actually set environment variable PGPASSWORD and then call pg_dump

#### Perform PostgreSQL DB restore:
> hkd db pg_restore --help

> PGPASSWORD=1234567 hkd db pg_restore db-2023-01-02.dump --superuser postgres --dbname example

> hkd db pg_restore db-2023-01-02.dump --host localhost --port 5432 --dbname example --username postgres --superuser postgres --password-file ~/password.txt

Notes:
- Either environment variable PGPASSWORD or flag --password-file is required (priority flag)
- Rely on pg_restore command to perform backup action for PostgreSQL, it actually set environment variable PGPASSWORD and then call pg_restore

#### Checking tools used by house-keeper
> hkd verify-tools

###### This project uses Go Application Template v4.3 (by Escan)