import logging  # 二开部分，针对oa登录报错问题，记录返回的code
import urllib.parse
from dataclasses import dataclass
from typing import Optional

import requests
from flask import current_app


@dataclass
class OAuthUserInfo:
    id: str
    name: str
    email: str


class OAuth:
    def __init__(self, client_id: str, client_secret: str, redirect_uri: str):
        self.client_id = client_id
        self.client_secret = client_secret
        self.redirect_uri = redirect_uri

    def get_authorization_url(self):
        raise NotImplementedError()

    def get_access_token(self, code: str):
        raise NotImplementedError()

    def get_raw_user_info(self, token: str):
        raise NotImplementedError()

    def get_user_info(self, token: str) -> OAuthUserInfo:
        raw_info = self.get_raw_user_info(token)
        return self._transform_user_info(raw_info)

    def _transform_user_info(self, raw_info: dict) -> OAuthUserInfo:
        raise NotImplementedError()


class GitHubOAuth(OAuth):
    _AUTH_URL = "https://github.com/login/oauth/authorize"
    _TOKEN_URL = "https://github.com/login/oauth/access_token"
    _USER_INFO_URL = "https://api.github.com/user"
    _EMAIL_INFO_URL = "https://api.github.com/user/emails"

    def get_authorization_url(self, invite_token: Optional[str] = None):
        params = {
            "client_id": self.client_id,
            "redirect_uri": self.redirect_uri,
            "scope": "user:email",  # Request only basic user information
        }
        if invite_token:
            params["state"] = invite_token
        return f"{self._AUTH_URL}?{urllib.parse.urlencode(params)}"

    def get_access_token(self, code: str):
        data = {
            "client_id": self.client_id,
            "client_secret": self.client_secret,
            "code": code,
            "redirect_uri": self.redirect_uri,
        }
        headers = {"Accept": "application/json"}
        response = requests.post(self._TOKEN_URL, data=data, headers=headers)

        response_json = response.json()
        access_token = response_json.get("access_token")

        if not access_token:
            raise ValueError(f"Error in GitHub OAuth: {response_json}")

        return access_token

    def get_raw_user_info(self, token: str):
        headers = {"Authorization": f"token {token}"}
        response = requests.get(self._USER_INFO_URL, headers=headers)
        response.raise_for_status()
        user_info = response.json()

        email_response = requests.get(self._EMAIL_INFO_URL, headers=headers)
        email_info = email_response.json()
        primary_email = next((email for email in email_info if email["primary"] == True), None)

        return {**user_info, "email": primary_email["email"]}

    def _transform_user_info(self, raw_info: dict) -> OAuthUserInfo:
        email = raw_info.get("email")
        if not email:
            email = f"{raw_info['id']}+{raw_info['login']}@users.noreply.github.com"
        return OAuthUserInfo(id=str(raw_info["id"]), name=raw_info["name"], email=email)


class GoogleOAuth(OAuth):
    _AUTH_URL = "https://accounts.google.com/o/oauth2/v2/auth"
    _TOKEN_URL = "https://oauth2.googleapis.com/token"
    _USER_INFO_URL = "https://www.googleapis.com/oauth2/v3/userinfo"

    def get_authorization_url(self, invite_token: Optional[str] = None):
        params = {
            "client_id": self.client_id,
            "response_type": "code",
            "redirect_uri": self.redirect_uri,
            "scope": "openid email",
        }
        if invite_token:
            params["state"] = invite_token
        return f"{self._AUTH_URL}?{urllib.parse.urlencode(params)}"

    def get_access_token(self, code: str):
        data = {
            "client_id": self.client_id,
            "client_secret": self.client_secret,
            "code": code,
            "grant_type": "authorization_code",
            "redirect_uri": self.redirect_uri,
        }
        headers = {"Accept": "application/json"}
        response = requests.post(self._TOKEN_URL, data=data, headers=headers)

        response_json = response.json()
        access_token = response_json.get("access_token")

        if not access_token:
            raise ValueError(f"Error in Google OAuth: {response_json}")

        return access_token

    def get_raw_user_info(self, token: str):
        headers = {"Authorization": f"Bearer {token}"}
        response = requests.get(self._USER_INFO_URL, headers=headers)
        response.raise_for_status()
        return response.json()

    def _transform_user_info(self, raw_info: dict) -> OAuthUserInfo:
        return OAuthUserInfo(id=str(raw_info["sub"]), name=None, email=raw_info["email"])


class OaOAuth(OAuth):
    _AUTH_URL = ""
    _Host = ""
    _TOKEN_URL = ""
    _USER_INFO_URL = ""

    def get_authorization_url(self, invite_token: Optional[str] = None):
        params = {
            "client_id": self.client_id,
            "redirect_uri": self.redirect_uri,
        }
        with current_app.app_context():
            self._Host = current_app.config.get("OAUTH2_CLIENT_URL")
            self._TOKEN_URL = current_app.config.get("OAUTH2_TOKEN_URL")
            self._USER_INFO_URL = current_app.config.get("OAUTH2_USER_INFO_URL")
        return f"{self._Host}{self._AUTH_URL}?{urllib.parse.urlencode(params)}"

    def get_access_token(self, code: str):
        data = {
            "client_id": self.client_id,
            "client_secret": self.client_secret,
            "code": code,
            "grant_type": "authorization_code",
            "redirect_uri": self.redirect_uri,
        }
        with current_app.app_context():
            self._Host = current_app.config.get("OAUTH2_CLIENT_URL")
        headers = {"Accept": "application/json"}
        response = requests.post(self._Host + self._TOKEN_URL, data=data, headers=headers)
        response.encoding = "utf-8"
        if response.status_code != 200:
            return ""
        try:
            response_json = response.json()
        except:
            return ""
        access_token = response_json.get("access_token")

        return access_token

    def get_raw_user_info(self, token: str):
        with current_app.app_context():
            self._Host = current_app.config.get("OAUTH2_CLIENT_URL")
        headers = {"Authorization": f"Bearer {token}"}
        response = requests.get(self._Host + self._USER_INFO_URL, headers=headers)
        response.raise_for_status()
        return response.json()

    def _transform_user_info(self, raw_info: dict) -> OAuthUserInfo:
        # 检查 raw_info 是否为空或为 None
        if not raw_info or not isinstance(raw_info, dict):
            return OAuthUserInfo(
                id="",
                name="",
                email="",
            )

        # 如果data为空说明报错了
        data = raw_info.get("data")
        if data == {} or data is None or not isinstance(data, dict):
            code = raw_info.get("code", "")
            msg = raw_info.get("info", "")
            logging.info(f"raw_info {raw_info}")
            return OAuthUserInfo(
                id="",
                name="",
                email="",
            )

        username = data.get("username")
        name = data.get("name")
        email = data.get("email")
        if not username:
            raise ValueError("OA系统返回用户数据格式不正确。请返回进行重新登录。")

        return OAuthUserInfo(
            id=str(username) if username is not None else None,
            name=str(name) if name is not None else None,
            email=str(email) if email is not None else None,
        )
