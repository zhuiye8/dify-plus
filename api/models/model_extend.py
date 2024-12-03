from extensions.ext_database import db

from .types import StringUUID


class EndUserAccountJoinsExtend(db.Model):
    __tablename__ = "end_user_account_joins_extend"
    __table_args__ = (
        db.PrimaryKeyConstraint("id", name="end_user_account_joins_pkey"),
        db.Index("end_user_account_joins_account_id_idx", "account_id"),
        db.Index("end_user_account_joins_end_user_id_app_id_idx", "end_user_id", "app_id"),
    )

    id = db.Column(StringUUID, server_default=db.text("uuid_generate_v4()"))
    end_user_id = db.Column(StringUUID, nullable=False)
    account_id = db.Column(StringUUID, nullable=False)
    app_id = db.Column(StringUUID, nullable=False)
    created_at = db.Column(db.DateTime, nullable=False, server_default=db.text("CURRENT_TIMESTAMP(0)"))
    updated_at = db.Column(db.DateTime, nullable=False, server_default=db.text("CURRENT_TIMESTAMP(0)"))
