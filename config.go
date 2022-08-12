package atta

// Config for AttaReport
type Config struct {
	AttaID     string `yaml:"atta_id"`     // atta ID
	AttaToken  string `yaml:"atta_token"`  // atta token
	RetryTime  int    `yaml:"retry_time"`  // 上报atta失败重试次数
	AutoEscape bool   `yaml:"auto_escape"` // 是否打开自动转义
}

// Option 声明cache的option
type Option func(*Config)

// WithAttaInfo 设置Atta信息
func WithAttaInfo(id string, token string) Option {
	return func(c *Config) {
		c.AttaID = id
		c.AttaToken = token
	}
}

// WithRetryTime 设置重试次数
func WithRetryTime(i int) Option {
	return func(c *Config) {
		c.RetryTime = i
	}
}

// WithAutoEscape 设置是否自动转义
func WithAutoEscape(b bool) Option {
	return func(c *Config) {
		c.AutoEscape = b
	}
}
