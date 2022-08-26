from sqlalchemy import Column, Integer, String, ForeignKey, Date
from datetime import date
from sqlalchemy.orm import relationship

from api.database import Base


class ArticleLink(Base):
    __tablename__ = "articleLink"

    id = Column(Integer, primary_key=True)
    user = Column(Integer, ForeignKey('user.id'), nullable=False)
    title = Column(String(128), nullable=False)
    link = Column(String(256), nullable=False)
    published_at = Column(Date, default=date.today, nullable=False)
