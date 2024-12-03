import datetime
import time

import click

import app
from extensions.ext_database import db
from models.account_money_extend import AccountMoneyExtend
from models.account_money_monthly_stat_extend import AccountMoneyMonthlyStatExtend


@app.celery.task(queue="extend_low")
def update_account_used_quota_extend():
    click.echo(click.style("Start 更新账号余额额度", fg="green"))
    start_at = time.perf_counter()

    # 快照额度
    account_money_monthly_stats = []
    account_money_extend = db.session.query(AccountMoneyExtend).all()
    for account in account_money_extend:
        account_money_monthly_stats.append(
            AccountMoneyMonthlyStatExtend(
                account_id=account.account_id,
                total_quota=account.total_quota,
                used_quota=account.used_quota,
                stat_at=datetime.datetime.now(),
                updated_at=datetime.datetime.now(),
                created_at=datetime.datetime.now(),
            )
        )
    db.session.add_all(account_money_monthly_stats)

    # 重置用户额度
    db.session.query(AccountMoneyExtend).update(
        {
            AccountMoneyExtend.used_quota: 0,
            AccountMoneyExtend.updated_at: datetime.datetime.now(),
        }
    )

    db.session.commit()

    end_at = time.perf_counter()
    click.echo(click.style("更新账号余额额度：success latency: {}".format(end_at - start_at), fg="green"))
