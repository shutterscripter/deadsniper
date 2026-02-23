package config

type ConfigHttp struct {
	URL        string
	Timeout    int
	Delay      float64
	Threads    int
	OutputType int
	Verbose    bool
	Help       bool
	Recursive  bool
	MaxDepth   int
}

var DefaultConfig ConfigHttp = ConfigHttp{
	URL:        "",
	Timeout:    10,
	Delay:      0.5,
	Threads:    1,
	OutputType: 2,
	Recursive:  true,
	MaxDepth:   3,
}
