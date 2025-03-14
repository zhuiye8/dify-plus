"""
Nature 文档工具提供者实现
"""
from typing import Dict, Any

from core.tools.provider.builtin_tool_provider import BuiltinToolProviderController


class NatureToolProvider(BuiltinToolProviderController):
    """
    Nature 文档工具提供者
    
    提供文档和表格生成工具
    """
    
    def _validate_credentials(self, credentials: Dict[str, Any]) -> None:
        """
        验证凭证
        
        Args:
            credentials: 凭证字典
        """
        # 不需要验证凭证
        pass 