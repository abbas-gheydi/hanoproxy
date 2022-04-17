package confs

import (
	"fmt"

	"github.com/Abbas-gheydi/hanoproxy/dns/server"
	"github.com/Abbas-gheydi/hanoproxy/model"
	"github.com/spf13/viper"
)

var (
	HaNoProxy     *Configs
	Configuration Configs
)

type Configs struct {
	GlobalOptions server.ServerConfs
	DnsRecords    []model.Record
}

func Read() {
	viper.SetConfigName("HANoProxy")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/hanoproxy/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	err = viper.Unmarshal(&HaNoProxy)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	*server.Configs = HaNoProxy.GlobalOptions

}
