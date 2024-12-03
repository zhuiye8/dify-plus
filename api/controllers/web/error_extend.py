from libs.exception import BaseHTTPException


class AccountNoMoneyErrorExtend(BaseHTTPException):
    error_code = "Insufficient balance, call failed."
    description = "余额不足，调用失败！"
    code = 403


class WebAuthRequiredErrorExtend(BaseHTTPException):
    error_code = "unauthorized"
    description = "Unauthorized"
    code = 401
