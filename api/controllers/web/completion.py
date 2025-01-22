import logging

from flask import request  # ----------------- start You must log in to access your account extend ---------------
from flask_restful import reqparse  # type: ignore
from werkzeug.exceptions import InternalServerError, NotFound

import services
from controllers.web import api
from controllers.web.error import (
    AppUnavailableError,
    CompletionRequestError,
    ConversationCompletedError,
    NotChatAppError,
    NotCompletionAppError,
    ProviderModelCurrentlyNotSupportError,
    ProviderNotInitializeError,
    ProviderQuotaExceededError,
)
from controllers.web.error import InvokeRateLimitError as InvokeRateLimitHttpError
from controllers.web.error_extend import (
    AccountNoMoneyErrorExtend,  # You must log in to access your account extend
    WebAuthRequiredErrorExtend,
)
from controllers.web.wraps import WebApiResource
from core.app.apps.base_app_queue_manager import AppQueueManager
from core.app.entities.app_invoke_entities import InvokeFrom
from core.errors.error import (
    ModelCurrentlyNotSupportError,
    ProviderTokenNotInitError,
    QuotaExceededError,
)
from core.model_runtime.errors.invoke import InvokeError
from extensions.ext_database import db
from libs import helper
from libs.helper import uuid_value
from libs.passport import (
    PassportService,  # ----------------- start You must log in to access your account extend ---------------
)
from models.account_money_extend import AccountMoneyExtend
from models.model import AppMode
from services.account_service import (
    AccountService,  # ----------------- start You must log in to access your account extend ---------------
)
from services.app_generate_service import AppGenerateService
from services.app_generate_service_extend import AppGenerateServiceExtend
from services.errors.llm import InvokeRateLimitError


# ----------------- start You must log in to access your account extend ---------------
def is_end_login(end_user):
    user_info = None
    try:
        auth_token = request.headers.get("Authorization-extend")
        decoded = PassportService().verify(auth_token)
        user_info = AccountService.load_logged_in_account(account_id=decoded.get("user_id"))
        if user_info is not None:
            if end_user.external_user_id is None:
                end_user.external_user_id = decoded.get("user_id")
    except:
        logging.error("load_logged_in_account error")
        pass
    # no login
    return user_info


# ----------------- stop You must log in to access your account extend ---------------


# ----------------- 二开部分Begin - 额度限制 ---------------
def is_money_limit(end_user) -> bool:
    try:
        # TODO 需要写入缓存，读缓存
        account_money = (
            db.session.query(AccountMoneyExtend).filter(AccountMoneyExtend.account_id == end_user.id).first()
        )
        if not account_money:
            return False

        if account_money.used_quota >= account_money.total_quota:
            return True
        return False
    except:
        return True


# ----------------- 二开部分End - 额度限制  ---------------


# define completion api for user
class CompletionApi(WebApiResource):
    def post(self, app_model, end_user):
        if app_model.mode != "completion":
            raise NotCompletionAppError()

        # ----------------- start You must log in to access your account extend ---------------
        # no login
        if is_end_login(end_user) is None:
            raise WebAuthRequiredErrorExtend()
        # ----------------- stop You must log in to access your account extend ---------------

        # ----------------- 二开部分Begin - 余额判断-----------------
        if is_money_limit(end_user):
            raise AccountNoMoneyErrorExtend()
        # ----------------- 二开部分End - 余额判断-----------------

        parser = reqparse.RequestParser()
        parser.add_argument("inputs", type=dict, required=True, location="json")
        parser.add_argument("query", type=str, location="json", default="")
        parser.add_argument("files", type=list, required=False, location="json")
        parser.add_argument("response_mode", type=str, choices=["blocking", "streaming"], location="json")
        parser.add_argument("retriever_from", type=str, required=False, default="web_app", location="json")

        args = parser.parse_args()

        streaming = args["response_mode"] == "streaming"
        args["auto_generate_name"] = False

        try:
            AppGenerateServiceExtend.calculate_cumulative_usage(
                app_model=app_model,
                args=args,
            )  # Extend: App Center -
            # Recommended list sorted by usage frequency
            response = AppGenerateService.generate(
                app_model=app_model, user=end_user, args=args, invoke_from=InvokeFrom.WEB_APP, streaming=streaming
            )

            return helper.compact_generate_response(response)
        except services.errors.conversation.ConversationNotExistsError:
            raise NotFound("Conversation Not Exists.")
        except services.errors.conversation.ConversationCompletedError:
            raise ConversationCompletedError()
        except services.errors.app_model_config.AppModelConfigBrokenError:
            logging.exception("App model config broken.")
            raise AppUnavailableError()
        except ProviderTokenNotInitError as ex:
            raise ProviderNotInitializeError(ex.description)
        except QuotaExceededError:
            raise ProviderQuotaExceededError()
        except ModelCurrentlyNotSupportError:
            raise ProviderModelCurrentlyNotSupportError()
        except InvokeError as e:
            raise CompletionRequestError(e.description)
        except ValueError as e:
            raise e
        except Exception as e:
            logging.exception("internal server error.")
            raise InternalServerError()


class CompletionStopApi(WebApiResource):
    def post(self, app_model, end_user, task_id):
        if app_model.mode != "completion":
            raise NotCompletionAppError()

        AppQueueManager.set_stop_flag(task_id, InvokeFrom.WEB_APP, end_user.id)

        return {"result": "success"}, 200


class ChatApi(WebApiResource):
    def post(self, app_model, end_user):
        # ----------------- start You must log in to access your account extend ---------------
        # no login
        if is_end_login(end_user) is None:
            raise WebAuthRequiredErrorExtend()
        # ----------------- stop You must log in to access your account extend ---------------

        # ----------------- 二开部分Begin - 余额判断-----------------
        if is_money_limit(end_user):
            raise AccountNoMoneyErrorExtend()
        # ----------------- 二开部分End - 余额判断-----------------

        app_mode = AppMode.value_of(app_model.mode)
        if app_mode not in {AppMode.CHAT, AppMode.AGENT_CHAT, AppMode.ADVANCED_CHAT}:
            raise NotChatAppError()

        parser = reqparse.RequestParser()
        parser.add_argument("inputs", type=dict, required=True, location="json")
        parser.add_argument("query", type=str, required=True, location="json")
        parser.add_argument("files", type=list, required=False, location="json")
        parser.add_argument("files_extend", type=list, required=False, location="json")
        parser.add_argument("response_mode", type=str, choices=["blocking", "streaming"], location="json")
        parser.add_argument("conversation_id", type=uuid_value, location="json")
        parser.add_argument("parent_message_id", type=uuid_value, required=False, location="json")
        parser.add_argument("retriever_from", type=str, required=False, default="web_app", location="json")

        args = parser.parse_args()

        streaming = args["response_mode"] == "streaming"
        args["auto_generate_name"] = False

        try:
            AppGenerateServiceExtend.calculate_cumulative_usage(
                app_model=app_model,
                args=args,
            )  # Extend: App
            # Center - Recommended list sorted by usage frequency
            response = AppGenerateService.generate(
                app_model=app_model, user=end_user, args=args, invoke_from=InvokeFrom.WEB_APP, streaming=streaming
            )

            return helper.compact_generate_response(response)
        except services.errors.conversation.ConversationNotExistsError:
            raise NotFound("Conversation Not Exists.")
        except services.errors.conversation.ConversationCompletedError:
            raise ConversationCompletedError()
        except services.errors.app_model_config.AppModelConfigBrokenError:
            logging.exception("App model config broken.")
            raise AppUnavailableError()
        except ProviderTokenNotInitError as ex:
            raise ProviderNotInitializeError(ex.description)
        except QuotaExceededError:
            raise ProviderQuotaExceededError()
        except ModelCurrentlyNotSupportError:
            raise ProviderModelCurrentlyNotSupportError()
        except InvokeRateLimitError as ex:
            raise InvokeRateLimitHttpError(ex.description)
        except InvokeError as e:
            raise CompletionRequestError(e.description)
        except ValueError as e:
            raise e
        except Exception as e:
            logging.exception("internal server error.")
            raise InternalServerError()


class ChatStopApi(WebApiResource):
    def post(self, app_model, end_user, task_id):
        app_mode = AppMode.value_of(app_model.mode)
        if app_mode not in {AppMode.CHAT, AppMode.AGENT_CHAT, AppMode.ADVANCED_CHAT}:
            raise NotChatAppError()

        AppQueueManager.set_stop_flag(task_id, InvokeFrom.WEB_APP, end_user.id)

        return {"result": "success"}, 200


api.add_resource(CompletionApi, "/completion-messages")
api.add_resource(CompletionStopApi, "/completion-messages/<string:task_id>/stop")
api.add_resource(ChatApi, "/chat-messages")
api.add_resource(ChatStopApi, "/chat-messages/<string:task_id>/stop")
