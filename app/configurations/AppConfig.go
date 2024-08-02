package configurations

import (
	"log"
	"sync"
	"sync/atomic"

	"github.com/spf13/viper"
)

type ConfigArgs struct {
	Name string
	Path []string
}

var defaultConfigName = "appConfigs"
var (
	defaultConfigPath = []string{"./configs", "${HOME}/configs"}
	configs           *viper.Viper
	done              int32
	m                 sync.Mutex
)

func Setup(configArgs *ConfigArgs) {
	if atomic.LoadInt32(&done) == 1 {
		return
	}
	m.Lock()
	defer m.Unlock()
	if done == 0 {
		defer atomic.StoreInt32(&done, 1)
		doSetup(configArgs)
	}

}

func doSetup(args *ConfigArgs) {
	log.Println("doing set up of configurations")
	configs = viper.New()
	name := defaultConfigName
	path := defaultConfigPath

	if args != nil {
		if args.Name != "" {
			name = args.Name
		}
		if args.Path != nil && len(args.Path) > 0 {
			path = args.Path
		}
	}
	configs.SetConfigName(name)
	configs.SetConfigType("toml")
	for _, val := range path {
		configs.AddConfigPath(val)
	}

	if err := configs.ReadInConfig(); err != nil {
		log.Printf("Error reading in config file %v", err)
	}
	log.Printf("Application is using config %s", configs.ConfigFileUsed())
}

func Get(key string) interface{} {
	if atomic.LoadInt32(&done) == 0 {
		log.Panicln("Configurations have not been initialized till now")
	}
	return configs.Get(key)
}
