package deps

import (
	"api-gateway/internal/client"
	"api-gateway/internal/resilience"
)

type Dependencies struct {
	AuthClient     *client.AuthClient
	UserClient     *client.UserClient
	ResumeClient   *client.ResumeClient
	VacancyClient  *client.VacancyClient
	ResponseClient *client.ResponseClient
}

func NewDependencies(cfg *Config) *Dependencies {
	authCB := resilience.NewCircuitBreaker("auth-service", 3, 0.6, 30)
	userCB := resilience.NewCircuitBreaker("user-service", 3, 0.6, 30)
	resumeCB := resilience.NewCircuitBreaker("resume-service", 3, 0.6, 30)
	vacancyCB := resilience.NewCircuitBreaker("vacancy-service", 3, 0.6, 30)
	responseCB := resilience.NewCircuitBreaker("response-service", 3, 0.6, 30)

	return &Dependencies{
		AuthClient:     client.NewAuthClient(cfg.AuthServiceURL, authCB),
		UserClient:     client.NewUserClient(cfg.UserServiceURL, userCB),
		ResumeClient:   client.NewResumeClient(cfg.ResumeServiceURL, resumeCB),
		VacancyClient:  client.NewVacancyClient(cfg.VacancyServiceURL, vacancyCB),
		ResponseClient: client.NewResponseClient(cfg.ResponseServiceURL, responseCB),
	}
}

type Config struct {
	AuthServiceURL     string
	UserServiceURL     string
	ResumeServiceURL   string
	VacancyServiceURL  string
	ResponseServiceURL string
}
