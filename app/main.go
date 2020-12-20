package main

import (
    "log"
    "fmt"
    "database/sql"
    "time"

    _ "github.com/lib/pq"
    "github.com/gin-gonic/gin"

    "github.com/shellrean/extraordinary-raport/config"
    _userRepo "github.com/shellrean/extraordinary-raport/services/user/repository/postgres"
    _userUsecase "github.com/shellrean/extraordinary-raport/services/user/usecase"
    httpHandler "github.com/shellrean/extraordinary-raport/interface/http/handler"
)

func main() {
    cfgPath, err := config.ParseFlags()
    if err != nil {
        log.Fatal(err)
    }

    cfg, err := config.NewConfig(cfgPath)
    if err != nil {
        log.Fatal(err)
    }

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
        cfg.Database.Host, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName, cfg.Database.Port, cfg.Database.Timezone,
    )

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        err := db.Close()
        if err != nil {
            log.Fatal(err)
        }
    }()

    timeoutContext := time.Duration(cfg.Context.Timeout) * time.Second

    userRepo := _userRepo.NewPostgresUserRepository(db)
    userUsecase := _userUsecase.NewUserUsecase(userRepo, timeoutContext)

    if cfg.Release == true {
        gin.SetMode(gin.ReleaseMode)
    }
    
    r := gin.Default()

    httpHandler.NewUserHandler(r, userUsecase)
    
    // Let's run our server
    fmt.Printf("Extraordinary-raport serve on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
    err = r.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
    if err != nil {
        log.Fatal(err)
    }
}