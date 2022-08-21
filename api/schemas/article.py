from typing import Optional

from pydantic import BaseModel, Field


class ArticleBase(BaseModel):
    user_id: int = Field(None, description="user")
    title: Optional[str] = Field(None, example="sample title")

class ArticleCreate(ArticleBase):
    content: Optional[str] = Field(None, example="sample content")
    pass


class ArticleCreateResponse(ArticleCreate):
    id: int

    class Config:
        orm_mode = True


class Article(ArticleBase):
    id: int

    class Config:
        orm_mode = True
