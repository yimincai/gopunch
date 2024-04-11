package domain

type BaseResponse struct {
	Code   int    `json:"code"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

type LoginResult struct {
	UserID       string `json:"userId"`
	Account      string `json:"account"`
	UserName     string `json:"userName"`
	Role         string `json:"role"`
	Avatar       string `json:"avatar"`
	OrganizeID   string `json:"organizeId"`
	OrganizeName string `json:"organizeName"`
	AccessToken  string `json:"accessToken"`
}

type LoginResponse struct {
	BaseResponse
	Result LoginResult `json:"result"`
}

type PunchResult struct {
	IsSuccess bool   `json:"isSuccess"`
	ErrorType string `json:"errorType"`
	UserName  string `json:"userName"`
	Avatar    string `json:"avatar"`
}

type PunchResponse struct {
	BaseResponse
	Result PunchResult `json:"result"`
}

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type PunchUsecase interface {
	Punch(discordUserID string) error
}
