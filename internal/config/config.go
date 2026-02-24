package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type FuncCounters struct {
	ConfigLoadingFuncCounter int
	LoggerFuncCounter        int
	DatabaseFuncCounter      int
	ServiceFuncCounter       int
	HandlerFuncCounter       int
}

type ServerConfig struct {
	AppName  string `mapstructure:"app_name"`
	Env      string `mapstructure:"env"`
	Debug    bool   `mapstructure:"debug"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	TimeZone string `mapstructure:"time_zone"`
}

type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	AdminUser string `mapstructure:"admin_user"`
	AdminPassword string `mapstructure:"admin_password"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Debug    bool   `mapstructure:"debug"`
}

type JWTConfig struct {
	Secret   string `mapstructure:"secret"` //签名密钥
	Issuer   string `mapstructure:"issuer"` // 签发者
	Audience string `mapstructure:"audience"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	// 配置文件名
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("../../") // 根目录读取 config.yaml
	fmt.Println("✓ Loading config file done")

	// 支持环境变量覆盖，例如：SERVER_APP_NAME
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)
	fmt.Println("✓ Setting default configuration done")

	// 读取 YAML
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Error: ⚠ config.yaml not found, using default values and environment variables")
	} else {
		fmt.Printf("✓ Using config file: %s, not using default configuration\n", v.ConfigFileUsed())
	}

	// 确保 jwt.secret 存在且强度足够，否则生成并写入配置文件
	if err := EnsureJWTSecret(v); err != nil {
		return nil, fmt.Errorf("Error: can't generate and write in jwt.secret: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Error: parsing config failed: %w", err)
	}
	fmt.Println("✓ Unmarshalling config file done")

	// 校验
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	fmt.Println("✓ Validating config file done")

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.app_name", "secure_file_box")
	v.SetDefault("server.env", "development")
	v.SetDefault("server.debug", true)
	v.SetDefault("server.host", "127.0.0.1")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.time_zone", "Asia/Shanghai")

	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "0827")      // you own password, change it
	v.SetDefault("database.name", "secure_file_box") //create database with name "secure_file_box"
	v.SetDefault("database.debug", false)

	v.SetDefault("jwt.secret", "PLEASE_CHANGE_ME_32_CHARS_MINIMUM")
	v.SetDefault("jwt.issuer", "secure_file_box")
	v.SetDefault("jwt.audience", "secure_users")
}

func validateConfig(cfg *Config) error {
	if len(cfg.JWT.Secret) < 32 {
		return fmt.Errorf("Error: jwt.secret must be greater than 32 fugures")
	}
	if cfg.Database.Name == "" {
		return fmt.Errorf("Error: database.name can't be empty")
	}
	return nil
}

// GenerateJWTSecret 生成一个高强度的随机字符串，使用 base64 URL-safe 编码。
// 参数 nBytes 指定随机字节数，推荐至少 32（256 bits）。
// base64 encoding
func GenerateJWTSecret(nBytes int) (string, error) {
	if nBytes <= 0 {
		nBytes = 32
	}
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// 使用 RawURLEncoding 去掉 padding，使字符串在 URL/headers 中更安全
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// EnsureJWTSecret 检查 viper 中的 jwt.secret，若缺失或强度不足则生成并写回配置文件。
func EnsureJWTSecret(v *viper.Viper) error {
	cur := v.GetString("jwt.secret")
	if len(cur) >= 32 && cur != "PLEASE_CHANGE_ME_32_CHARS_MINIMUM" {
		return nil
	}

	secret, err := GenerateJWTSecret(32) // 32 bytes = 256 bits, 可写33，34...等更高
	if err != nil {
		return err
	}

	v.Set("jwt.secret", secret)

	// 如果已经加载了配置文件，写回同一文件；否则在运行目录写入 config.yaml
	if v.ConfigFileUsed() != "" {
		return v.WriteConfig() // 覆盖原有配置文件
	}
	return v.WriteConfigAs("config.yaml")
}

// 函数使用计数器
func NewFuncCounters() *FuncCounters {
	return &FuncCounters{
		ConfigLoadingFuncCounter: 0,
		LoggerFuncCounter:        0,
		DatabaseFuncCounter:      0,
		ServiceFuncCounter:       0,
		HandlerFuncCounter:       0,
	}
}
