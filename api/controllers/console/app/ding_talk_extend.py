from flask import redirect
from flask_restful import Resource, reqparse

from controllers.console.app.error_extend import DingTalkNotExist
from services.ding_talk_extend import DingTalkService

from .. import api


class DingTalk(Resource):
    def get(self):
        """
        DingTalk login
        """
        parser = reqparse.RequestParser()
        parser.add_argument("code", type=str, required=True, location="args")
        args = parser.parse_args()
        if not (0 < len(args.code) < 500):
            raise DingTalkNotExist
        token, err = DingTalkService.get_user_info(args.code)
        if len(err) > 0:
            raise DingTalkNotExist(err)
        return redirect(token)


api.add_resource(DingTalk, "/ding-talk/login")
