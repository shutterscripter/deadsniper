package config

type ConfigHttp struct {
	URL        string
	Timeout    int
	Delay      float64
	Threads    int
	OutputType int // 6: json, 2: csv, 4: xml, 8: text

}

var DefaultConfig ConfigHttp  =   ConfigHttp{
	URL: "",
	Timeout: 10,
	Delay: 0.5,
	Threads: 1,
	OutputType: 2,
}
