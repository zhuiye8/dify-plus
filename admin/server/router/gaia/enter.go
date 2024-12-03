package gaia

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	DashboardRouter
	QuotaRouter
	ProvidersRouter
	TenantsRouter
	TestRouter
}

var (
	dashboardApi = api.ApiGroupApp.GaiaApiGroup.DashboardApi
	providersApi = api.ApiGroupApp.GaiaApiGroup.ProvidersApi
	tenantsApi   = api.ApiGroupApp.GaiaApiGroup.TenantsApi
)
var quotaApi = api.ApiGroupApp.GaiaApiGroup.QuotaApi
var testApi = api.ApiGroupApp.GaiaApiGroup.TestApi
