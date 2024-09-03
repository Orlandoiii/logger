package logger

type ConfigLogger struct {
	Console           bool
	BeutifyConsoleLog bool
	File              bool
	Ruta              string
	MinLevel          string
	RotationMaxSizeMB int
	MaxAgeDay         int
	MaxBackups        int
	Compress          bool
}
