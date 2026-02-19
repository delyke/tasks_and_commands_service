package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type HTTPConfig interface {
	Host() string
	Port() string
	Address() string
	ReadTimeout() time.Duration
}

type RedisConfig interface {
	Host() string
	Port() int
	Password() string
	DB() int
	Address() string
}

type JWTConfig interface {
	AccessTokenTTL() time.Duration
	Issuer() string
	SecretKey() string
}

type MySQLConfig interface {
	Host() string
	Port() int
	User() string
	Password() string
	Database() string
	MaxOpenConns() int
	MaxIdleConns() int
	ConnMaxLifetime() time.Duration
	DSN() string
}
