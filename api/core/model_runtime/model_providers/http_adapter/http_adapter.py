import logging

from core.model_runtime.entities.model_entities import ModelType
from core.model_runtime.errors.validate import CredentialsValidateFailedError
from core.model_runtime.model_providers.__base.model_provider import ModelProvider

logger = logging.getLogger(__name__)


class HTTPAdapterProvider(ModelProvider):
    def validate_provider_credentials(self, credentials: dict) -> None:
        """
        验证提供商凭据
        如果验证失败，抛出异常

        :param credentials: 提供商凭据，包含base_url和api_key
        """
        if not credentials.get("base_url"):
            raise CredentialsValidateFailedError("API 基础地址 (base_url) 是必填项")
            
        if not credentials.get("api_key"):
            raise CredentialsValidateFailedError("API 密钥 (api_key) 是必填项")
            
        try:
            # 获取LLM模型实例
            model_instance = self.get_model_instance(ModelType.LLM)
            
            # 准备验证凭据
            # 默认添加必要的配置项，减少用户输入
            enhanced_credentials = {
                **credentials,
                "function_calling_type": "tool_call",  # 默认支持函数调用
                "stream_function_calling": "supported",
                "vision_support": "not_support",  # 默认不支持多模态
                "mode": "chat",  # 默认聊天模式
                "context_size": "4096"  # 默认上下文长度
            }
            
            # 使用默认模型名称或用户提供的模型进行验证
            model = credentials.get("model", "deepseek-chat")
            
            # 执行验证
            model_instance.validate_credentials(model=model, credentials=enhanced_credentials)
        except CredentialsValidateFailedError as ex:
            raise ex
        except Exception as ex:
            logger.exception(f"{self.get_provider_schema().provider} 凭据验证失败")
            raise CredentialsValidateFailedError(f"验证失败: {str(ex)}") 