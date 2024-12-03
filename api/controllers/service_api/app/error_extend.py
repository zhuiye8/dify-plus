from libs.exception import BaseHTTPException


class AccountNoMoneyErrorExtend(BaseHTTPException):
    error_code = "Insufficient balance, call failed."
    description = "余额不足，调用失败！"
    code = 403


class ApiTokenDayNoMoneyErrorExtend(BaseHTTPException):
    error_code = "The daily call limit for this key has been reached, the call failed!"
    description = "该密钥每日调用额度已达上限，调用失败！"
    code = 403


class ApiTokenMonthNoMoneyErrorExtend(BaseHTTPException):
    error_code = "The monthly call limit for this key has been reached, the call failed!"
    description = "该密钥每月调用额度已达上限，调用失败！"
    code = 403
