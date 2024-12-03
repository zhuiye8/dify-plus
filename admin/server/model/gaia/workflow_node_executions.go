package gaia

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type WorkflowNodeExecutions struct {
	ID                uuid.UUID  `gorm:"column:id;primaryKey" json:"id"`                        // 唯一标识
	TenantID          uuid.UUID  `gorm:"column:tenant_id" json:"tenant_id"`                     // 工作空间ID
	AppID             uuid.UUID  `gorm:"column:app_id" json:"app_id"`                           // 应用ID
	WorkflowID        uuid.UUID  `gorm:"column:workflow_id" json:"workflow_id"`                 // 工作流ID
	TriggeredFrom     string     `gorm:"column:triggered_from" json:"triggered_from"`           // 触发来源
	WorkflowRunID     *uuid.UUID `gorm:"column:workflow_run_id" json:"workflow_run_id"`         // 工作流运行ID
	Index             int        `gorm:"column:index" json:"index"`                             // 索引
	PredecessorNodeID *string    `gorm:"column:predecessor_node_id" json:"predecessor_node_id"` // 前置节点ID
	NodeID            string     `gorm:"column:node_id" json:"node_id"`                         // 节点ID
	NodeType          string     `gorm:"column:node_type" json:"node_type"`                     // 节点类型
	Title             string     `gorm:"column:title" json:"title"`                             // 标题
	Inputs            *string    `gorm:"column:inputs" json:"inputs"`                           // 输入数据
	ProcessData       *string    `gorm:"column:process_data" json:"process_data"`               // 处理数据
	Outputs           *string    `gorm:"column:outputs" json:"outputs"`                         // 输出数据
	Status            string     `gorm:"column:status" json:"status"`                           // 状态
	Error             *string    `gorm:"column:error" json:"error"`                             // 错误信息
	ElapsedTime       float64    `gorm:"column:elapsed_time" json:"elapsed_time"`               // 执行耗时
	ExecutionMetadata *string    `gorm:"column:execution_metadata" json:"execution_metadata"`   // 执行元数据
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`                   // 创建时间
	CreatedByRole     string     `gorm:"column:created_by_role" json:"created_by_role"`         // 创建者角色
	CreatedBy         uuid.UUID  `gorm:"column:created_by" json:"created_by"`                   // 创建者ID
	FinishedAt        *time.Time `gorm:"column:finished_at" json:"finished_at"`                 // 完成时间
	NodeExecutionID   *string    `gorm:"column:node_execution_id" json:"node_execution_id"`     // 节点执行ID
}
