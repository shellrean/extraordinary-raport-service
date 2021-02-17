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

    Security struct {
        // CORS Configuration
        CORS struct {
            // Host that you allowed to access the api
            Host    string      `yaml:"host"`
            
            // Method that allowed to access the api
            Method  string      `yaml:"method"`
        } `yaml:"cors"`
    } `yaml:"security"`

    Database struct {
        // Username is the database machine user
        Username    string  `yaml:"username"`

        // Password is the database machine password
        Password    string  `yaml:"password"`

        // DBName is the database name
        DBName      string  `yaml:"dbname"`

        // Host is local machine IP Address to bind the database
        Host        string  `yaml:"host"`

        // Port is local machie TCP Port to bind the database
        Port        string  `yaml:"port"`
        
        // Timezone is the your timezone's database
        Timezone    string  `yaml:"timezone"`
    } `yaml:"database"`

    Redis struct {
        // Enable will use redis as cache
        Enable      bool    `yaml:"enable"`
        Host        string  `yaml:"host"`
        Port        string  `yaml:"port"`
        Password    string  `yaml:"password"`
        DBName      int  `yaml:"dbname"`
    } `yaml:"redis"`

    Context struct {
        Timeout     int     `yaml:"timeout"`
    } `yaml:"context"`

    Release bool `yaml:"release"`

    JWTAccessKey  string

    JWTRefreshKey string

    JWTFileKey string

    Storage struct {
        TmpPath     string
        PublicPath  string
    }
}