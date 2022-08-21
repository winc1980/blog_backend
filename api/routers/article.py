from fastapi import APIRouter, Depends
import api.schemas.article as article_schema
from sqlalchemy.ext.asyncio import AsyncSession
from typing import List

import api.cruds.article as article_crud
from api.database import get_db

router = APIRouter()


@router.get("/articles", response_model=List[article_schema.Article])
async def list_articles(db: AsyncSession = Depends(get_db)):
    return await article_crud.get_articles(db)

@router.post("/articles", response_model=article_schema.ArticleCreateResponse)
async def create_article(
    article_body: article_schema.ArticleCreate, db: AsyncSession = Depends(get_db)
):
    return await article_crud.create_article(db, article_body)
