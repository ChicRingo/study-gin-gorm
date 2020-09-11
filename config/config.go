package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局变量指针，用来保存程序的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	*ServerConfig `mapstructure:"server"`
	*LogConfig    `mapstructure:"log"`
	*MySQLConfig  `mapstructure:"mysql"`
}
type ServerConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapstructure:"port"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

// log配置信息
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// mysql配置信息
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Charset      string `mapstructure:"charset"`
	Loc          string `mapstructure:"loc"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

func Init() (err error) {
	// 方式1：直接指定配置文件的路径（相对路径或者绝对路径）
	// main函数目录的相对文件位置，注意使用此函数，不要使用下面的方式2，会覆盖你的设置，建议二选一
	viper.SetConfigFile("config/config.yaml")

	// 方式2：指定配置文件名和配置文件的位置，viper自行查找可用的配置文件
	//viper.AddConfigPath("./config/")                    // 指定查找配置文件的路径,可配置多个
	//viper.SetConfigName("config")               // 指定配置文件(不需要带后缀)

	// 方式3：远程获取配置文件的指定文件类型
	//因为字节流中没有文件扩展名，所以支持的扩展名是“ json”，“ toml”，“ yaml”，“ yml”，“ properties”，“ props”，“ prop”，“ env”，“ dotenv”
	//viper.SetConfigType("yaml") // 指定配置文件(专用于从远程获取配置信息时指定配置文件类型的)

	// 读取配置信息
	if err = viper.ReadInConfig(); err != nil {
		// 读取配置文件失败
		fmt.Printf("viper.ReadInConfig() failed, err:%v\n", err)
		return
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
		return
	}
	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal() failed, err:%v\n", err)
			return
		}
	})
	return
}
