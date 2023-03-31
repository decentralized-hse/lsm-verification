package config

type Config struct {
	RunEnv interface{} `yaml:"run_env"`
	Env    Env         `yaml:"env"`
	Db     Db          `yaml:"db"`
}
type Env struct {
	EnvDb Db  `yaml:"db"`
	Rsa   Rsa `yaml:"rsa"`
}
type EnvDb struct {
	ServerAddress interface{} `yaml:"server_address"`
	ReplicaID     interface{} `yaml:"replica_id"`
}
type Rsa struct {
	PublicKey  interface{} `yaml:"public_key"`
	PrivateKey interface{} `yaml:"private_key"`
}
type Db struct {
	BatchSize interface{} `yaml:"batch_size"`
}
