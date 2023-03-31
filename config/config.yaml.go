package config
type Config struct {
	RunEnv interface{} `yaml:"run_env"`
	Env Env `yaml:"env"`
}
type Env struct {
	ServerAddress interface{} `yaml:"server_address"`
	ReplicaID interface{} `yaml:"replica_id"`
	BatchSize interface{} `yaml:"batch_size"`
}
