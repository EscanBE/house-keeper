package constants

// Define constants in this file

//goland:noinspection GoSnakeCaseUsage
const (
	// TODO Sample of template

	APP_NAME    = "Go App Template"
	APP_DESC    = "This is description of Goo App Template app, please change it with your own description"
	BINARY_NAME = "goappnamed" // rename it, ends with 'd', eg: evmosd

	// Do not change bellow

	DEFAULT_HOME             = "." + BINARY_NAME
	DEFAULT_CONFIG_FILE_NAME = CONFIG_FILE_NAME_PART + "." + CONFIG_TYPE
	CONFIG_FILE_NAME_PART    = "config"
	CONFIG_TYPE              = "yaml"
)

//goland:noinspection GoSnakeCaseUsage
const (
	FLAG_HOME = "home"
)

//goland:noinspection GoSnakeCaseUsage
const (
	FILE_PERMISSION     = 0o600
	FILE_PERMISSION_STR = "600"
)
