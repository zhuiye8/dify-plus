package gaia

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/gaia"
	gaiaReq "github.com/flipped-aurora/gin-vue-admin/server/model/gaia/request"
)

type TenantsService struct{}

// GetTenants 根据id获取tenants表记录
func (tenantsService *TenantsService) GetTenants(id string) (tenants gaia.Tenants, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&tenants).Error
	return
}

// GetTenantsInfoList 分页获取tenants表记录
func (tenantsService *TenantsService) GetTenantsInfoList(info gaiaReq.TenantsSearch) (list []gaia.Tenants, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&gaia.Tenants{})
	var tenantss []gaia.Tenants
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&tenantss).Error
	return tenantss, total, err
}

// GetAllTenants 获取所有工作区
func (tenantsService *TenantsService) GetAllTenants() (tenants []gaia.Tenants, err error) {
	err = global.GVA_DB.Find(&tenants).Error
	return
}
