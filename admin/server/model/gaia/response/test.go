package response

type GetAppRequestTestDataResponse struct {
	Name        string  `json:"name" gorm:"comment:应用名"`
	Status      bool    `json:"status" gorm:"comment:状态"`
	Inputs      string  `json:"inputs" gorm:"comment:输入"`
	Outputs     string  `json:"outputs" gorm:"comment:输出"`
	Error       string  `json:"error" gorm:"comment:错误信息"`
	Comparison  string  `json:"comparison" gorm:"comment:历史对照"`
	LogTime     float64 `json:"log_time" gorm:"not null;default:0;comment:旧耗时"`
	ElapsedTime float64 `json:"elapsed_time" gorm:"not null;default:0;comment:耗时"`
}
