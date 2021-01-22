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
    _majorRepo "github.com/shellrean/extraordinary-raport/services/major/repository/postgres"
    _majorUsecase "github.com/shellrean/extraordinary-raport/services/major/usecase"
    _settingRepo "github.com/shellrean/extraordinary-raport/services/setting/repository/postgres"
    _settingUsecase "github.com/shellrean/extraordinary-raport/services/setting/usecase"
    _classroomAcademicRepo "github.com/shellrean/extraordinary-raport/services/classroom_academic/repository/postgres"
    _classroomAcademicUsecase "github.com/shellrean/extraordinary-raport/services/classroom_academic/usecase"
    _subjectRepo "github.com/shellrean/extraordinary-raport/services/subject/repository/postgres"
    _subjectUsecase "github.com/shellrean/extraordinary-raport/services/subject/usecase"
    _classroomStudentRepo "github.com/shellrean/extraordinary-raport/services/classroom_student/repository/postgres"
    _classroomStudentUsecase "github.com/shellrean/extraordinary-raport/services/classroom_student/usecase"
    _classroomSubjectRepo "github.com/shellrean/extraordinary-raport/services/classroom_subject/repository/postgres"
    _classroomSubjectUsecase "github.com/shellrean/extraordinary-raport/services/classroom_subject/usecase"
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
    
    majorRepo := _majorRepo.NewPostgresMajorRepository(db)
    majorUsecase := _majorUsecase.NewMajorUsecase(majorRepo, timeoutContext, cfg)

    classroomRepo := _classroomRepo.NewPostgresClassroomRepository(db)
    classroomUsecase := _classroomUsecase.NewClassroomUsecase(classroomRepo, majorRepo, timeoutContext, cfg)

    settingRepo := _settingRepo.NewPostgresSettingRepository(db)
    settingUsecase := _settingUsecase.NewSettingUsecase(settingRepo, timeoutContext, cfg)

    classroomAcademicRepo := _classroomAcademicRepo.NewPostgresClassroomAcademicRepository(db)
    classroomAcademicUsecase := _classroomAcademicUsecase.NewClassroomAcademicUsecase(
        classroomAcademicRepo, 
        settingRepo, 
        userRepo,
        classroomRepo,
        timeoutContext, 
        cfg,
    )

    subjectRepo := _subjectRepo.NewPostgresSubjectRepository(db)
    subjectUsecase := _subjectUsecase.NewSubjectUsecase(subjectRepo, timeoutContext, cfg)

    classroomStudentRepo := _classroomStudentRepo.NewPostgresClassroomStudentRepository(db)
    classroomStudentUsecase := _classroomStudentUsecase.NewClassroomStudentUsecase(
        classroomStudentRepo, 
        classroomAcademicRepo,
        settingRepo,
        timeoutContext, 
        cfg,
    )

    classroomSubjectRepo := _classroomSubjectRepo.NewPostgresClassroomSubjectRepository(db)
    classroomSubjectUsecase := _classroomSubjectUsecase.NewClassroomSubjectUsecase(
        classroomSubjectRepo, 
        classroomAcademicRepo,
        subjectRepo,
        userRepo,
        timeoutContext, 
        cfg,
    )

    if cfg.Release == true {
        gin.SetMode(gin.ReleaseMode)
    }
    
    r := gin.Default()

    mddl := _middleware.InitMiddleware(cfg)

    r.Use(mddl.CORS())

    httpHandler.NewUserHandler(r, userUsecase, cfg, mddl)
    httpHandler.NewStudentHandler(r, studentUsecase, cfg, mddl)
    httpHandler.NewAcademicHandler(r, academicUsecase, cfg, mddl)
    httpHandler.NewClassroomHandler(r, classroomUsecase, cfg, mddl)
    httpHandler.NewMajorHandler(r, majorUsecase, cfg, mddl)
    httpHandler.NewClassAcademicHandler(r, classroomAcademicUsecase, cfg, mddl)
    httpHandler.NewSubjectHandler(r, subjectUsecase, cfg, mddl)
    httpHandler.NewClassroomStudentHandler(r, classroomStudentUsecase, cfg, mddl)
    httpHandler.NewClassroomSubjectHandler(r, classroomSubjectUsecase, cfg, mddl)
    httpHandler.NewSettingHandler(r, settingUsecase, cfg, mddl)

    // Let's run our extraordinary-raport server
    fmt.Printf("Extraordinary-raport serve on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
    err = r.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
    if err != nil {
        log.Fatal(err)
    }
}