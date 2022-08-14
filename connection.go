package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	ENV_SERVER_ADDR = "TOMQ_SERVER_ADDR"
	ENV_SERVER_PORT = "TOMQ_SERVER_PORT"
)

type config struct {
	server string
	port   string
}

type configuration interface {
	readEnv() error
	getConnectionString() string
}

func (c *config) getConnectionString() string {
	return c.server + ":" + c.port
}

func (c *config) readEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	c.server = os.Getenv(ENV_SERVER_ADDR)
	c.port = os.Getenv(ENV_SERVER_PORT)
	return nil
}
