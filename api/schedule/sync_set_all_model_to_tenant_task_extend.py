import logging
from datetime import datetime

import click

import app
from extensions.ext_database import db
from extensions.ext_redis import redis_client
from models import Tenant
from services.account_service_extend import TenantExtendService
from services.model_provider_service_extend import ModelProviderExtendService


@app.celery.task(queue="extend_low")
def sync_set_all_model_to_tenant_task_extend():
    # 获取上次执行的时间
    last_execution_time_str = redis_client.get('sync_set_all_model_to_tenant_time')

    if last_execution_time_str is not None:
        last_execution_time = datetime.fromisoformat(last_execution_time_str.decode('utf-8'))
        tenants = db.session.query(Tenant).filter(Tenant.created_at > last_execution_time).order_by(Tenant.created_at.asc()).all()
    else:
        tenants = db.session.query(Tenant).order_by(Tenant.created_at.asc()).all()

    if tenants is None or tenants == []:
        logging.info("没有工作区需要同步")

    for tenant in tenants:
        tenant_id = tenant.id
        logging.info(f"开始同步所有模型到工作区: {tenant_id}")
        model_provider_service_extend = ModelProviderExtendService()
        # 同步供应商+模型名称的模型数据
        provider_model_records = TenantExtendService.get_sync_all_model()
        for provider_model_record in provider_model_records:
            logging.info(f"{tenant_id} 同步的模型所在工作区: {provider_model_record.tenant_id}")
            logging.info(f"{tenant_id} 同步的模型provider_name: {provider_model_record.provider_name}")
            logging.info(f"{tenant_id} 同步的模型model_type {provider_model_record.model_type}")
            logging.info(f"{tenant_id} 同步的模型model_name: {provider_model_record.model_name}")
            origin_credentials = model_provider_service_extend.get_model_credentials_obfuscated(
                tenant_id=provider_model_record.tenant_id,
                provider=provider_model_record.provider_name,
                model_type=provider_model_record.model_type,
                model=provider_model_record.model_name,
            )
            # 查不到相应的凭证
            if origin_credentials is None:
                logging.error(click.style(f"{tenant_id} 同步失败: {provider_model_record.provider_name}，{provider_model_record.model_name}", fg="red"))

                continue

            model_id = model_provider_service_extend.save_model_credentials_without_validate(
                tenant_id=tenant_id,
                provider=provider_model_record.provider_name,
                model_type=provider_model_record.model_type,
                model=provider_model_record.model_name,
                credentials=origin_credentials,
            )

            model_provider_service_extend.create_tenant_model_sync_if_not_exist(
                tenant_id=tenant_id, model_id=model_id, origin_model_id=provider_model_record.id, is_all=True
            )
            logging.info(
                f"{tenant_id} 同步成功: {provider_model_record.provider_name}，{provider_model_record.model_name}"
            )

        # 同步只有供应商的模型数据
        provider_records = TenantExtendService.get_sync_all_provider()
        for provider_record in provider_records:
            logging.info(f"{tenant_id} 同步的模型所在工作区: {provider_record.tenant_id}")
            logging.info(f"{tenant_id} 同步的模型provider_name: {provider_record.provider_name}")
            origin_provider_credentials = model_provider_service_extend.get_provider_credentials_obfuscated(
                tenant_id=provider_record.tenant_id,
                provider=provider_record.provider_name,
            )

            # 查不到相应的凭证
            if origin_provider_credentials is None:
                logging.error(click.style(f"{tenant_id} 同步失败: {provider_record.tenant_id}，{provider_record.provider_name}", fg="red"))
                continue

            provider_id = model_provider_service_extend.save_provider_credentials_without_validate(
                tenant_id=tenant_id, provider=provider_record.provider_name, credentials=origin_provider_credentials
            )
            model_provider_service_extend.create_tenant_model_sync_if_not_exist(
                tenant_id=tenant_id, model_id=provider_id, origin_model_id=provider_record.id, is_all=True
            )
            logging.info(f"{tenant_id} 同步成功: {provider_record.tenant_id}，{provider_record.provider_name}")

        # 保存执行时间到redis
        redis_client.set("sync_set_all_model_to_tenant_time", tenant.created_at.isoformat(), ex=None)

    logging.info(click.style(f"上次执行最后一个工作区的创建时间:{last_execution_time_str} sync_set_all_model_to_tenant_task_extend 执行完毕", fg="green"))

    return True
