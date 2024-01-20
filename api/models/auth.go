package models

import "time"

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,min=6"`
}

type AuthResponse struct {
	AccessToken string    `json:"access_token"`
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
}