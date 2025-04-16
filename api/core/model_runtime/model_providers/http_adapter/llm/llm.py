import logging
import requests
import json
from collections.abc import Generator
from typing import Optional, Union
from urllib.parse import urljoin

from core.model_runtime.entities.llm_entities import LLMResult
from core.model_runtime.entities.message_entities import (
    PromptMessage,
    PromptMessageTool,
)
from core.model_runtime.errors.validate import CredentialsValidateFailedError
from core.model_runtime.model_providers.openai_api_compatible.llm.llm import OAIAPICompatLargeLanguageModel

logger = logging.getLogger(__name__)


class HTTPAdapterLargeLanguageModel(OAIAPICompatLargeLanguageModel):
    """
    HTTP适配器大语言模型类，默认支持DeepSeek API。
    """

    def _invoke(
        self,
        model: str,
        credentials: dict,
        prompt_messages: list[PromptMessage],
        model_parameters: dict,
        tools: Optional[list[PromptMessageTool]] = None,
        stop: Optional[list[str]] = None,
        stream: bool = True,
        user: Optional[str] = None,
    ) -> Union[LLMResult, Generator]:
        """
        调用LLM模型
        
        :param model: 模型ID
        :param credentials: 凭据信息
        :param prompt_messages: 提示消息
        :param model_parameters: 模型参数
        :param tools: 工具列表
        :param stop: 停止词
        :param stream: 是否流式输出
        :param user: 用户ID
        :return: LLM结果或生成器
        """
        # 预处理凭据
        self._preprocess_credentials(credentials, model)
        
        # 调用基类方法
        return super()._invoke(model, credentials, prompt_messages, model_parameters, tools, stop, stream, user)

    def validate_credentials(self, model: str, credentials: dict) -> None:
        """
        验证模型API凭据
        
        :param model: 模型名称，默认 deepseek-chat
        :param credentials: API凭据，包含base_url和api_key
        """
        # 确保必要的参数存在
        if not credentials.get("base_url"):
            raise CredentialsValidateFailedError("API 基础地址 (base_url) 是必填项")
        
        if not credentials.get("api_key"):
            raise CredentialsValidateFailedError("API 密钥 (api_key) 是必填项")
        
        # 处理api_key中可能的空格问题
        credentials["api_key"] = credentials.get("api_key", "").strip()
        
        # 预处理凭据
        self._preprocess_credentials(credentials, model)
        
        # 验证API连接
        self._validate_api_connection(credentials, model)
            
    def _validate_api_connection(self, credentials: dict, model: str) -> None:
        """
        验证API连接
        
        :param credentials: 预处理后的凭据
        :param model: 模型名称
        """
        # 获取基础URL和API路径
        base_url = credentials.get("endpoint_url", "").rstrip("/")
        api_key = credentials.get("api_key", "")
        
        # 构建/models端点URL
        models_url = urljoin(base_url, "/models")
        
        # 准备请求头
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {api_key}"
        }
        
        try:
            # 尝试调用/models端点
            logger.info(f"验证API连接: {models_url}")
            
            models_response = requests.get(
                models_url,
                headers=headers,
                timeout=10
            )
            
            if models_response.status_code == 200:
                logger.info("API验证成功: /models 端点可访问")
                # 尝试获取可用模型列表
                try:
                    models_data = models_response.json()
                    if models_data and 'data' in models_data and isinstance(models_data['data'], list):
                        available_models = [m.get('id') for m in models_data['data'] if isinstance(m, dict) and 'id' in m]
                        if available_models:
                            logger.info(f"可用模型: {', '.join(available_models)}")
                            credentials["available_models"] = available_models
                except Exception as e:
                    logger.warning(f"解析模型列表失败: {str(e)}")
                return
            
            # 如果/models端点失败，尝试chat/completions端点
            if models_response.status_code in [404, 401, 403, 500]:
                logger.warning(f"/models端点不可用，尝试 /chat/completions 端点...")
                
                # 构建chat/completions端点URL
                chat_url = urljoin(base_url, "/chat/completions")
                
                # 准备最小化请求体
                payload = {
                    "model": model,
                    "messages": [{"role": "user", "content": "Hi"}],
                    "max_tokens": 1,
                    "stream": False
                }
                
                # 发送请求
                chat_response = requests.post(
                    chat_url,
                    headers=headers,
                    json=payload,
                    timeout=10
                )
                
                if chat_response.status_code == 200:
                    logger.info("API验证成功: /chat/completions 端点可访问")
                    return
                
                # 如果chat/completions也失败，报告详细错误
                error_message = self._get_error_message(chat_response)
                raise CredentialsValidateFailedError(
                    f"API验证失败，状态码: {chat_response.status_code}。错误: {error_message}"
                )
            
            # 处理/models端点的错误
            error_message = self._get_error_message(models_response)
            raise CredentialsValidateFailedError(
                f"API验证失败，状态码: {models_response.status_code}。错误: {error_message}"
            )
                
        except requests.exceptions.RequestException as e:
            logger.error(f"API请求异常: {str(e)}", exc_info=True)
            raise CredentialsValidateFailedError(f"连接API失败: {str(e)}")
            
    def _get_error_message(self, response) -> str:
        """从响应中提取错误信息"""
        try:
            error_json = response.json()
            if isinstance(error_json, dict):
                # 首先尝试获取标准错误格式
                if 'error' in error_json:
                    error_content = error_json['error']
                    if isinstance(error_content, dict) and 'message' in error_content:
                        return error_content['message']
                    return str(error_content)
                    
                # 尝试其他可能的错误字段
                for field in ['message', 'detail', 'description']:
                    if field in error_json:
                        return str(error_json[field])
            
            # 如果无法解析结构化错误，返回整个JSON
            return json.dumps(error_json)
        except Exception:
            # 如果JSON解析失败，尝试返回文本内容
            try:
                return response.text[:200]  # 限制长度
            except Exception:
                return "无法获取错误详情"

    def _preprocess_credentials(self, credentials: dict, model: str) -> None:
        """
        预处理凭据，添加必要的默认值
        
        :param credentials: 原始凭据
        :param model: 模型名称
        """
        # 处理base_url
        base_url = credentials.get("base_url", "").rstrip("/")
        if not base_url.startswith(("http://", "https://")):
            base_url = f"https://{base_url}"
        
        # 设置endpoint_url，确保以/结尾
        credentials["endpoint_url"] = f"{base_url}/"
        
        # 设置模型名称
        if not credentials.get("model_name"):
            credentials["model_name"] = model
            
        # 默认参数设置
        if "function_calling_type" not in credentials:
            credentials["function_calling_type"] = "tool_call"
            
        if "stream_function_calling" not in credentials:
            credentials["stream_function_calling"] = "supported"
            
        if "vision_support" not in credentials:
            credentials["vision_support"] = "not_support"
            
        if "mode" not in credentials:
            credentials["mode"] = "chat"
            
        if "context_size" not in credentials:
            credentials["context_size"] = "4096" 