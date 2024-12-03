package response

// OaUserInfoRes 请求OA用户信息接口返回值
type OaUserInfoRes struct {
	Code int    `json:"code"`
	Info string `json:"info"`
	Data struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	} `json:"data"`
}

// OaAccessTokenRes 请求OA access-token接口返回值
type OaAccessTokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
