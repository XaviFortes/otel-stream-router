package routerprocessor

import (
	"errors"
)

// Config define la configuración de nuestro procesador custom.
type Config struct {
	// DefaultStream es el fallback por si falla todo lo demás
	DefaultStream string `mapstructure:"default_stream"`
}

// Validate es una función obligatoria de OTel para comprobar que la
// configuración que el usuario ha puesto en el YAML tiene sentido antes de arrancar.
func (c *Config) Validate() error {
	if c.DefaultStream == "" {
		return errors.New("default_stream no puede estar vacío")
	}
	return nil
}
