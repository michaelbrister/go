package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var myConfig *Config

type Config struct {
	EngineVersion string
	ClusterIds    []string
}

func init() {

}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}

	viper.Unmarshal(&myConfig)

	pflag.StringSliceVar(&myConfig.ClusterIds, "cluster", myConfig.ClusterIds, "CSV list of clusters")

	pflag.Parse()

	// clusterIds = viper.GetStringSlice("clusterIds")

	for _, cluster := range myConfig.ClusterIds {
		fmt.Println(cluster)
	}
}
