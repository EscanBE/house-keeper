# Go Application Template
version: 4.3
### Features:
   I. Supported commands
   - init: generate configuration file
   - check: show & validate configuration file
   - version: check binary version
   - start: start application to process main business

   II. Supported business logic behaviors:
   - Config file (can be generated with `init` command), `read cmd/init.go` for more information
   - Worker, for parallel execution, if you don't need this, remove the `work` package with related source code
   - Report via Telegram
   - Database CRUD with transaction

   III. Supported coding features:
   - Flag
   - Init & Read configuration file
   - Logging
   - Telegram
   - Database (PostgreSQL) + Tx
   - Handle panic, execute exit func and panic again (search code `defer libapp.TryRecoverAndExecuteExitFunctionIfRecovered`) 

## How to use convert this app into your app
#### Let say you are going to use this template and convert it into your own telegram bot with app name `tbot` and source code located at `https://github.com/EscanBE/tbot.git`
1. Search & Replace package name `github.com/EscanBE/go-app-name` into your own
    - Eg: github.com/EscanBE/tbot
2. Replace directory name `cmd/goappnamed` into your app name with `d` as suffix
    - Eg: `cmd/tbotd` _( "tbot" + "d" )_
3. Replace the following constants in `constants/constants.go` with your app name
```go
const (
    APP_NAME = "Telegram Bot 1"
    APP_DESC = "Telegram Bot helps managing your system"
    BINARY_NAME = "tbotd" // `tbot` + `d`
)
```
4. Open `Makefile` and replace `goappnamed` with your own
5. Search everywhere again to make sure all the following patterns are replaced with your own, thus no longer exists
   - `github.com/EscanBE/go-app-name`
   - `goappnamed`
6. Update your configuration file template at `cmd/init.go`
7. Search pattern `// TODO Sample of template` and update with your own logic here
8. Now verify
```bash
go mod tidy
make build
./build/tbotd version
# To install binary, use `make install`
```

#### Notes
- Do not delete `// Legacy TODO xxx`, this is coding convention for TODO that keeps forever to append related code
- Always `defer libapp.TryRecoverAndExecuteExitFunctionIfRecovered` on every go-routines to release resource to prevent resources leak
- When want to exit app gracefully (eg: `os.Exit`), remember to call `libapp.ExecuteExitFunction()`
- Write the following text into README.md file of the new project
> This project follows Go Application Template version x.y