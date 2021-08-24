package auth

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInGoogleReq struct {
	GoogleAccessToken string `json:"googleAccessToken"`
}

type SignInResp struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}

type SignUpReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResp struct {
	User  User  `json:"user"`
	Token Token `json:"token"`
}
