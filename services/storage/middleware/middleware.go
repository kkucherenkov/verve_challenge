package middleware

import "verve_challenge_storage/service"

type ServiceMiddleware func(service.Service) service.Service
