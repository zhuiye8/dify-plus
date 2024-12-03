package gaia

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia/response"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type ProvidersService struct{}

// SyncProviders 同步模型记录
// Author [yourname](https://github.com/yourname)
func (providersService *ProvidersService) SyncProviders(jwtToken string, req gaiaReq.SyncProviderReq) (err error) {
	//var providers gaia.Providers
	//err = global.GVA_DB.Model(&gaia.Providers{}).Where("id = ?", providers.Id).Updates(&providers).Error

	token := jwtToken

	resp, err := providersService.syncAzureOpenAIModels(token, req)
	if resp.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("请求Gaia-API鉴权不通过，响应信息：%s", resp.String())
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("请求Gaia-API失败，响应信息：%s", resp.String())
	}
	if err != nil {
		return err
	}
	return err
}

func (providersService *ProvidersService) syncAzureOpenAIModels(token string, req gaiaReq.SyncProviderReq) (*resty.Response, error) {
	// 构建请求体和请求域名
	url := global.GVA_CONFIG.Gaia.Url
	var requestBody = make(map[string]interface{})
	if req.ProviderModelName != "" {
		url += fmt.Sprintf("/console/api/workspaces/current/model-providers/%s/models/sync", req.ProviderName)

		// TODO 这里model_type 有个映射关系，暂时先直接赋值，后续再做处理
		var modelType = ""
		switch req.ModelType {
		case "text-generation":
			modelType = "llm"
		}

		requestBody = map[string]interface{}{
			"model":      req.ProviderModelName,
			"model_type": modelType,
			"is_all":     req.IsAll,
			"tenant_ids": req.TenantIds,
		}
	} else {
		url += fmt.Sprintf("/console/api/workspaces/current/model-providers/%s/sync", req.ProviderName)
		requestBody = map[string]interface{}{
			"is_all":     req.IsAll,
			"tenant_ids": req.TenantIds,
		}
	}
	global.GVA_LOG.Info(url)
	global.GVA_LOG.Info("requestBody", zap.Any("requestBody", requestBody))

	client := resty.New()

	// TODO jwt-token
	//token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMjRhMzVjMTctZTNiMS00ZmIzLTliYmQtOTFlZmM2ZDg5ZTY4IiwiZXhwIjoxNzMyMDA4MjU3LCJpc3MiOiJTRUxGX0hPU1RFRCIsInN1YiI6IkNvbnNvbGUgQVBJIFBhc3Nwb3J0In0.p1FjxGy34nFYNGLYYJjsjNr-vUzT6YQiZ3tA9ghUMpI"

	// 发送请求
	resp, err := client.R().
		SetHeader("Accept", "*/*").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Cache-Control", "no-cache").
		SetHeader("Connection", "keep-alive").
		SetHeader("Origin", "http://127.0.0.1:3000").
		SetHeader("Pragma", "no-cache").
		SetHeader("Referer", "http://127.0.0.1:3000/").
		SetHeader("User-Agent", "Gaia-Admin-API").
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(url)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetProviders 根据id获取providers表记录
// Author [yourname](https://github.com/yourname)
func (providersService *ProvidersService) GetProviders(id string) (res response.FindProvidersRes, err error) {
	var providers gaia.Providers
	err = global.GVA_DB.Where("id = ?", id).First(&providers).Error
	var providerModels gaia.ProviderModels
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = global.GVA_DB.Where("id = ?", id).First(&providerModels).Error
		if err != nil {
			return
		}
	}
	var resId, resProviderName, resModelName, resModelType, resTenantID string
	if providers.Id != "" {
		resId = providers.Id
		resProviderName = providers.ProviderName
		resModelName = ""
		resModelType = ""
		resTenantID = providers.TenantId
	} else {
		resId = providerModels.ID
		resProviderName = providerModels.ProviderName
		resModelName = providerModels.ModelName
		resModelType = providerModels.ModelType
		resTenantID = providerModels.TenantID
	}

	// 查询同步配置
	var modelSyncConfigExtend gaia.ModelSyncConfigExtend
	err = global.GVA_DB.Where("model_id = ?", resId).Find(&modelSyncConfigExtend).Error
	if err != nil {
		return
	}

	// 查询已同步的工作区
	var tenantModelSyncExtends []gaia.TenantModelSyncExtend
	err = global.GVA_DB.Model(&gaia.TenantModelSyncExtend{}).Where("origin_model_id = ?", resId).Find(&tenantModelSyncExtends).Error
	if err != nil {
		return
	}
	var syncedTenants []string
	for _, tenantModelSync := range tenantModelSyncExtends {
		syncedTenants = append(syncedTenants, tenantModelSync.TenantID)
	}
	// 拼接结果
	res = response.FindProvidersRes{
		ID:                resId,
		TenantID:          resTenantID,
		ProviderName:      resProviderName,
		ProviderModelName: resModelName,
		ModelType:         resModelType,
		IsAll:             modelSyncConfigExtend.IsAll,
		TenantIDs:         syncedTenants,
	}

	return
}

// GetProvidersInfoList 分页获取providers表记录
// Author [yourname](https://github.com/yourname)
func (providersService *ProvidersService) GetProvidersInfoList(info gaiaReq.ProvidersSearch) (list []response.FindProvidersRes, total int64, err error) {
	if global.GVA_CONFIG.Gaia.SuperAdminTenantId == "" {
		err = fmt.Errorf("请在配置里设置系统默认工作区ID SUPER_ADMIN_TENANT_ID")
		return
	}
	adminTenantId := global.GVA_CONFIG.Gaia.SuperAdminTenantId
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&gaia.Providers{}).Where("tenant_id = ?", adminTenantId)
	var providerss []gaia.Providers
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&providerss).Error
	if err != nil {
		return
	}
	var providerIDs []string
	for _, provider := range providerss {
		providerIDs = append(providerIDs, provider.Id)
	}

	// 查询已同步的工作区
	var tenantProviderSyncExtends []gaia.TenantModelSyncExtend
	err = global.GVA_DB.Model(&gaia.TenantModelSyncExtend{}).Where("origin_model_id in ?", providerIDs).Find(&tenantProviderSyncExtends).Error
	if err != nil {
		return
	}
	var syncedTenants = make(map[string][]string)
	for _, tenantModelSync := range tenantProviderSyncExtends {
		syncedTenants[tenantModelSync.OriginModelID] = append(syncedTenants[tenantModelSync.OriginModelID], tenantModelSync.TenantID)
	}

	// 查询同步配置
	var providerSyncConfigExtend []gaia.ModelSyncConfigExtend
	err = global.GVA_DB.Where("model_id in ?", providerIDs).Find(&providerSyncConfigExtend).Error
	if err != nil {
		return
	}
	var providerIsAllConfig = make(map[string]bool)
	for _, providerSyncConfig := range providerSyncConfigExtend {
		providerIsAllConfig[providerSyncConfig.ModelID] = providerSyncConfig.IsAll
	}

	for _, provider := range providerss {
		list = append(list, response.FindProvidersRes{
			ID:                provider.Id,
			TenantID:          provider.TenantId,
			ProviderName:      provider.ProviderName,
			ProviderModelName: "",
			ModelType:         "",
			TenantIDs:         syncedTenants[provider.Id],
			IsAll:             providerIsAllConfig[provider.Id],
		})
	}

	// 再从provider
	var providerModels []gaia.ProviderModels
	err = global.GVA_DB.Where("tenant_id = ?", adminTenantId).Find(&providerModels).Error
	if err != nil {
		return
	}
	var providerModelIDs []string
	for _, providerModel := range providerModels {
		providerModelIDs = append(providerModelIDs, providerModel.ID)
	}

	// 查询同步配置
	var providerModelSyncConfigExtend []gaia.ModelSyncConfigExtend
	err = global.GVA_DB.Where("model_id in ?", providerModelIDs).Find(&providerModelSyncConfigExtend).Error
	if err != nil {
		return
	}
	var providerModelIsAllConfig = make(map[string]bool)
	for _, providerSyncConfig := range providerModelSyncConfigExtend {
		providerModelIsAllConfig[providerSyncConfig.ModelID] = providerSyncConfig.IsAll
	}

	// 查询已同步的工作区
	var tenantModelSyncExtends []gaia.TenantModelSyncExtend
	err = global.GVA_DB.Model(&gaia.TenantModelSyncExtend{}).Where("origin_model_id in ?", providerModelIDs).Find(&tenantModelSyncExtends).Error
	if err != nil {
		return
	}
	for _, tenantModelSync := range tenantModelSyncExtends {
		syncedTenants[tenantModelSync.OriginModelID] = append(syncedTenants[tenantModelSync.OriginModelID], tenantModelSync.TenantID)
	}

	for _, providerModel := range providerModels {
		list = append(list, response.FindProvidersRes{
			ID:                providerModel.ID,
			TenantID:          providerModel.TenantID,
			ProviderName:      providerModel.ProviderName,
			ProviderModelName: providerModel.ModelName,
			ModelType:         providerModel.ModelType,
			TenantIDs:         syncedTenants[providerModel.ID],
			IsAll:             providerModelIsAllConfig[providerModel.ID],
		})
	}

	return list, total, err
}
