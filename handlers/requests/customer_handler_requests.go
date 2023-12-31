package requests

type SignupRequest struct {
	Name  string `json:"name" binding:"required"`
	Plan  string `json:"plan" binding:"required,oneof=basic advanced enterprise"`
	Email string `json:"email" binding:"required,email"`
}
type ReportsRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

