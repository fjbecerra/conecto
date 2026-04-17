package core

type Config interface{
	 Load(configPath string) Config
}