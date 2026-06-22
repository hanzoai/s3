package util

import (
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/hanzoai/s3/s3/glog"
)

var (
	ConfigurationFileDirectory DirectoryValueType
	loadSecurityConfigOnce     sync.Once
)

type DirectoryValueType string

func (s *DirectoryValueType) Set(value string) error {
	*s = DirectoryValueType(value)
	return nil
}
func (s *DirectoryValueType) String() string {
	return string(*s)
}

type Configuration interface {
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetStringSlice(key string) []string
	SetDefault(key string, value interface{})
}

func LoadSecurityConfiguration() {
	loadSecurityConfigOnce.Do(func() {
		LoadConfiguration("security", false)
	})
}

func LoadConfiguration(configFileName string, required bool) (loaded bool) {

	// find a filer store
	viper.SetConfigName(configFileName)                                   // name of config file (without extension)
	viper.AddConfigPath(ResolvePath(ConfigurationFileDirectory.String())) // path to look for the config file in
	viper.AddConfigPath(".")                                              // optionally look for config in the working directory
	viper.AddConfigPath("$HOME/.s3")                               // call multiple times to add many search paths
	viper.AddConfigPath("/usr/local/etc/hanzo/")                      // search path for bsd-style config directory in
	viper.AddConfigPath("/etc/hanzo/")                                // path to look for the config file in

	if err := viper.MergeInConfig(); err != nil { // Handle errors reading the config file
		if strings.Contains(err.Error(), "Not Found") {
			glog.V(1).Infof("Reading %s: %v", viper.ConfigFileUsed(), err)
		} else {
			// If the config is required, fail immediately
			if required {
				glog.Fatalf("Reading %s: %v", viper.ConfigFileUsed(), err)
			}
			// If the config is optional, log a warning but don't crash
			glog.Warningf("Reading %s: %v. Skipping optional configuration.", viper.ConfigFileUsed(), err)
		}
		if required {
			glog.Fatalf("Failed to load %s.toml file from current directory, or $HOME/.s3/, or /etc/hanzo/"+
				"\n\nPlease use this command to generate the default %s.toml file\n"+
				"    s3 scaffold -config=%s -output=.\n\n\n",
				configFileName, configFileName, configFileName)
		}
		return false
	}
	glog.V(1).Infof("Reading %s.toml from %s", configFileName, viper.ConfigFileUsed())

	return true
}

type ViperProxy struct {
	*viper.Viper
	sync.Mutex
}

var (
	vp = &ViperProxy{}
)

func (vp *ViperProxy) SetDefault(key string, value interface{}) {
	vp.Lock()
	defer vp.Unlock()
	vp.Viper.SetDefault(key, value)
}

func (vp *ViperProxy) GetString(key string) string {
	vp.Lock()
	defer vp.Unlock()
	return vp.Viper.GetString(key)
}

func (vp *ViperProxy) GetBool(key string) bool {
	vp.Lock()
	defer vp.Unlock()
	return vp.Viper.GetBool(key)
}

func (vp *ViperProxy) GetInt(key string) int {
	vp.Lock()
	defer vp.Unlock()
	return vp.Viper.GetInt(key)
}

func (vp *ViperProxy) GetStringSlice(key string) []string {
	vp.Lock()
	defer vp.Unlock()
	return vp.Viper.GetStringSlice(key)
}

func GetViper() *ViperProxy {
	vp.Lock()
	defer vp.Unlock()

	if vp.Viper == nil {
		vp.Viper = viper.GetViper()
		vp.AutomaticEnv()
		vp.SetEnvPrefix("s3")
		vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}

	return vp
}
