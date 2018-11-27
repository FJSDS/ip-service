# ip-service
提供ip地址地理位置查询

## 初始化
```golang
go get -u github.com/FJSDS/ip-service
```

## 使用说明
目前只有一个http接口  __"/"__ 返回json结构体，可选query参数
- language=[cn,en],默认cn,
- format=[json,string],默认json
### 成功返回如下结构体
```golang
{
  Success bool `json:"success"`
  IPInfo struct {
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
  } `json:"ip_info"`
}
```
#### 例子
**request:**
```
http://ip:25000/
http://ip:25000/?language=cn
```
**response**
```json
{
  "ip_info": {
    "ip": "118.113.146.10",
    "country": "中国",
    "province": "四川省",
    "city": "成都",
    "location": {
      "accuracy_radius": 200,
      "latitude": 30.6667,
      "longitude": 104.0667,
      "metro_code": 0,
      "time_zone": "Asia/Shanghai"
    }
  },
  "success": "true"
}
```

**request:**
```
http://ip:25000/?format=string
http://ip:25000/?language=cn&format=string
```
**response**
```json
{
  "ip_info": {
    "ip": "118.113.146.10",
    "area": "中国 四川省 成都"
  },
  "success": "true"
}
```

### 失败时
success=false
reason = [failed reason]
