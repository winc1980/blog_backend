from fastapi import APIRouter, Depends
import api.schemas.articleLink as articleLink_schema
from sqlalchemy.ext.asyncio import AsyncSession
from typing import List

import api.cruds.articleLink as articleLink_crud
from api.database import get_db

router = APIRouter()


@router.get("/articleLinks", response_model=List[articleLink_schema.ArticleLink])
async def list_articleLinks(db: AsyncSession = Depends(get_db)):
    return await articleLink_crud.get_articleLinks(db)

@router.post("/articleLinks", response_model=articleLink_schema.ArticleLinkCreateResponse)
async def create_articleLink(
    articleLink_body: articleLink_schema.ArticleLinkCreate, db: AsyncSession = Depends(get_db)
):
    return await articleLink_crud.create_articleLink(db, articleLink_body)
