## House Keeper
#### A tool to manage Escan
> go install -v github.com/EscanBE/house-keeper/cmd/hkd@v0.3.0

### Commands:

#### Listing files
> hkd files list --help

> hkd files list --working-directory '/tmp' --order-by date --contains '0' --skip 1

> hkd files list --working-directory '/tmp/backup-db' --order-by name --desc --contains '.dump' --skip 3 --silent --delete

###### This project uses Go Application Template v4.3 (by Escan)