package configuration

import (
	"context"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/jackc/pgx/v5/pgxpool"

	"jdgonzalez907/saas-api/internal/users/application"
	"jdgonzalez907/saas-api/internal/users/infrastructure/controllers"
	"jdgonzalez907/saas-api/internal/users/infrastructure/database"
)

type ContentInformation struct {
	title   string
	content string
}

type ContentInformationParams struct {
	Title   string
	Content string
}

func NewContentInformation(params ContentInformationParams) (ContentInformation, error) {
	// Validaciones aquí

	return ContentInformation{
		title:   params.Title,
		content: params.Content,
	}, nil
}

func (c ContentInformation) Title() string {
	return c.title
}

func (c ContentInformation) Content() string {
	return c.content
}

type Config struct {
	HTTPPort    string `env:"HTTP_PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
}

type Application struct {
	server *http.Server
	db     *pgxpool.Pool
}

func NewApplication() (*Application, error) {
	cfg := new(Config)

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := NewPostgresConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	userRepository := database.NewPostgresUserRepository(db)

	findUserByIDUseCase := application.NewFindUserByIDUseCase(userRepository)
	createUserUseCase := application.NewCreateUserUseCase(userRepository)
	deleteUserUseCase := application.NewDeleteUserUseCase(userRepository)
	updatePersonalInformationUseCase := application.NewUpdateUserPersonalInformationUseCase(userRepository)
	changeEmailUseCase := application.NewChangeUserEmailUseCase(userRepository)
	changePhoneUseCase := application.NewChangeUserPhoneUseCase(userRepository)

	findUserByID := controllers.NewFindUserByIDController(findUserByIDUseCase)
	createUser := controllers.NewCreateUserController(createUserUseCase)
	deleteUser := controllers.NewDeleteUserController(deleteUserUseCase)
	updatePersonalInformation := controllers.NewUpdateUserPersonalInformationController(updatePersonalInformationUseCase)
	changeEmail := controllers.NewChangeUserEmailController(changeEmailUseCase)
	changePhone := controllers.NewChangeUserPhoneController(changePhoneUseCase)

	router := controllers.NewRouter(controllers.RouterParams{
		FindUserByID:              findUserByID,
		CreateUser:                createUser,
		DeleteUser:                deleteUser,
		UpdatePersonalInformation: updatePersonalInformation,
		ChangeEmail:               changeEmail,
		ChangePhone:               changePhone,
	})

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	return &Application{
		server: server,
		db:     db,
	}, nil
}

func (a *Application) Run() error {
	return a.server.ListenAndServe()
}

func (a *Application) Close() error {
	if a.db != nil {
		a.db.Close()
	}

	return nil
}
