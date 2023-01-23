package utils

type DbModeEnum uint8

const (
	DB_MODE_TEST DbModeEnum = iota
	DB_MODE_PROD
)
