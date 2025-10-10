package nodes

const (
	methodStartNatTraversal = "nodes.start_nat_traversal"
)

type Config struct {
	LogPings bool `yaml:"log_pings"`
}

var defaultConfig = Config{}
