package constants

// Define constants in this file

//goland:noinspection GoSnakeCaseUsage
const (
	APP_NAME    = "House Keeper"
	APP_DESC    = "House Keeper does multiple tasks"
	BINARY_NAME = "hkd" // rename it, ends with 'd', eg: evmosd

	// Do not change bellow

	DEFAULT_HOME             = "." + BINARY_NAME
	DEFAULT_CONFIG_FILE_NAME = CONFIG_FILE_NAME_PART + "." + CONFIG_TYPE
	CONFIG_FILE_NAME_PART    = "config"
	CONFIG_TYPE              = "yaml"
)

//goland:noinspection GoSnakeCaseUsage
const (
	FLAG_HOME        = "home"
	FLAG_WORKING_DIR = "working-directory"
	FLAG_ORDER_BY    = "order-by"
	FLAG_CONTAINS    = "contains"
	FLAG_DESC        = "desc"
	FLAG_SILENT      = "silent"
	FLAG_SKIP        = "skip"
)

//goland:noinspection GoSnakeCaseUsage
const (
	FILE_PERMISSION     = 0o600
	FILE_PERMISSION_STR = "600"
)
