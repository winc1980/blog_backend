from typing import Optional

from pydantic import BaseModel, Field



class UserBase(BaseModel):
    name: Optional[str] = Field(None, example="sample user")


class UserCreate(UserBase):
    pass


class UserCreateResponse(UserCreate):
    id: int

    class Config:
        orm_mode = True


class User(UserBase):
    id: int

    class Config:
        orm_mode = True