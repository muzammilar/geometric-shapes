/*
	The logger package extends the logrus.Logger with helpful reuseable functions
*/
package logs

/*
 * Structs
 */
type Configuration struct {
	Name           string // the name of the logger
	Level          string // the logging level of the logger
	Path           string // the path for storing the log
	Permissions    uint32 // the permissions of the file
	ConsoleLogging bool   // whether to add logging to the console or not
}

// NewConfiguration creates a new configuration from strings and add defaults
func NewConfiguration(name string, level string, path string, console bool) *Configuration {

	conf := &Configuration{
		Name:           name,
		Level:          level,
		Path:           path,
		Permissions:    defaultPermissions,
		ConsoleLogging: console,
	}
	// add defaults
	conf.AddDefaults()

	return conf
}

func (c *Configuration) AddDefaults() {
	// default name is "default"
	if c.Name == "" {
		c.Name = defaultLoggerName
	}
	// default level is info
	if c.Level == "" {
		c.Level = "info"
	}
	// default output path
	if c.Path == "" {
		c.Path = "/var/log/goapplication.log"
	}
	// default permissions
	if c.Permissions <= 0 {
		c.Permissions = defaultPermissions
	}
}
