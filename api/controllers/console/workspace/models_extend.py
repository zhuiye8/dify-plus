from flask_login import current_user
from flask_restful import Resource, inputs, marshal, reqparse
from werkzeug.exceptions import Forbidden

from configs import dify_config
from controllers.console import api
from controllers.console.workspace.workspace import setup_required, tenants_fields
from controllers.console.wraps import account_initialization_required
from core.model_runtime.entities.model_entities import ModelType
from extensions.ext_database import db
from libs.login import login_required
from models.provider import ProviderModel
from models.tenant_model_sync_extend import ModelSyncConfigExtend
from services.account_service_extend import TenantExtendService
from services.model_provider_service import ModelProviderService
from services.model_provider_service_extend import ModelProviderExtendService


class ModelProviderModelSyncApi(Resource):
    @setup_required
    @login_required
    @account_initialization_required
    def post(self, provider: str):

        tenant_extend_service = TenantExtendService
        super_admin_tenant_id = tenant_extend_service.get_super_admin_tenant_id().id

        if super_admin_tenant_id is None:
            return {"error": "请设置默认工作区ID"}, 400

        parser = reqparse.RequestParser()
        parser.add_argument("model", type=str, required=True, nullable=False, location="json")
        parser.add_argument(
            "model_type",
            type=str,
            required=True,
            nullable=False,
            choices=[mt.value for mt in ModelType],
            location="json",
        )
        parser.add_argument("tenant_ids", type=str, required=False, location="json", action="append")
        parser.add_argument("is_all", type=inputs.boolean)
        args = parser.parse_args()

        model = args["model"]
        model_type = args["model_type"]
        request_tenant_ids = args["tenant_ids"]
        is_all = args["is_all"]

        model_provider_service_extend = ModelProviderExtendService()

        # 查询原始模型记录
        provider_model_record = (
            db.session.query(ProviderModel)
            .filter(
                ProviderModel.tenant_id == super_admin_tenant_id,
                ProviderModel.provider_name == provider,
                ProviderModel.model_name == model,
                ProviderModel.model_type == ModelType.value_of(model_type).to_origin_model_type(),
            )
            .first()
        )

        # TODO 以下逻辑到时候丢到异步队列里去
        tenant_extend_service = TenantExtendService()
        if is_all:
            # 获取所有工作空间
            tenants = tenant_extend_service.get_all_tenants()

            all_tenant_ids = []
            for tenant in tenants:
                all_tenant_ids.append(tenant.id)
            tenant_ids = all_tenant_ids

            tenant_extend_service.create_model_sync_config_if_not_exist(provider_model_record.id)
        else:
            tenant_extend_service.delete_model_sync_config(provider_model_record.id)
            # 根据传过来的工作空间ID
            tenant_ids = request_tenant_ids

        # 获取当前已同步的工作区间ID
        sync_tenant_ids = []
        sync_tenants = model_provider_service_extend.get_current_syned_tenants(provider_model_record.id)
        for tenant_syned in sync_tenants:
            sync_tenant_ids.append(tenant_syned.tenant_id)

        # 如果 tenant_ids 为空，则删除所有当前同步的工作区
        if not tenant_ids:
            to_delete = set(sync_tenant_ids)
            to_add = set()
        else:
            # 找出需要删除和新增的工作区ID
            to_delete = set(sync_tenant_ids) - set(tenant_ids)
            to_add = set(tenant_ids) - set(sync_tenant_ids)

        # 删除不再需要同步的工作区
        model_provider_service = ModelProviderService()
        for tenant_id in to_delete:
            model_provider_service.remove_model_credentials(
                tenant_id=tenant_id, provider=provider, model=model, model_type=model_type
            )
            model_provider_service_extend.delete_syned_tenants(origin_model_id=provider_model_record.id, tenant_id=tenant_id)

        # 只对需要新增的工作区进行操作
        for tenant_id in to_add:
            if tenant_id == super_admin_tenant_id:
                continue
            origin_credentials = model_provider_service_extend.get_model_credentials_obfuscated(
                tenant_id=super_admin_tenant_id, provider=provider, model_type=model_type, model=model
            )
            # 查不到相应的凭证
            if origin_credentials is None:
                raise ValueError("Credentials cannot be None")

            # 保存模型数据
            model_id = model_provider_service_extend.save_model_credentials_without_validate(
                tenant_id=tenant_id,
                provider=provider,
                model_type=model_type,
                model=model,
                credentials=origin_credentials,
            )

            # 创建模型同步关系
            model_provider_service_extend.create_tenant_model_sync_if_not_exist(
                tenant_id=tenant_id, model_id=model_id, origin_model_id=provider_model_record.id
            )

        return {"result": "success"}


class ModelProviderModelSyncWorkspacesApi(Resource):
    @setup_required
    @login_required
    @account_initialization_required
    def get(self, provider: str):
        if not current_user.is_admin_or_owner:
            raise Forbidden()

        parser = reqparse.RequestParser()
        parser.add_argument("model", type=str, required=True, nullable=False, location="args")
        parser.add_argument(
            "model_type",
            type=str,
            required=True,
            nullable=False,
            choices=[mt.value for mt in ModelType],
            location="args",
        )
        args = parser.parse_args()

        model = args["model"]
        model_type = args["model_type"]

        # 查询原始模型记录
        provider_model_record = (
            db.session.query(ProviderModel)
            .filter(
                ProviderModel.tenant_id == current_user.current_tenant_id,
                ProviderModel.provider_name == provider,
                ProviderModel.model_name == model,
                ProviderModel.model_type == ModelType.value_of(model_type).to_origin_model_type(),
            )
            .first()
        )

        is_all = False
        if provider_model_record is None:
            return {"workspaces": [], "is_all": is_all}, 200

        tenants = TenantExtendService.get_model_sync_join_tenants(
            origin_model_id=provider_model_record.id, current_role=current_user.current_role, account_id=current_user.id
        )
        if not tenants:
            return {"workspaces": [], "is_all": is_all}, 200

        # is_all 获取
        model_sync_config_record = (
            db.session.query(ModelSyncConfigExtend)
            .filter(ModelSyncConfigExtend.model_id == provider_model_record.id)
            .first()
        )
        if model_sync_config_record:
            is_all = True

        return {"workspaces": marshal(tenants, tenants_fields), "is_all": is_all}, 200


api.add_resource(
    ModelProviderModelSyncWorkspacesApi,
    "/workspaces/current/model-providers/<string:provider>/models/get-model-sync-workspaces",
    endpoint="model-provider-model-sync-workspaces",
)
api.add_resource(
    ModelProviderModelSyncApi,
    "/workspaces/current/model-providers/<string:provider>/models/sync",
    endpoint="model-provider-model-sync",
)
