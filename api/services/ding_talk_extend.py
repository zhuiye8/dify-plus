import json
import logging
import time

import requests
from alibabacloud_dingtalk.oauth2_1_0 import models as dingtalkoauth_2__1__0_models
from alibabacloud_dingtalk.oauth2_1_0.client import Client as dingtalkoauth2_1_0Client
from alibabacloud_tea_openapi import models as open_api_models
from alibabacloud_tea_util.client import Client as UtilClient
from flask import request

from configs import dify_config
from extensions.ext_database import db
from libs.helper import extract_remote_ip
from models.account import Account, AccountIntegrate
from services.account_service import AccountService, RegisterService, TenantService
from services.account_service_extend import TenantExtendService

logger = logging.getLogger(__name__)

ACCESS_TOKEN = {
    "time": 0,
    "token": "",
}


class DingTalkService:
    @classmethod
    def create_client(cls) -> dingtalkoauth2_1_0Client:
        """
        使用 Token 初始化账号Client
        @return: Client
        @throws Exception
        """
        config = open_api_models.Config()
        config.protocol = "https"
        config.region_id = "central"
        return dingtalkoauth2_1_0Client(config)

    @classmethod
    def get_access_token(cls) -> (str, str):
        global ACCESS_TOKEN
        if ACCESS_TOKEN["time"] > time.time():
            return ACCESS_TOKEN["token"], ""
        # get token
        client = cls.create_client()
        get_access_token_request = dingtalkoauth_2__1__0_models.GetAccessTokenRequest(
            app_secret=dify_config.NEXT_PUBLIC_DING_TALK_CORP_SECRET,
            app_key=dify_config.NEXT_PUBLIC_DING_TALK_CORP_ID,
        )
        try:
            token_request = client.get_access_token(get_access_token_request)
            if token_request.status_code == 200:
                ACCESS_TOKEN["token"] = token_request.body.access_token
                ACCESS_TOKEN["time"] = int(time.time()) + 3600
                return token_request.body.access_token, ""
            else:
                return "", token_request.body
        except Exception as err:
            if not UtilClient.empty(err.code) and not UtilClient.empty(err.message):
                # err 中含有 code 和 message 属性，可帮助开发定位问题
                return "", f"Failed to retrieve token:${err.code}, {err.message}"
            return "", "Failed to retrieve token"

    @classmethod
    def get_user_info(cls, code: str) -> (str, str):
        oa = dify_config.OAUTH2_CLIENT_URL
        host = "https://oapi.dingtalk.com/topapi/v2/user"
        token, err = cls.get_access_token()
        if err != "":
            return "", f"Failed to obtain token: {err}"
        response = requests.post(
            f"{host}/getuserinfo?access_token={token}",
            json={
                "code": code,
            },
        )
        # Check the response status code
        if response.status_code != 200:
            return "", f"Request failed, status code: {response.status_code}"
        # Print the response content
        req = response.json()
        if req["errcode"] != 0:
            return "", "Request failed: " + req["errmsg"]
        responses = requests.post(
            f'{oa}/serv/?c=user&a=getUserInfoByUserId&userId={req["result"]["userid"]}',
            json={
                "userid": req["result"]["userid"],
            },
        )
        # Check the response status code
        if responses.status_code != 200:
            return "", f"Request for user information failed, status code: {responses.status_code}"
        reqs = responses.json()
        if reqs["code"] != 200 and len(reqs["data"]) > 0:
            return "", "Request for user information failed: " + json.dumps(req) + " " + json.dumps(reqs)
        # Check if the user exists
        user = reqs["data"][0]
        integrate: AccountIntegrate = (
            db.session.query(AccountIntegrate).filter(AccountIntegrate.open_id == user["userName"]).first()
        )
        if integrate is None:
            # registered user
            try:
                account = RegisterService.register(
                    email=f"{user['userName']}@{dify_config.EMAIL_DOMAIN}",
                    open_id=user["userName"],
                    name=user["name"],
                    provider="oauth2",
                    password=None,
                )
            except EOFError as a:
                logging.info(f"register user error: {str(a)}， info {json.loads(user)}")
                return "", "register error"

            tenant_extend_service = TenantExtendService
            super_admin_id = tenant_extend_service.get_super_admin_id().id
            super_admin_tenant_id = tenant_extend_service.get_super_admin_tenant_id().id
            if super_admin_id and super_admin_tenant_id:
                isCreate = TenantExtendService.create_default_tenant_member_if_not_exist(
                    super_admin_tenant_id, account.id
                )
                if isCreate:
                    TenantService.switch_tenant(account, super_admin_tenant_id)
        else:
            account: Account = db.session.query(Account).filter(Account.id == integrate.account_id).first()
        # token jwt
        token = AccountService.login(account, ip_address=extract_remote_ip(request))

        return f"{dify_config.CONSOLE_WEB_URL}/explore/apps-center-extend?console_token={token.access_token}&&refresh_token={token.refresh_token}", ""
