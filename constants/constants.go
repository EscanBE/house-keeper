package constants

// Define constants in this file

//goland:noinspection GoSnakeCaseUsage
const (
	APP_DESC    = "House Keeper does multiple tasks"
	BINARY_NAME = "hkd" // rename it, ends with 'd', eg: evmosd
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
