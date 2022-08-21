from sqlalchemy.ext.asyncio import AsyncSession
from typing import List, Tuple

from sqlalchemy import select
from sqlalchemy.engine import Result

import api.models.user as user_model
import api.schemas.user as user_schema


async def create_user(
    db: AsyncSession, user_create: user_schema.UserCreate
) -> user_model.User:
    user = user_model.User(**user_create.dict())
    db.add(user)
    await db.commit()
    await db.refresh(user)
    return user


async def get_users(db: AsyncSession) -> List[Tuple[int, str, bool]]:
    result: Result = await (
        db.execute(
            select(
                user_model.User.id,
                user_model.User.name
            )
        )
    )
    return result.all()
