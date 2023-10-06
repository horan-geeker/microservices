package config

func NewConfig[T any](conf *T) error {
	if err := ParseConfigFromEnvFile[T](conf); err != nil {
		return err
	}
	return nil
}
