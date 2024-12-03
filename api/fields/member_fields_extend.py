from flask_restful import fields

account_money_fields = {
    "total_quota": fields.Float,
    "used_quota": fields.Float,
}
