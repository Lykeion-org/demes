package main

import ()

type Config struct {
	UseCli bool `yaml:"use_cli"`
	UseGrpcServer bool `yaml:"use_grpc_server"`
	UserServicePort string `yaml:"user_service_port"`
	DatabaseHost string `yaml:"database_host"`
	DatabasePort int `yaml:"database_port"`
	DatabaseName string `yaml:"database_name"`
	DatabaseUserName string `yaml:"database_username"`
	DatabasePassword string `yaml:"database_password"`
}