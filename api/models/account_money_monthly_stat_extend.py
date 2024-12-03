from extensions.ext_database import db

from .types import StringUUID


class AccountMoneyMonthlyStatExtend(db.Model):
    __tablename__ = "account_money_monthly_stat_extend"
    __table_args__ = (
        db.PrimaryKeyConstraint("id", name="account_money_monthly_stat_pkey"),
        db.Index("idx_account_money_monthly_stat_account_id", "account_id"),
    )

    id = db.Column(StringUUID, server_default=db.text("uuid_generate_v4()"))
    account_id = db.Column(StringUUID, nullable=False)
    total_quota = db.Column(db.Numeric(16, 7))
    used_quota = db.Column(db.Numeric(16, 7))
    stat_at = db.Column(db.DateTime, nullable=False)
    created_at = db.Column(db.DateTime, nullable=False, server_default=db.text("CURRENT_TIMESTAMP(0)"))
    updated_at = db.Column(db.DateTime, nullable=False, server_default=db.text("CURRENT_TIMESTAMP(0)"))
