package config

import (
	"os"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
)

type Config struct {
	RunEnv string `yaml:"run_env,omitempty"`
	Env    Env
	Db     Db `yaml:"db,omitempty"`
}
type Env struct {
	Db  EnvDb
	Rsa Rsa
}
type EnvDb struct {
	ServerAddress string
	ReplicaID     int32
}
type Rsa struct {
	PublicKey  string
	PrivateKey string
}
type Db struct {
	BatchSize *uint32 `yaml:"batch_size,omitempty"`
}

func loadEnvVar(envVar string) string {
	variable, exists := os.LookupEnv(envVar)
	if !exists {
		log.Fatalln("Env variable not found", envVar)
	}
	return variable
}

func LoadConfig(path string) Config {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	config.Env.Db.ServerAddress = loadEnvVar("dbServerAddress")
	replicaId, err := strconv.Atoi(loadEnvVar("dbReplicaID"))
	if err != nil {
		log.Fatalln("Failed to convert replica id")
	}
	config.Env.Db.ReplicaID = int32(replicaId)
	config.Env.Rsa.PublicKey = loadEnvVar("rsaPublicKey")
	if privateKey, exists := os.LookupEnv("rsaPrivateKey"); exists {
		config.Env.Rsa.PrivateKey = privateKey
	} else {
		if config.RunEnv != "Validation" {
			log.Fatalln("rsaPrivateKey key not found, trying to start in mode: ", config.RunEnv)
		}
		log.Println("Starting in validation mode")
	}
	log.Println("Config loaded")
	return config
}
