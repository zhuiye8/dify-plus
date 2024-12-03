package response

// GetAccountQuotaRankingDataRes 获取账户配额排名数据的响应结构
type GetAccountQuotaRankingDataRes struct {
	Ranking    int     `json:"ranking"`     // 排名
	Name       string  `json:"name"`        // 姓名
	UsedQuota  float64 `json:"used_quota"`  // 已使用配额
	TotalQuota float64 `json:"total_quota"` // 总配额
}

// GetAppQuotaRankingDataRes 获取应用配额排名数据的响应结构
type GetAppQuotaRankingDataRes struct {
	Ranking      int     `json:"ranking"`      // 排名
	Name         string  `json:"name"`         // 应用名称
	Mode         string  `json:"mode"`         // 应用类型
	TenantName   string  `json:"tenant_name"`  // 工作区名称
	AccountName  string  `json:"account_name"` // 账号名称
	UsedQuota    float64 `json:"used_quota"`   // 已使用配额
	AppID        string  `json:"app_id"`
	TotalCost    float64 `json:"total_cost"`
	MessageCost  float64 `json:"message_cost"`
	WorkflowCost float64 `json:"workflow_cost"`
	RecordNum    float64 `json:"record_num"`
	UseNum       int     `json:"use_num"`
}

// GetAppTokenQuotaRankingDataRes 获取应用密钥配额排名数据的响应结构
type GetAppTokenQuotaRankingDataRes struct {
	Ranking          int     `json:"ranking"`           // 排名
	Name             string  `json:"name"`              // 对应应用名称
	AppToken         string  `json:"app_token"`         // 密钥（需要加密显示）
	AccumulatedQuota float64 `json:"accumulated_quota"` // 累计使用
	DayUsedQuota     float64 `json:"day_used_quota"`    // 日使用
	MonthUsedQuota   float64 `json:"month_used_quota"`  // 月使用
	DayLimitQuota    float64 `json:"day_limit_quota"`   // 日限额
	MonthLimitQuota  float64 `json:"month_limit_quota"` // 月限额
}

type GetAppTokenDailyQuotaDataRes struct {
	StatDate  string  `json:"stat_date"`
	TotalUsed float64 `json:"total_used"`
}

// GetAiImageQuotaRankingRes 获取AI图片配额排名数据的响应结构
type GetAiImageQuotaRankingRes struct {
	Ranking   int     `json:"ranking"`    // 排名
	Address   string  `json:"address"`    // 域名
	Path      string  `json:"path"`       // 路径
	Model     string  `json:"model"`      // 模型
	TotalCost float64 `json:"total_cost"` // 总花费
	RecordNum int     `json:"record_num"` // 调用次数
}
