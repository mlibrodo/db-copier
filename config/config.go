package config

type Config struct {
	Logging LoggingConfig `mapstructure:"logging"`
	AWS     AWSConfig     `mapstructure:"aws"`
	Backup  BackupConfig  `mapstructure:"backup"`
}

type BackupConfig struct {
	S3Bucket string `mapstructure:"s3_bucket"`
}

type AWSConfig struct {
	ServiceConfigProfile string `mapstructure:"serviceConfigProfile"`
	AccountId            string `mapstructure:"accountId"`
}

type LoggingConfig struct {
	Level   string `mapstructure:"level"`
	JSONLog bool   `mapstructure:"jsonLogs"`
}
