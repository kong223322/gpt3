package config

type Config struct {
	Hosts struct {
		Etcd struct {
			Address          []string
			RegisterTTL      int
			RegisterInterval int
		}
		Mysql struct {
			Host      string
			User      string
			Pass      string
			DBName    string
			IfShowSql bool
			IfSyncDB  bool
		}
		Redis struct {
			Address []string
			Pass    string
			DB      int
		}
	}
	Project     string
	ServiceName string
	Env         string
	Version     string
}
