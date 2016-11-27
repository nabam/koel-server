package main

type config struct {
	Http     httpConfig   `toml:"http"`
	Sqlite   sqliteConfig `toml:"sqlite"`
	Database string
}

type httpConfig struct {
	Host string
	Port int16
}

type sqliteConfig struct {
	Path string
}
