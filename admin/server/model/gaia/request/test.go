package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

const GetAppRequestFilterSuccess = 1                                       // 筛选成功
const GetAppRequestFilterFailure = 2                                       // 筛选失败
const PostgreSQLDataLimit = 1000                                           // 查询数据限制
const PostgreSQLDataTypeUUID = "uuid"                                      // uuid类型
const PostgreSQLDataTypeCharacterVarying = "character varying"             // 可变字符类型
const PostgreSQLDataTypeText = "text"                                      // 文本类型
const PostgreSQLDataTypeJSON = "json"                                      // JSON类型
const PostgreSQLDataTypeInteger = "integer"                                // 整数类型
const PostgreSQLDataTypeDoublePrecision = "double precision"               // 双精度浮点数类型
const PostgreSQLDataTypeNumeric = "numeric"                                // 数值类型
const PostgreSQLDataTypeTimestampWithoutTZ = "timestamp without time zone" // 不带时区的时间戳类型
const PostgreSQLDataTypeBoolean = "boolean"                                // 布尔类型
const PostgreSQLDefaultSchema = "public"                                   // 默认环境

// SyncDatabaseTableData 同步数据库表数据
type SyncDatabaseTableData struct {
	LogTable  string `json:"log_table" form:"log_table" gorm:"comment:旧表"`
	NewTable  string `json:"new_table" form:"new_table" gorm:"comment:要同步数据到的旧表"`
	KeyName   string `json:"key_name" form:"key_name" gorm:"comment:表主键名"`
	OrderName string `json:"order_name" form:"order_name" gorm:"comment:排序索引名"`
	GroupName string `json:"group_name" form:"group_name" gorm:"comment:表分组名(可为空)"`
}

type GetAppRequestTestRequest struct {
	request.PageInfo
	Apps    []string `json:"apps[]" form:"apps[]" gorm:"comment:检索app"`
	Status  uint     `json:"status" form:"status" gorm:"index;comment:状态"`
	BatchId uint     `json:"batch_id" form:"batch_id" gorm:"index;comment:批次ID"`
}

type DatabaseTableColumn struct {
	ColumnName string `json:"column_name" form:"column_name" gorm:"comment:列名"`
	DataType   string `json:"data_type" form:"data_type" gorm:"comment:数据类型"`
	IsNullable bool   `json:"is_nullable" form:"is_nullable" gorm:"comment:是否为空"`
}
