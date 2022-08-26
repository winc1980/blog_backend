from datetime import date
from doctest import Example
from typing import Optional

from pydantic import BaseModel, Field


class ArticleLinkBase(BaseModel):
    user: int = Field(None, description="user")
    link: Optional[str] = Field(None, example="sample url")
    title: Optional[str] = Field(None, example="sample title")
    published_at: Optional[date] = Field(None)

class ArticleLinkCreate(ArticleLinkBase):
    pass


class ArticleLinkCreateResponse(ArticleLinkCreate):
    id: int

    class Config:
        orm_mode = True


class ArticleLink(ArticleLinkBase):
    id: int

    class Config:
        orm_mode = True
