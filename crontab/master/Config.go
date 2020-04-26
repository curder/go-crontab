package master

import (
    "encoding/json"
    "io/ioutil"
)

// 配置结构体
type Config struct {
    APiPort         int      `json:"apiPort"`
    ApiReadTimeout  int      `json:"apiReadTimeout"`
    ApiWriteTimeout int      `json:"apiWriteTimeout"`
    EtcdEndPoints   []string `json:"etcdEndPoints"`
    EtcdDialTimeout int      `json:"etcdDialTimeout"`
    WebRoot         string   `json:"webRoot"`
}

var (
    GConfig *Config
)

func InitConfig(filename string) (err error) {
    var (
        content []byte
        config  Config
    )

    // 读取配置文件
    if content, err = ioutil.ReadFile(filename); err != nil {
        return
    }

    // JSON反序列化
    if err = json.Unmarshal(content, &config); err != nil {
        return
    }

    // 赋值单例
    GConfig = &config

    return
}
