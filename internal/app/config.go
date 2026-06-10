package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ai-mate-client-gateway 的完整运行配置。
type Config struct {
	// 服务名。
	AppName string `mapstructure:"appName"`
	// 运行环境。
	Env string `mapstructure:"env"`
	// HTTP 服务配置。
	HTTP HTTPConfig `mapstructure:"http"`
	// 日志输出配置。
	Log LogConfig `mapstructure:"log"`
	// ai-mate-core gRPC 调用配置。
	Core CoreConfig `mapstructure:"core"`
	// OpenTelemetry 链路追踪配置。
	Observ ObservConfig `mapstructure:"observability"`
}

// HTTP 服务监听和超时参数。
type HTTPConfig struct {
	// HTTP 服务监听地址。
	Addr string `mapstructure:"addr"`
	// 请求读取超时时间。
	ReadTimeout time.Duration `mapstructure:"readTimeout"`
	// 响应写入超时时间。
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	// 空闲连接保留时间。
	IdleTimeout time.Duration `mapstructure:"idleTimeout"`
	// 优雅关闭超时时间。
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

// 日志输出格式和等级。
type LogConfig struct {
	// 日志等级。
	Level string `mapstructure:"level"`
	// JSON 格式日志开关。
	JSON bool `mapstructure:"json"`
}

// ai-mate-core gRPC 调用配置。
type CoreConfig struct {
	// ai-mate-core gRPC 地址。
	GRPCTarget string `mapstructure:"grpcTarget"`
	// core RPC 调用超时时间。
	Timeout time.Duration `mapstructure:"timeout"`
}

// OpenTelemetry trace 导出配置。
type ObservConfig struct {
	// trace 导出开关。
	Enabled bool `mapstructure:"enabled"`
	// OTLP HTTP exporter 地址。
	OTLPEndpoint string `mapstructure:"otlpEndpoint"`
	// trace 采样比例。
	SampleRatio float64 `mapstructure:"sampleRatio"`
}

// 读取 yaml 配置、环境变量覆盖项，并完成配置校验。
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName(resolveConfigName())
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath("../../configs")
	setDefaults(v)
	v.SetEnvPrefix("CLIENT_GATEWAY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnvs(v)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func bindEnvs(v *viper.Viper) {
	_ = v.BindEnv("appName")
	_ = v.BindEnv("env")
	_ = v.BindEnv("http.addr")
	_ = v.BindEnv("core.grpcTarget")
	_ = v.BindEnv("observability.enabled")
	_ = v.BindEnv("observability.otlpEndpoint")
	_ = v.BindEnv("observability.sampleRatio")
}

func resolveConfigName() string {
	env := strings.TrimSpace(os.Getenv("CLIENT_GATEWAY_ENV"))
	if env == "" {
		env = strings.TrimSpace(os.Getenv("ENV"))
	}
	if env == "" {
		env = "local"
	}
	return "config." + env
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("appName", "ai-mate-client-gateway")
	v.SetDefault("env", "local")
	v.SetDefault("http.addr", ":8083")
	v.SetDefault("http.readTimeout", "10s")
	v.SetDefault("http.writeTimeout", "60s")
	v.SetDefault("http.idleTimeout", "60s")
	v.SetDefault("http.shutdownTimeout", "10s")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.json", true)
	v.SetDefault("core.grpcTarget", "127.0.0.1:9090")
	v.SetDefault("core.timeout", "10s")
	v.SetDefault("observability.enabled", false)
	v.SetDefault("observability.sampleRatio", 1.0)
}

// 校验启动所需的关键配置，失败时阻止服务启动。
func (c *Config) Validate() error {
	if strings.TrimSpace(c.AppName) == "" {
		return fmt.Errorf("appName is required")
	}
	if strings.TrimSpace(c.HTTP.Addr) == "" {
		return fmt.Errorf("http.addr is required")
	}
	if c.Observ.SampleRatio < 0 || c.Observ.SampleRatio > 1 {
		return fmt.Errorf("observability.sampleRatio must be between 0 and 1")
	}
	return nil
}
