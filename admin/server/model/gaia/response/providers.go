package response

type GetProvidersInfoListRes struct {
	ID string `json:"id"`
}

// FindProvidersRes 查询模型相关信息
type FindProvidersRes struct {
	ID                string   `json:"id"`
	TenantID          string   `json:"tenant_id"`           // 工作区ID
	IsAll             bool     `json:"is_all"`              // 是否同步所有
	TenantIDs         []string `json:"tenant_ids"`          // 已同步工作空间ID集合
	ProviderName      string   `json:"provider_name"`       // 模型供应商名称
	ProviderModelName string   `json:"provider_model_name"` // 模型名称
	ModelType         string   `json:"model_type"`          // 模型类型
}
