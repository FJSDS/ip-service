package main

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/oschwald/geoip2-golang"
)

type IPInfo struct {
	IP       string   `json:"ip"`       //ip地址
	Country  string   `json:"country"`  //国家
	Province string   `json:"province"` //省份
	City     string   `json:"city"`     //城市
	Location struct { //经纬度
		AccuracyRadius uint16  `json:"accuracy_radius"`
		Latitude       float64 `json:"latitude"`
		Longitude      float64 `json:"longitude"`
		MetroCode      uint    `json:"metro_code"`
		TimeZone       string  `json:"time_zone"`
	} `json:"location"`
}

var db *geoip2.Reader
var localNetworkNames = map[string]string{
	"zh-CN": "局域网",
	"en":    "local network",
}

func init() {
	var err error
	db, err = geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}

}
func main() {
	engine := gin.Default()
	engine.Any("/", HandleIPInfo)
	engine.Run(":25000")
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false

		default:
			return true
		}
	}
	return false
}

func HandleIPInfo(c *gin.Context) {
	language := c.Query("language")
	if language != "cn" && language != "en" {
		language = "cn"
	}
	if language == "cn" {
		language = "zh-CN"
	}
	format := c.Query("format")
	if format != "json" && format != "string" {
		format = "json"
	}
	ipStr := c.ClientIP()
	ip := net.ParseIP(ipStr)
	if !IsPublicIP(ip) {
		Success(c, &IPInfo{
			IP:   ip.String(),
			City: localNetworkNames[language],
		}, format)
		return
	}
	city, err := db.City(ip)
	if err != nil {
		c.JSON(200, gin.H{
			"success": "false",
			"reason":  ipStr + " is a invalid ip",
		})
		return
	}
	ipInfo := &IPInfo{}
	ipInfo.IP = ip.String()
	ip.IsLoopback()
	ipInfo.Country = city.Country.Names[language]
	if len(city.Subdivisions) > 0 {
		ipInfo.Province = city.Subdivisions[0].Names[language]
	} else {
		ipInfo.Province = city.City.Names[language]
	}
	ipInfo.City = city.City.Names[language]

	ipInfo.Location.AccuracyRadius = city.Location.AccuracyRadius
	ipInfo.Location.Latitude = city.Location.Latitude
	ipInfo.Location.Longitude = city.Location.Longitude
	ipInfo.Location.MetroCode = city.Location.MetroCode
	ipInfo.Location.TimeZone = city.Location.TimeZone
	Success(c, ipInfo, format)
}

func Success(c *gin.Context, ipInfo *IPInfo, format string) {
	if format == "string" {
		var area []string
		if ipInfo.Country != "" {
			area = append(area, ipInfo.Country)
		}
		if ipInfo.Province != "" {
			area = append(area, ipInfo.Province)
		}
		if ipInfo.City != "" {
			area = append(area, ipInfo.City)
		}
		c.JSON(200, gin.H{
			"success": "true",
			"ip_info": struct {
				IP   string `json:"ip"`
				Area string `json:"area"`
			}{
				IP:   ipInfo.IP,
				Area: strings.Join(area, " "),
			},
		})
	} else {
		c.JSON(200, gin.H{
			"success": "true",
			"ip_info": ipInfo,
		})
	}
}
