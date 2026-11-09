package dto

type CreateIncidentDTO struct {
	Title       string `json:"title" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required"`
	Priority    string `json:"priority" binding:"required,oneof=low medium high critical"`
}

type UpdateStatusDTO struct {
	Status string `json:"status" binding:"required,oneof=open in_progress resolved closed"`
}

type RegisterDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}
