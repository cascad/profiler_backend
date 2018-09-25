package profiler

import (
	"os"
	"strconv"
)

type Config struct {
	Hostname     string
	MongoHost    string
	DBName       string
	LoadInterval int
}

func (c *Config) Init() {
	c.Hostname = os.Getenv("HOSTNAME")
	c.MongoHost = os.Getenv("MONGOHOST")
	c.DBName = os.Getenv("DBNAME")
	inverval, _ := strconv.Atoi(os.Getenv("LOAD_INTERVAL"))
	c.LoadInterval = inverval
}
