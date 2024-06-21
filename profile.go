package profilego

import (
	"errors"

	"dario.cat/mergo"
	"github.com/grafana/pyroscope-go"
)

type Config struct {
	ApplicationName string            `json:"application_name"` // ApplicationName specifies the name of the application.
	ServerAddress   string            `json:"server_address"`   // ServerAddress specifies the address of the profiling server.
	Type            string            `json:"type"`             // Type specifies the type of profiler. Valid values are "pyroscope".
	Tags            map[string]string `json:"tags"`             // Tags specifies the tags to be added to the profiler.
}

var defaultConfig = Config{
	ApplicationName: "my-app",
	ServerAddress:   "127.0.0.1:4040",
	Type:            "pyroscope",
	Tags:            map[string]string{},
}

func Init(config Config, additionalAttrs ...any) error {

	err := mergo.Merge(&defaultConfig, config, mergo.WithOverride)
	if err != nil {
		return err
	}

	if defaultConfig.ApplicationName == "" {
		return errors.New("application name not provided (ApplicationName)")
	}

	if defaultConfig.ServerAddress == "" {
		return errors.New("server address name not provided (ServerAddress)")
	}

	if defaultConfig.Type == "pyroscope" {
		_, err := initPyroscope(defaultConfig)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("profiler type not provided (Type)")

}

func initPyroscope(config Config) (*pyroscope.Profiler, error) {
	profiler, err := pyroscope.Start(pyroscope.Config{
		Logger:          ProfilingLogger{},
		ApplicationName: config.ApplicationName,

		// replace this with the address of pyroscope server
		ServerAddress: config.ServerAddress,

		// you can provide static tags via a map:
		Tags: config.Tags,

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	if err != nil {
		return nil, err
	}

	return profiler, nil
}
