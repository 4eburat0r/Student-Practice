userRepo := postgres.NewUserRepo(db)
studentRepo := postgres.NewStudentRepo(db)
employerRepo := postgres.NewEmployerRepo(db)

userService := service.NewUserService(userRepo)
studentService := service.NewStudentService(studentRepo, userRepo)
employerService := service.NewEmployerService(employerRepo, userRepo)

