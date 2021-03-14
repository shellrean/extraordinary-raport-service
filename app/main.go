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
    _classroomSubjectPlanRepo "github.com/shellrean/extraordinary-raport/services/classroom_subject_plan/repository/postgres"
    _classroomSubjectPlanUsecase "github.com/shellrean/extraordinary-raport/services/classroom_subject_plan/usecase"
    _classroomSubjectPlanResultRepo "github.com/shellrean/extraordinary-raport/services/classroom_subject_plan_result/repository/postgres"
    _classroomSubjectPlanResultUsecase "github.com/shellrean/extraordinary-raport/services/classroom_subject_plan_result/usecase"
    _exschoolRepo "github.com/shellrean/extraordinary-raport/services/exschool/repository/postgres"
    _exschoolUsecase "github.com/shellrean/extraordinary-raport/services/exschool/usecase"
    _exschoolStudentRepo "github.com/shellrean/extraordinary-raport/services/exschool_student/repository/postgres"
    _exschoolStudentUsecase "github.com/shellrean/extraordinary-raport/services/exschool_student/usecase"
    _studentNoteRepo "github.com/shellrean/extraordinary-raport/services/student_note/repository/postgres"
    _studentNoteUsecase "github.com/shellrean/extraordinary-raport/services/student_note/usecase"
    _attendanceRepo "github.com/shellrean/extraordinary-raport/services/attendance/repository/postgres"
    _attendanceUsecase "github.com/shellrean/extraordinary-raport/services/attendance/usecase"
    _middleware "github.com/shellrean/extraordinary-raport/interface/http/middleware"
    academicHandler "github.com/shellrean/extraordinary-raport/services/academic/interface/http"
    majorHandler "github.com/shellrean/extraordinary-raport/services/major/interface/http"
    classroomHandler "github.com/shellrean/extraordinary-raport/services/classroom/interface/http"
    classroomAcademicHandler "github.com/shellrean/extraordinary-raport/services/classroom_academic/interface/http"
    classroomStudentHandler "github.com/shellrean/extraordinary-raport/services/classroom_student/interface/http"
    classroomSubjectHandler "github.com/shellrean/extraordinary-raport/services/classroom_subject/interface/http"
    classroomSubjectPlanHandler "github.com/shellrean/extraordinary-raport/services/classroom_subject_plan/interface/http"
    classroomSubjectPlanResultHandler "github.com/shellrean/extraordinary-raport/services/classroom_subject_plan_result/interface/http"
    settingHandler "github.com/shellrean/extraordinary-raport/services/setting/interface/http"
    studentHandler "github.com/shellrean/extraordinary-raport/services/student/interface/http"
    subjectHandler "github.com/shellrean/extraordinary-raport/services/subject/interface/http"
    userHandler "github.com/shellrean/extraordinary-raport/services/user/interface/http"
    exschoolHandler "github.com/shellrean/extraordinary-raport/services/exschool/interface/http"
    exschoolStudentHandler "github.com/shellrean/extraordinary-raport/services/exschool_student/interface/http"
    studentNoteHandler "github.com/shellrean/extraordinary-raport/services/student_note/interface/http"
    attendanceHandler "github.com/shellrean/extraordinary-raport/services/attendance/interface/http"
)

func main() {
    cfgPath, err := config.ParseFlags()
    if err != nil {
        log.Fatal(err)
    }

    cfg, err := config.New(cfgPath)
    if err != nil {
        log.Fatal(err)
    }
    cfg.JWTAccessKey = "secret"
    cfg.JWTRefreshKey = "refreshsecret"
    cfg.JWTFileKey = "filesecret"
    cfg.Storage.TmpPath = "/storage/app/_tmp"
    cfg.Storage.PublicPath = "/storage/app/public"

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

    userRepo := _userRepo.New(db)
    userCacheRepo := _userCacheRepo.New(redis)
    userUsecase := _userUsecase.New(userRepo, userCacheRepo, timeoutContext, cfg)

    studentRepo := _studentRepo.New(db)
    studentUsecase := _studnetUsecase.New(studentRepo, timeoutContext, cfg)

    academicRepo := _academicRepo.New(db)
    academicUsecase := _academicUsecase.New(academicRepo, timeoutContext, cfg)
    
    majorRepo := _majorRepo.New(db)
    majorUsecase := _majorUsecase.New(majorRepo, timeoutContext, cfg)

    classroomRepo := _classroomRepo.New(db)
    classroomUsecase := _classroomUsecase.New(classroomRepo, majorRepo, timeoutContext, cfg)

    settingRepo := _settingRepo.New(db)
    settingUsecase := _settingUsecase.New(
        settingRepo, 
        academicRepo,
        timeoutContext, 
        cfg,
    )

    classroomAcademicRepo := _classroomAcademicRepo.New(db)
    classroomAcademicUsecase := _classroomAcademicUsecase.New(
        classroomAcademicRepo, 
        settingRepo, 
        userRepo,
        classroomRepo,
        timeoutContext, 
        cfg,
    )

    subjectRepo := _subjectRepo.New(db)
    subjectUsecase := _subjectUsecase.New(subjectRepo, timeoutContext, cfg)

    classroomStudentRepo := _classroomStudentRepo.New(db)
    classroomStudentUsecase := _classroomStudentUsecase.New(
        classroomStudentRepo, 
        classroomAcademicRepo,
        settingRepo,
        timeoutContext, 
        cfg,
    )

    classroomSubjectRepo := _classroomSubjectRepo.New(db)
    classroomSubjectUsecase := _classroomSubjectUsecase.New(
        classroomSubjectRepo, 
        classroomAcademicRepo,
        subjectRepo,
        userRepo,
        settingRepo,
        timeoutContext, 
        cfg,
    )

    classroomSubjectPlanRepo := _classroomSubjectPlanRepo.New(db)
    classroomSubjectPlanUsecase := _classroomSubjectPlanUsecase.New(
        classroomSubjectPlanRepo,
        userRepo,
        classroomSubjectRepo,
        classroomAcademicRepo,
        settingRepo,
        timeoutContext,
        cfg,
    )

    classroomSubjectPlanResultRepo := _classroomSubjectPlanResultRepo.New(db)
    classroomSubjectPlanResultUsecase := _classroomSubjectPlanResultUsecase.New(
        classroomSubjectPlanResultRepo,
        classroomStudentRepo,
        classroomSubjectRepo,
        classroomSubjectPlanRepo,
        timeoutContext,
        cfg,
    )

    exschoolRepo := _exschoolRepo.New(db)
    exschoolUsecase := _exschoolUsecase.New(exschoolRepo, timeoutContext, cfg)

    exschoolStudentRepo := _exschoolStudentRepo.New(db)
    exschoolStudentUsecase := _exschoolStudentUsecase.New(
        exschoolStudentRepo,
        exschoolRepo,
        classroomStudentRepo,
        classroomAcademicRepo,
        timeoutContext,
        cfg,
    )

    studentNoteRepo := _studentNoteRepo.New(db)
    studentNoteUsecase := _studentNoteUsecase.New(
        studentNoteRepo,
        classroomStudentRepo,
        userRepo,
        timeoutContext,
        cfg,
    )

    attendanceRepo := _attendanceRepo.New(db)
    attendanceUsecase := _attendanceUsecase.New(
        attendanceRepo,
        classroomStudentRepo,
        timeoutContext,
        cfg,
    )

    if cfg.Release == true {
        gin.SetMode(gin.ReleaseMode)
    }
    
    r := gin.Default()

    mddl := _middleware.Init(cfg)

    r.Use(mddl.CORS())
    
    userHandler.New(r, userUsecase, cfg, mddl)
    settingHandler.New(r, settingUsecase, cfg, mddl)
    majorHandler.New(r, majorUsecase, cfg, mddl)
    studentHandler.New(r, studentUsecase, cfg, mddl)
    subjectHandler.New(r, subjectUsecase, cfg, mddl)
    academicHandler.New(r, academicUsecase, cfg, mddl)
    classroomHandler.New(r, classroomUsecase, cfg, mddl)
    classroomAcademicHandler.New(r, classroomAcademicUsecase, cfg, mddl)
    classroomStudentHandler.New(r, classroomStudentUsecase, cfg, mddl)
    classroomSubjectHandler.New(r, classroomSubjectUsecase, cfg, mddl)
    classroomSubjectPlanHandler.New(r, classroomSubjectPlanUsecase, cfg, mddl)
    classroomSubjectPlanResultHandler.New(r, classroomSubjectPlanResultUsecase, cfg, mddl)
    exschoolHandler.New(r, exschoolUsecase, cfg, mddl)
    exschoolStudentHandler.New(r, exschoolStudentUsecase, cfg, mddl)
    studentNoteHandler.New(r, studentNoteUsecase, cfg, mddl)
    attendanceHandler.New(r, attendanceUsecase, cfg, mddl)

    // Let's run our extraordinary-raport server
    fmt.Printf("Extraordinary-raport serve on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
    err = r.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
    if err != nil {
        log.Fatal(err)
    }
}