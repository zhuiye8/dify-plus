import datetime
import time

import click

import app
from extensions.ext_database import db
from models.api_token_money_extend import ApiTokenMoneyExtend, ApiTokenMoneyMonthlyStatExtend


@app.celery.task(queue="extend_low")
def update_api_token_monthly_used_quota_task_extend():
    click.echo(click.style("Start 重置密钥月额度", fg="green"))
    start_at = time.perf_counter()

    # 快照额度
    api_token_money_monthly_stats = []
    api_token_money_extend = db.session.query(ApiTokenMoneyExtend).all()
    for api_token in api_token_money_extend:
        api_token_money_monthly_stats.append(
            ApiTokenMoneyMonthlyStatExtend(
                app_token_id=api_token.app_token_id,
                accumulated_quota=api_token.accumulated_quota,
                month_used_quota=api_token.month_used_quota,
                month_limit_quota=api_token.month_limit_quota,
                stat_at=datetime.datetime.now(),
                updated_at=datetime.datetime.now(),
                created_at=datetime.datetime.now(),
            )
        )
    db.session.add_all(api_token_money_monthly_stats)

    # 重置密钥月额度
    db.session.query(ApiTokenMoneyExtend).update(
        {
            ApiTokenMoneyExtend.month_used_quota: 0,
            ApiTokenMoneyExtend.updated_at: datetime.datetime.now(),
        }
    )

    db.session.commit()

    end_at = time.perf_counter()
    click.echo(click.style("重置密钥月额度：success latency: {}".format(end_at - start_at), fg="green"))
