package constants

// Define constants in this file

//goland:noinspection GoSnakeCaseUsage
const (
	APP_DESC    = "House Keeper does multiple tasks"
	BINARY_NAME = "hkd" // rename it, ends with 'd', eg: evmosd
)

//goland:noinspection GoSnakeCaseUsage
const (
	FLAG_WORKING_DIR = "working-directory"
	FLAG_ORDER_BY    = "order-by"
	FLAG_CONTAINS    = "contains"
	FLAG_DESC        = "desc"
	FLAG_SILENT      = "silent"
	FLAG_SKIP        = "skip"
	FLAG_DELETE      = "delete"

	FLAG_HOST          = "host"
	FLAG_PORT          = "port"
	FLAG_DB_NAME       = "dbname"
	FLAG_USER_NAME     = "username"
	FLAG_SUPER_USER    = "superuser"
	FLAG_SCHEMA        = "schema"
	FLAG_OUTPUT_FILE   = "output-file"
	FLAG_LOG_FILE      = "log-file"
	FLAG_EXCLUDE_DIRS  = "exclude-dirs"
	FLAG_PASSWORD_FILE = "password-file"
	FLAG_NO_PASSWORD   = "no-password"
	FLAG_TOOL_FILE     = "tool-file"
	FLAG_TOOL_OPTIONS  = "tool-options"
	FLAG_DATA_ONLY     = "data-only"
	FLAG_NO_PUB_SUB    = "no-pubsub"

	FLAG_REMOTE_TO_LOCAL = "remote-to-local"
	FLAG_LOCAL_TO_REMOTE = "local-to-remote"
	FLAG_LOCAL_TO_LOCAL  = "local-to-local"

	FLAG_SSHPASS_PASSPHRASE = "passphrase"

	FLAG_CACHE_AND_TRUST = "cache-and-trust"
)

//goland:noinspection GoSnakeCaseUsage
const (
	FILE_PERMISSION_700     = 0o700
	FILE_PERMISSION_700_STR = "700"

	FILE_PERMISSION_600     = 0o600
	FILE_PERMISSION_600_STR = "600"

	FILE_PERMISSION_400     = 0o400
	FILE_PERMISSION_400_STR = "400"
)

//goland:noinspection GoSnakeCaseUsage
const (
	ENV_PG_PASSWORD    = "PGPASSWORD"
	ENV_RSYNC_PASSWORD = "RSYNC_PASSWORD"
	ENV_SSHPASS        = "SSHPASS"
)
