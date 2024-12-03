package request

type SetUserQuotaRequest struct {
	Uid   string  `json:"uid" form:"uid"`     // 用户id
	Quota float64 `json:"quota" form:"quota"` // 额度
}
