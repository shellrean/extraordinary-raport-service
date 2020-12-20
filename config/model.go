package config

import "time"

type Config struct {
    Server struct {
        // Host is the local machine IP Address to bind the HTTP Server to
        Host    string `yaml:"host"`

        // Port is the local machine TCP Port to bind the HTTP Server to
        Port    string `yaml:"port"`
        Tiemout struct {
            // Server is the general server timeout to use
            // for graceful shutdowns
            Server  time.Duration `yaml:"server"`

            // Write is the amount of time to wait until an HTTP Server
            // write operation is cancelled
            Write   time.Duration `yaml:"write"`

            // Read is the amount of time to wait until an HTTP Server
            // read operation is cancelled
            Read    time.Duration `yaml:"read"`

            // Read is the amount of time to wait
            // until an IDLE HTTP Session is closed
            Idle    time.Duration `yaml:"idle"`
        } `yaml:"timeout"`
    } `yaml:"server"`

    Database struct {
        Username    string  `yaml:"username"`
        Password    string  `yaml:"password"`
        DBName      string  `yaml:"dbname"`
        Host        string  `yaml:"host"`
        Port        string  `yaml:"port"`
    } `yaml:"database"`

    Context struct {
        Timeout     int     `yaml:"timeout"`
    } `yaml:"context"`

    Release bool `yaml:"release"`
}