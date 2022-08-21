from sqlalchemy.ext.asyncio import AsyncSession
from typing import List, Tuple

from sqlalchemy import select
from sqlalchemy.engine import Result

import api.models.article as article_model
import api.schemas.article as article_schema


async def create_article(
    db: AsyncSession, article_create: article_schema.ArticleCreate
) -> article_model.Article:
    article = article_model.Article(**article_create.dict())
    db.add(article)
    await db.commit()
    await db.refresh(article)
    return article


async def get_articles(db: AsyncSession) -> List[Tuple[int, str, bool]]:
    result: Result = await (
        db.execute(
            select(
                article_model.Article.id,
                article_model.Article.title,
                article_model.Article.user_id
            )
        )
    )
    return result.all()
