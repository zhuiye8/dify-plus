package gaia

const TestDefaultNumber = 2     // 默认测试执行次数
const BatchStatusInProgress = 1 // 批次状态:执行中
const BatchStatusCompleted = 2  // 批次状态:已结束

// AppRequestTest APP请求测试表
type AppRequestTest struct {
	ID          uint    `json:"id" gorm:"primarykey;comment:主键"`
	AppID       string  `json:"app_id" gorm:"index;comment:应用ID"`
	BatchId     uint    `json:"batch_id" gorm:"index;comment:批次ID"`
	Status      string  `json:"status" gorm:"index;comment:状态"`
	Inputs      string  `json:"inputs" gorm:"comment:输入"`
	Outputs     string  `json:"outputs" gorm:"comment:输出"`
	Error       string  `json:"error" gorm:"comment:错误信息"`
	Comparison  string  `json:"comparison" gorm:"comment:历史对照"`
	LogTime     float64 `json:"log_time" gorm:"not null;default:0;comment:旧耗时"`
	ElapsedTime float64 `json:"elapsed_time" gorm:"not null;default:0;comment:耗时"`
}

// AppRequestTestBatch APP请求测试批次表
type AppRequestTestBatch struct {
	ID           uint  `json:"id" gorm:"primarykey;comment:主键"`
	Status       uint  `json:"status" gorm:"index;comment:状态"`
	App          uint  `json:"app" gorm:"comment:app测试数"`
	Sum          uint  `json:"sum" gorm:"comment:累计测试数"`
	CreateTime   int64 `json:"create_time" gorm:"comment:创建时间"`
	EndTime      int64 `json:"end_time" gorm:"comment:结束时间"`
	SuccessCount uint  `json:"success_count" gorm:"comment:成功数"`
	FailureCount uint  `json:"failure_count" gorm:"comment:失败数"`
}

func (AppRequestTest) TableName() string      { return "app_request_tests_extend" }
func (AppRequestTestBatch) TableName() string { return "app_request_test_batches_extend" }
