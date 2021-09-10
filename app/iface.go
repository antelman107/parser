package app

type CsvApp interface {
	Run() error
}

type CsvAppBuilder interface {
	Build(k uint32, cfg CsvParserAppConfig) CsvApp
}
