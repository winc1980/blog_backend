from fastapi import FastAPI

from api.routers import user

app = FastAPI()
app.include_router(user.router)