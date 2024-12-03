from functools import wraps

from flask_login import current_user

from controllers.console.error_extend import AccountNoMoneyErrorExtend
from extensions.ext_database import db
from models.account_money_extend import AccountMoneyExtend

#: A proxy for the current user. If no user is logged in, this will be an
#: anonymous user


def money_limit(view):
    """ """

    @wraps(view)
    def decorated(*args, **kwargs):
        account = current_user

        # TODO 需要写入缓存，读缓存
        account_money = db.session.query(AccountMoneyExtend).filter(AccountMoneyExtend.account_id == account.id).first()
        if not account_money:
            return view(*args, **kwargs)

        if account_money.used_quota >= account_money.total_quota:
            raise AccountNoMoneyErrorExtend()

        return view(*args, **kwargs)

    return decorated
