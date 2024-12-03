from flask import Blueprint

from libs.external_api import ExternalApi

from .app.app_import import AppImportApi, AppImportConfirmApi
from .files import FileApi, FilePreviewApi, FileSupportTypeApi
from .remote_files import RemoteFileInfoApi, RemoteFileUploadApi

bp = Blueprint("console", __name__, url_prefix="/console/api")
api = ExternalApi(bp)

# File
api.add_resource(FileApi, "/files/upload")
api.add_resource(FilePreviewApi, "/files/<uuid:file_id>/preview")
api.add_resource(FileSupportTypeApi, "/files/support-type")

# Remote files
api.add_resource(RemoteFileInfoApi, "/remote-files/<path:url>")
api.add_resource(RemoteFileUploadApi, "/remote-files/upload")

# Import App
api.add_resource(AppImportApi, "/apps/imports")
api.add_resource(AppImportConfirmApi, "/apps/imports/<string:import_id>/confirm")

# Import other controllers
from . import admin, apikey, extension, feature, ping, setup, version

# Import app controllers
from .app import (
    advanced_prompt_template,
    agent,
    annotation,
    app,
    app_extend, # 二开部分：新增同步应用到模版中心
    audio,
    completion,
    conversation,
    conversation_variables,
    ding_talk_extend,  # Extend: DingTalk Related APIs
    generator,
    message,
    model_config,
    ops_trace,
    passport_extend,  # 二开部分: 新增passport_extend(额度限制，应用web计费)
    site,
    statistic,
    workflow,
    workflow_app_log,
    workflow_run,
    workflow_statistic,
)

# Import auth controllers
from .auth import activate, data_source_bearer_auth, data_source_oauth, forgot_password, login, oauth, register_extend # 二开部分: 新增用户（调用dify注册接口）

# Import billing controllers
from .billing import billing

# Import datasets controllers
from .datasets import (
    data_source,
    datasets,
    datasets_document,
    datasets_segments,
    external,
    hit_testing,
    website,
)

# Import explore controllers
from .explore import (
    audio,
    completion,
    conversation,
    installed_app,
    message,
    parameter,
    recommended_app,
    saved_message,
    workflow,
)

# Import tag controllers
from .tag import tags

# Import workspace controllers
# 二开部分：新增models_extend,account_extend,model_providers_extend
from .workspace import (
    account,
    account_extend,
    load_balancing_config,
    members,
    model_providers,
    model_providers_extend,
    models,
    models_extend,
    tool_providers,
    workspace
)
