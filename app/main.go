package main

import (
    "log"
    "fmt"
    "database/sql"
    "time"
    "context"

    _ "github.com/lib/pq"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"

    "github.com/shellrean/extraordinary-raport/config"
    _userRepo "github.com/shellrean/extraordinary-raport/services/user/repository/postgres"
    _userCacheRepo "github.com/shellrean/extraordinary-raport/services/user/repository/redis"
    _userUsecase "github.com/shellrean/extraordinary-raport/services/user/usecase"
    _studentRepo "github.com/shellrean/extraordinary-raport/services/student/repository/postgres"
    _studnetUsecase "github.com/shellrean/extraordinary-raport/services/student/usecase"
    _academicRepo "github.com/shellrean/extraordinary-raport/services/academic/repository/postgres"
    _academicUsecase "github.com/shellrean/extraordinary-raport/services/academic/usecase"
    _classroomRepo "github.com/shellrean/extraordinary-raport/services/classroom/repository/postgres"
    _classroomUsecase "github.com/shellrean/extraordinary-raport/services/classroom/usecase"
    _middleware "github.com/shellrean/extraordinary-raport/interface/http/middleware"
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
    cfg.JWTAccessKey = "secret"
    cfg.JWTRefreshKey = "refreshsecret"

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
        
    redisDsn := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
    redis := redis.NewClient(&redis.Options{
        Addr: redisDsn,
        Password: cfg.Redis.Password,
        DB: cfg.Redis.DBName,
    })

    if cfg.Redis.Enable {
        _, err = redis.Ping(context.TODO()).Result()
        if err != nil {
            log.Fatal(err)
        }
    }

    timeoutContext := time.Duration(cfg.Context.Timeout) * time.Second

    userRepo := _userRepo.NewPostgresUserRepository(db)
    userCacheRepo := _userCacheRepo.NewRedisUserRepository(redis)
    userUsecase := _userUsecase.NewUserUsecase(userRepo, userCacheRepo, timeoutContext, cfg)

    studentRepo := _studentRepo.NewPostgresStudentRepository(db)
    studentUsecase := _studnetUsecase.NewStudentUsecase(studentRepo, timeoutContext, cfg)

    academicRepo := _academicRepo.NewPostgresAcademicRepository(db)
    academicUsecase := _academicUsecase.NewAcademicUsecase(academicRepo, timeoutContext, cfg)

    classroomRepo := _classroomRepo.NewPostgresClassroomRepository(db)
    classroomUsecase := _classroomUsecase.NewClassroomUsecase(classroomRepo, timeoutContext, cfg)

    if cfg.Release == true {
        gin.SetMode(gin.ReleaseMode)
    }
    
    r := gin.Default()

    mddl := _middleware.InitMiddleware(cfg)

    httpHandler.NewUserHandler(r, userUsecase, cfg, mddl)
    httpHandler.NewStudentHandler(r, studentUsecase, cfg, mddl)
    httpHandler.NewAcademicHandler(r, academicUsecase, cfg, mddl)
    httpHandler.NewClassroomHandler(r, classroomUsecase, cfg, mddl)

    // Let's run our extraordinary-raport server
    fmt.Printf("Extraordinary-raport serve on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
    err = r.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
    if err != nil {
        log.Fatal(err)
    }
}