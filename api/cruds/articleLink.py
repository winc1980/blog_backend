from sqlalchemy.ext.asyncio import AsyncSession
from typing import List, Tuple

from sqlalchemy import select
from sqlalchemy.engine import Result

import api.models.articleLink as articleLink_model
import api.schemas.articleLink as articleLink_schema


async def create_articleLink(
    db: AsyncSession, articleLink_create: articleLink_schema.ArticleLinkCreate
) -> articleLink_model.ArticleLink:

    articleLink = articleLink_model.ArticleLink(**articleLink_create.dict())
    db.add(articleLink)
    await db.commit()
    await db.refresh(articleLink)
    return articleLink


async def get_articleLinks(db: AsyncSession) -> List[Tuple[int, str, bool]]:
    result: Result = await (
        db.execute(
            select(
                articleLink_model.ArticleLink.id,
                articleLink_model.ArticleLink.title,
                articleLink_model.ArticleLink.user,
                articleLink_model.ArticleLink.link,
                articleLink_model.ArticleLink.published_at
            )
        )
    )
    return result.all()
