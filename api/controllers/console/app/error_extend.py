from libs.exception import BaseHTTPException


class WebSSOAuthRequiredError(BaseHTTPException):
    error_code = "web_sso_auth_required"
    description = "Web SSO authentication required."
    code = 401


class DingTalkNotExist(BaseHTTPException):
    error_code = "login_not_exist"
    description = "DingTalk login failed."
    code = 400
