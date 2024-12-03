package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ProvidersSearch struct {
	request.PageInfo
	TenantId string `json:"tenant_id" form:"tenant_id"`
}

// SyncProviderReq 同步模型请求值
type SyncProviderReq struct {
	ProviderName      string   `json:"provider_name" form:"provider_name"`
	ProviderModelName string   `json:"provider_model_name" form:"provider_model_name"`
	ModelType         string   `json:"model_type" form:"model_type"`
	IsAll             bool     `json:"is_all" form:"is_all"`
	TenantIds         []string `json:"tenant_ids" form:"tenant_ids"`
}
