package models

//implementing Error interface
type ErrorResponse struct{
	Message string `json:"message" example:"An Error Occured"`
	StatusCode int `json:"status_code" example:"404"`
}

func (e *ErrorResponse) Error() string{
	return e.Message
}

//ONLY FOR SWAGGER VISUALIZATION PURPOSES
type BadRequestResponse struct{
	Message string `json:"message" example:"Name may not be empty"`
	StatusCode int `json:"status_code" example:"400"`
}

type UnauthorizedResponse struct{
	Message string `json:"message" example:"Invalid or expired token"`
	StatusCode int `json:"status_code" example:"401"`
}

type ForbiddenResponse struct{
	Message string `json:"message" example:"Insufficient access level"`
	StatusCode int `json:"status_code" example:"403"`
}