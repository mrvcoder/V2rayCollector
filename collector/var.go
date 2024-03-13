package collector

import (
	"flag"
	"github.com/valyala/fasthttp"
)

var (
	client       = &fasthttp.Client{}
	MaxMessages  = 100
	ConfigsNames = "@Vip_Security join us"
	Configs      = map[string]string{
		"ss":     "",
		"vmess":  "",
		"trojan": "",
		"vless":  "",
		"mixed":  "",
	}
	ConfigFileIds = map[string]int32{
		"ss":     0,
		"vmess":  0,
		"trojan": 0,
		"vless":  0,
		"mixed":  0,
	}
	Myregex = map[string]string{
		"ss":     `(?m)(...ss:|^ss:)\/\/.+?(%3A%40|#)`,
		"vmess":  `(?m)vmess:\/\/.+`,
		"trojan": `(?m)trojan:\/\/.+?(%3A%40|#)`,
		"vless":  `(?m)vless:\/\/.+?(%3A%40|#)`,
	}
	Sort = flag.Bool("sort", false, "sort from latest to oldest (default : false)")
)
