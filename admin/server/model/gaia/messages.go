package gaia

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

const ChatRequestTypeApi = "api"   // 聊天请求类型为api
const MessagesSucceeded = "normal" // 聊天状态成功

type Messages struct {
	ID                      uuid.UUID `gorm:"column:id;primary_key;default:uuid_generate_v4()" json:"id"`                           // 消息ID
	AppID                   uuid.UUID `gorm:"column:app_id;not null" json:"app_id"`                                                 // 应用ID
	ModelProvider           string    `gorm:"column:model_provider" json:"model_provider"`                                          // 模型提供商
	ModelID                 string    `gorm:"column:model_id" json:"model_id"`                                                      // 模型ID
	OverrideModelConfigs    string    `gorm:"column:override_model_configs" json:"override_model_configs"`                          // 覆盖模型配置
	ConversationID          uuid.UUID `gorm:"column:conversation_id;not null" json:"conversation_id"`                               // 对话ID
	Inputs                  string    `gorm:"column:inputs" json:"inputs"`                                                          // 输入数据
	Query                   string    `gorm:"column:query;not null" json:"query"`                                                   // 查询内容
	Message                 string    `gorm:"column:message;not null" json:"message"`                                               // 消息内容
	MessageTokens           int       `gorm:"column:message_tokens;not null;default:0" json:"message_tokens"`                       // 消息令牌数
	MessageUnitPrice        float64   `gorm:"column:message_unit_price;not null" json:"message_unit_price"`                         // 消息单价
	Answer                  string    `gorm:"column:answer;not null" json:"answer"`                                                 // 回答内容
	AnswerTokens            int       `gorm:"column:answer_tokens;not null;default:0" json:"answer_tokens"`                         // 回答令牌数
	AnswerUnitPrice         float64   `gorm:"column:answer_unit_price;not null" json:"answer_unit_price"`                           // 回答单价
	ProviderResponseLatency float64   `gorm:"column:provider_response_latency;not null;default:0" json:"provider_response_latency"` // 提供商响应延迟
	TotalPrice              float64   `gorm:"column:total_price" json:"total_price"`                                                // 总价格
	Currency                string    `gorm:"column:currency;not null" json:"currency"`                                             // 货币
	FromSource              string    `gorm:"column:from_source;not null" json:"from_source"`                                       // 来源
	FromEndUserID           uuid.UUID `gorm:"column:from_end_user_id" json:"from_end_user_id"`                                      // 来自最终用户ID
	FromAccountID           uuid.UUID `gorm:"column:from_account_id" json:"from_account_id"`                                        // 来自账户ID
	CreatedAt               time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`               // 创建时间
	UpdatedAt               time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`               // 更新时间
	AgentBased              bool      `gorm:"column:agent_based;not null;default:false" json:"agent_based"`                         // 是否基于代理
	MessagePriceUnit        float64   `gorm:"column:message_price_unit;not null;default:0.001" json:"message_price_unit"`           // 消息价格单位
	AnswerPriceUnit         float64   `gorm:"column:answer_price_unit;not null;default:0.001" json:"answer_price_unit"`             // 回答价格单位
	WorkflowRunID           uuid.UUID `gorm:"column:workflow_run_id" json:"workflow_run_id"`                                        // 工作流运行ID
	Status                  string    `gorm:"column:status;not null;default:normal" json:"status"`                                  // 状态
	Error                   string    `gorm:"column:error" json:"error"`                                                            // 错误信息
	MessageMetadata         string    `gorm:"column:message_metadata" json:"message_metadata"`                                      // 消息元数据
	InvokeFrom              string    `gorm:"column:invoke_from" json:"invoke_from"`                                                // 调用来源
}

func (Messages) TableName() string {
	return "messages"
}
