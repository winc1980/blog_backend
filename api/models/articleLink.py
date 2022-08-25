from sqlalchemy import Column, Integer, String, ForeignKey, DateTime
from datetime import datetime
from sqlalchemy.orm import relationship

from api.database import Base


class ArticleLink(Base):
    __tablename__ = "articleLink"

    id = Column(Integer, primary_key=True)
    user = Column(Integer, ForeignKey('user.id'))
    link = Column(String(256))
    published_at = Column(DateTime, default=datetime.now(), nullable=False)
