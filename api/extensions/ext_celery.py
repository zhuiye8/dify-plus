from datetime import timedelta

import pytz
from celery import Celery, Task
from celery.schedules import crontab

from configs import dify_config
from dify_app import DifyApp


def init_app(app: DifyApp) -> Celery:
    class FlaskTask(Task):
        def __call__(self, *args: object, **kwargs: object) -> object:
            with app.app_context():
                return self.run(*args, **kwargs)

    broker_transport_options = {}

    if dify_config.CELERY_USE_SENTINEL:
        broker_transport_options = {
            "master_name": dify_config.CELERY_SENTINEL_MASTER_NAME,
            "sentinel_kwargs": {
                "socket_timeout": dify_config.CELERY_SENTINEL_SOCKET_TIMEOUT,
            },
        }

    celery_app = Celery(
        app.name,
        task_cls=FlaskTask,
        broker=dify_config.CELERY_BROKER_URL,
        backend=dify_config.CELERY_BACKEND,
        task_ignore_result=True,
    )

    # Add SSL options to the Celery configuration
    ssl_options = {
        "ssl_cert_reqs": None,
        "ssl_ca_certs": None,
        "ssl_certfile": None,
        "ssl_keyfile": None,
    }

    celery_app.conf.update(
        result_backend=dify_config.CELERY_RESULT_BACKEND,
        broker_transport_options=broker_transport_options,
        broker_connection_retry_on_startup=True,
        worker_log_format=dify_config.LOG_FORMAT,
        worker_task_log_format=dify_config.LOG_FORMAT,
        worker_hijack_root_logger=False,
        timezone=pytz.timezone(dify_config.LOG_TZ),
    )

    if dify_config.BROKER_USE_SSL:
        celery_app.conf.update(
            broker_use_ssl=ssl_options,  # Add the SSL options to the broker configuration
        )

    if dify_config.LOG_FILE:
        celery_app.conf.update(
            worker_logfile=dify_config.LOG_FILE,
        )

    celery_app.set_default()
    app.extensions["celery"] = celery_app

    imports = [
        "schedule.clean_embedding_cache_task",
        "schedule.clean_unused_datasets_task",
        "schedule.create_tidb_serverless_task",
        "schedule.update_tidb_serverless_status_task",
        "schedule.clean_messages",
        "schedule.update_account_used_quota_extend",  # 二开部分 每月重置账号额度
        "schedule.update_api_token_daily_used_quota_task_extend",  # 二开部分 重置密钥日额度
        "schedule.update_api_token_monthly_used_quota_task_extend",  # 二开部分 重置密钥月额度
        "schedule.sync_set_all_model_to_tenant_task_extend",  # 二开部分 每分钟同步模型给新工作区间
    ]
    day = dify_config.CELERY_BEAT_SCHEDULER_TIME
    beat_schedule = {
        "clean_embedding_cache_task": {
            "task": "schedule.clean_embedding_cache_task.clean_embedding_cache_task",
            "schedule": timedelta(days=day),
        },
        "clean_unused_datasets_task": {
            "task": "schedule.clean_unused_datasets_task.clean_unused_datasets_task",
            "schedule": timedelta(days=day),
        },
        "create_tidb_serverless_task": {
            "task": "schedule.create_tidb_serverless_task.create_tidb_serverless_task",
            "schedule": crontab(minute="0", hour="*"),
        },
        "update_tidb_serverless_status_task": {
            "task": "schedule.update_tidb_serverless_status_task.update_tidb_serverless_status_task",
            "schedule": timedelta(minutes=10),
        },
        "clean_messages": {
            "task": "schedule.clean_messages.clean_messages",
            "schedule": timedelta(days=day),
        },
        # ---------------------------- 二开部分 Begin ----------------------------
        # 每月1号00:00，重置账号额度
        "update_account_used_quota": {
            "task": "schedule.update_account_used_quota_extend.update_account_used_quota_extend",
            "schedule": crontab(minute="0", hour="0", day_of_month="1"),
        },
        # 每天，重置密钥日额度
        "update_api_token_daily_used_quota_task_extend": {
            "task": "schedule.update_api_token_daily_used_quota_task_extend.update_api_token_daily_used_quota_task_extend",
            "schedule": timedelta(days=1),
        },
        # 每月1号00:00，重置密钥月额度
        "update_api_token_monthly_used_quota_task_extend": {
            "task": "schedule.update_api_token_monthly_used_quota_task_extend.update_api_token_monthly_used_quota_task_extend",
            # "schedule": crontab(minute="0", hour="0", day_of_month="1"),
            "schedule": crontab(minute="0", hour="22", day_of_month="7"),  # TODO 临时改到7号22点执行
        },
        # 每1分钟执行一次，同步模型到新增工作区间
        "sync_set_all_model_to_tenant_task_extend": {
            "task": "schedule.sync_set_all_model_to_tenant_task_extend.sync_set_all_model_to_tenant_task_extend",
            "schedule": timedelta(minutes=1),
        },
        # ---------------------------- 二开部分 End ----------------------------
    }
    celery_app.conf.update(beat_schedule=beat_schedule, imports=imports)

    return celery_app
