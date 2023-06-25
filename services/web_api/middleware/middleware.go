package middleware

import "verve_challenge_web_api/service"

type ServiceMiddleware func(service.Service) service.Service
