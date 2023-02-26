package option

type Option struct {
	Cache Cache `mapstructure:"cache"`
	App   App   `mapstructure:"app"`
}

type Cache struct {
	Address *string `mapstructure:"address"`
}

type App struct {
	Port string `mapstructure:"port"`
}
