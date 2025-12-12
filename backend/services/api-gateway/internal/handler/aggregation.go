package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"api-gateway/pkg/logger"
	"api-gateway/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) GetFullProfile(c *gin.Context) {
	token := c.GetString("token")
	role := c.GetString("role")

	var wg sync.WaitGroup
	var profileData, resumeData []byte
	var profileErr, resumeErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		path := "/api/users/students/me"
		if role == "employer" {
			path = "/api/users/employers/me"
		}
		profileData, _, profileErr = h.deps.UserClient.ProxyRequest(
			c.Request.Context(),
			"GET",
			path,
			nil,
			token,
		)
	}()

	if role == "student" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resumeData, _, resumeErr = h.deps.ResumeClient.ProxyRequest(
				c.Request.Context(),
				"GET",
				"/api/resumes/me",
				nil,
				token,
			)
		}()
	}

	wg.Wait()

	if profileErr != nil {
		logger.Error("Failed to fetch profile", zap.Error(profileErr))
		response.Error(c, http.StatusServiceUnavailable, "Failed to fetch profile")
		return
	}

	var profile map[string]interface{}
	if err := json.Unmarshal(profileData, &profile); err != nil {
		logger.Error("Failed to parse profile", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Failed to parse profile")
		return
	}

	result := map[string]interface{}{
		"profile": profile,
	}

	if role == "student" && resumeErr == nil && len(resumeData) > 0 {
		var resumes map[string]interface{}
		if err := json.Unmarshal(resumeData, &resumes); err == nil {
			result["resumes"] = resumes
		}
	}

	response.Success(c, http.StatusOK, result)
}

func (h *Handler) GetDashboard(c *gin.Context) {
	token := c.GetString("token")
	role := c.GetString("role")

	var wg sync.WaitGroup
	var profileData, resumeData, vacancyData, responseData []byte
	var profileErr, resumeErr, vacancyErr, responseErr error

	if role == "student" {
		wg.Add(3)

		go func() {
			defer wg.Done()
			profileData, _, profileErr = h.deps.UserClient.ProxyRequest(
				c.Request.Context(),
				"GET",
				"/api/users/students/me",
				nil,
				token,
			)
		}()

		go func() {
			defer wg.Done()
			resumeData, _, resumeErr = h.deps.ResumeClient.ProxyRequest(
				c.Request.Context(),
				"GET",
				"/api/resumes/me",
				nil,
				token,
			)
		}()

		go func() {
			defer wg.Done()
			responseData, _, responseErr = h.deps.ResponseClient.ProxyRequest(
				c.Request.Context(),
				"GET",
				"/api/responses/me",
				nil,
				token,
			)
		}()

		wg.Wait()

		result := map[string]interface{}{}

		if profileErr == nil && len(profileData) > 0 {
			var profile map[string]interface{}
			if err := json.Unmarshal(profileData, &profile); err == nil {
				result["profile"] = profile
			}
		}

		if resumeErr == nil && len(resumeData) > 0 {
			var resumes map[string]interface{}
			if err := json.Unmarshal(resumeData, &resumes); err == nil {
				result["resumes"] = resumes
			}
		}

		if responseErr == nil && len(responseData) > 0 {
			var responses map[string]interface{}
			if err := json.Unmarshal(responseData, &responses); err == nil {
				result["responses"] = responses
			}
		}

		response.Success(c, http.StatusOK, result)

	} else if role == "employer" {
		wg.Add(3)

		go func() {
			defer wg.Done()
			profileData, _, profileErr = h.deps.UserClient.ProxyRequest(
				c.Request.Context(),
				"GET",
				"/api/users/employers/me",
				nil,
				token,
			)
		}()

		go func() {
			defer wg.Done()
			vacancyData, _, vacancyErr = h.deps.VacancyClient.ProxyRequest(
				c.Request.Context(),
				"GET",
				"/api/vacancies/me",
				nil,
				token,
			)
		}()

		go func() {
			defer wg.Done()
			responseData, _, responseErr = h.deps.ResponseClient.ProxyRequest(
				c.Request.Context(),
				"GET",
				"/api/responses/candidates",
				nil,
				token,
			)
		}()

		wg.Wait()

		result := map[string]interface{}{}

		if profileErr == nil && len(profileData) > 0 {
			var profile map[string]interface{}
			if err := json.Unmarshal(profileData, &profile); err == nil {
				result["profile"] = profile
			}
		}

		if vacancyErr == nil && len(vacancyData) > 0 {
			var vacancies map[string]interface{}
			if err := json.Unmarshal(vacancyData, &vacancies); err == nil {
				result["vacancies"] = vacancies
			}
		}

		if responseErr == nil && len(responseData) > 0 {
			var candidates map[string]interface{}
			if err := json.Unmarshal(responseData, &candidates); err == nil {
				result["candidates"] = candidates
			}
		}

		response.Success(c, http.StatusOK, result)

	} else {
		response.Error(c, http.StatusForbidden, "Access denied")
	}
}
