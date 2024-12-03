import logging

from flask_restful import reqparse
from werkzeug.exceptions import InternalServerError

from controllers.web import api
from controllers.web.completion import is_end_login, is_money_limit  # You must log in to access your account extend
from controllers.web.error import (
    CompletionRequestError,
    NotWorkflowAppError,
    ProviderModelCurrentlyNotSupportError,
    ProviderNotInitializeError,
    ProviderQuotaExceededError,
)
from controllers.web.error_extend import (
    AccountNoMoneyErrorExtend,  # You must log in to access your account extend
    WebAuthRequiredErrorExtend,
)
from controllers.web.wraps import WebApiResource
from core.app.apps.base_app_queue_manager import AppQueueManager
from core.app.entities.app_invoke_entities import InvokeFrom
from core.errors.error import ModelCurrentlyNotSupportError, ProviderTokenNotInitError, QuotaExceededError
from core.model_runtime.errors.invoke import InvokeError
from libs import helper
from models.model import App, AppMode, EndUser
from services.app_generate_service import AppGenerateService
from services.app_generate_service_extend import (
    AppGenerateServiceExtend,  # Extend: App Center - Recommended list sorted by usage frequency
)

logger = logging.getLogger(__name__)


class WorkflowRunApi(WebApiResource):
    def post(self, app_model: App, end_user: EndUser):
        """
        Run workflow
        """
        app_mode = AppMode.value_of(app_model.mode)
        if app_mode != AppMode.WORKFLOW:
            raise NotWorkflowAppError()

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
        parser.add_argument("inputs", type=dict, required=True, nullable=False, location="json")
        parser.add_argument("files", type=list, required=False, location="json")
        args = parser.parse_args()

        try:
            AppGenerateServiceExtend.calculate_cumulative_usage(
                app_model=app_model,
                args=args,
            )  # Extend: App
            # Center - Recommended list sorted by usage frequency
            response = AppGenerateService.generate(
                app_model=app_model, user=end_user, args=args, invoke_from=InvokeFrom.WEB_APP, streaming=True
            )

            return helper.compact_generate_response(response)
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


class WorkflowTaskStopApi(WebApiResource):
    def post(self, app_model: App, end_user: EndUser, task_id: str):
        """
        Stop workflow task
        """
        app_mode = AppMode.value_of(app_model.mode)
        if app_mode != AppMode.WORKFLOW:
            raise NotWorkflowAppError()

        AppQueueManager.set_stop_flag(task_id, InvokeFrom.WEB_APP, end_user.id)

        return {"result": "success"}


api.add_resource(WorkflowRunApi, "/workflows/run")
api.add_resource(WorkflowTaskStopApi, "/workflows/tasks/<string:task_id>/stop")
