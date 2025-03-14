"""
Natrue 文档工具提供者实现
"""
from typing import Dict, Any

from core.tools.provider.builtin_tool_provider import BuiltinToolProviderController


class NatrueToolProvider(BuiltinToolProviderController):
    """
    Natrue 文档工具提供者
    
    提供文档和表格生成工具
    """
    
    def _validate_credentials(self, credentials: Dict[str, Any]) -> None:
        """
        验证凭证
        
        Args:
            credentials: 凭证字典
        """
        # 验证 Dify API Key
        dify_api_key = credentials.get("dify_api_key")
        if not dify_api_key:
            raise ValueError("Dify API Key 未提供")
            
        # 可以在这里添加更多的验证逻辑，例如测试 API 连接 