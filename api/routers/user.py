from fastapi import APIRouter, Depends
import api.schemas.user as user_schema
from sqlalchemy.ext.asyncio import AsyncSession
from typing import List

import api.cruds.user as user_crud
from api.database import get_db

router = APIRouter()


@router.get("/users", response_model=List[user_schema.User])
async def list_users(db: AsyncSession = Depends(get_db)):
    return await user_crud.get_users(db)

@router.post("/users", response_model=user_schema.UserCreateResponse)
async def create_user(
    user_body: user_schema.UserCreate, db: AsyncSession = Depends(get_db)
):
    return await user_crud.create_user(db, user_body)
