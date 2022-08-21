from fastapi import FastAPI

from api.routers import user, article

app = FastAPI()
app.include_router(user.router)
app.include_router(article.router)