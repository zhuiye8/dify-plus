"""empty message

Revision ID: 4faed5bbdb91
Revises: ca79d9b5973b, 01d6889832f7
Create Date: 2024-12-03 04:10:33.796159

"""
from alembic import op
import models as models
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '4faed5bbdb91'
down_revision = ('ca79d9b5973b', '01d6889832f7')
branch_labels = None
depends_on = None


def upgrade():
    pass


def downgrade():
    pass
