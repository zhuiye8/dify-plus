package gaia

type ServiceGroup struct {
	DashboardService
	QuotaService
	ProvidersService
	TenantsService
	TestService
}
