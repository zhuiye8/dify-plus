package gaia

import "time"

const WorkflowSucceeded = "succeeded" // 工作流状态成功
const WorkflowRunning = "running"     // 工作流状态运行中

// WorkflowRun 工作流运行
type WorkflowRun struct {
	ID             string    `json:"id" gorm:"index;comment:工作流运行ID"`
	TenantID       string    `json:"tenant_id" gorm:"not null;comment:租户ID"`
	AppID          string    `json:"app_id" gorm:"not null;comment:应用ID"`
	SequenceNumber int       `json:"sequence_number" gorm:"not null;comment:序列号"`
	WorkflowID     string    `json:"workflow_id" gorm:"not null;comment:工作流ID"`
	Type           string    `json:"type" gorm:"not null;comment:类型"`
	TriggeredFrom  string    `json:"triggered_from" gorm:"not null;comment:触发来源"`
	Version        string    `json:"version" gorm:"not null;comment:版本"`
	Graph          string    `json:"graph" gorm:"comment:图形表示"`
	Inputs         string    `json:"inputs" gorm:"comment:输入"`
	Status         string    `json:"status" gorm:"not null;comment:状态"`
	Outputs        string    `json:"outputs" gorm:"comment:输出"`
	Error          string    `json:"error" gorm:"comment:错误信息"`
	ElapsedTime    float64   `json:"elapsed_time" gorm:"not null;default:0;comment:耗时"`
	TotalTokens    int       `json:"total_tokens" gorm:"not null;default:0;comment:总令牌数"`
	TotalSteps     int       `json:"total_steps" gorm:"default:0;comment:总步骤数"`
	CreatedByRole  string    `json:"created_by_role" gorm:"not null;comment:创建者角色"`
	CreatedBy      string    `json:"created_by" gorm:"not null;comment:创建者ID"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP(0);comment:创建时间"`
	FinishedAt     time.Time `json:"finished_at" gorm:"comment:完成时间"`
}

func (WorkflowRun) TableName() string { return "workflow_runs" }
