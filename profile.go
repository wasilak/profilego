package profilego

import (
	"errors"

	"dario.cat/mergo"
	"github.com/grafana/pyroscope-go"
)

type ProfileGoConfig struct {
	ApplicationName string
	ServerAddress   string
	Type            string
	Tags            map[string]string
}

var defaultConfig = ProfileGoConfig{
	ApplicationName: "",
	ServerAddress:   "",
	Type:            "pyroscope",
	Tags:            map[string]string{},
}

func Init(config ProfileGoConfig, additionalAttrs ...any) error {

	err := mergo.Merge(&defaultConfig, config, mergo.WithOverride)
	if err != nil {
		return err
	}

	if defaultConfig.ApplicationName == "" {
		return errors.New("application name not provided (ApplicationName)")
	}

	if defaultConfig.ServerAddress == "" {
		return errors.New("srver address name not provided (ServerAddress)")
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

func initPyroscope(config ProfileGoConfig) (*pyroscope.Profiler, error) {
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
