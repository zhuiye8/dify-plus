package request

type OaLoginReq struct {
	AuthorizeCode string `json:"authorize_code" form:"authorize_code"` // OA返回的授权验证码，用于请求用户信息
}
