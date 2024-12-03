package response

type GetQuotaManagementDataResponse struct {
	Uid        string  `json:"uid"`         // 用户id
	Ranking    int     `json:"ranking"`     // 排名
	Name       string  `json:"name"`        // 姓名
	UsedQuota  float64 `json:"used_quota"`  // 已使用配额
	TotalQuota float64 `json:"total_quota"` // 总配额
}
