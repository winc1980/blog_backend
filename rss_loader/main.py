import os
import sys
import time

import feedparser
from sqlalchemy.ext.asyncio import AsyncSession

sys.path.append(os.pardir)
import asyncio
import datetime as dt
import time

from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import declarative_base, sessionmaker

import api.cruds.articleLink as articleLink_crud
import api.models.articleLink as articleLink_model

ASYNC_DB_URL = "mysql+aiomysql://root:root@db:3306/winc_blog?charset=utf8"

async_engine = create_async_engine(ASYNC_DB_URL, echo=True)
async_session = sessionmaker(
    autocommit=False, autoflush=False, bind=async_engine, class_=AsyncSession
)

Base = declarative_base()

sample_users = [
    {"id": 1, "name": "hibiki", "zenn_id": "hibiki_kato"},
    {"id": 2, "name": "zenn", "zenn_id": "zenn"},
]


async def job():
    db_session = async_session()
    for user in sample_users:
        await zenn(user, db_session)


async def zenn(user, db):
    RSS_URL = "https://zenn.dev/" + user["zenn_id"] + "/feed?all=1"

    d = feedparser.parse(RSS_URL)
    for entry in d.entries:

        published_at = dt.datetime.fromtimestamp(
            time.mktime(entry["published_parsed"])
        ).strftime("%Y-%m-%d")
        articleLink_body = {
            "user": user["id"],
            "link": entry["link"],
            "title": entry["title"],
            "published_at": published_at,
        }
        articleLink = articleLink_model.ArticleLink(**articleLink_body)

        db.add(articleLink)
        await db.commit()
        await db.refresh(articleLink)
        return articleLink


# @asyncio.coroutine
# async def periodic(period):
#     def g_tick():
#         t = time.time()
#         count = 0
#         while True:
#             count += 1
#             yield max(t + count * period - time.time(), 0)
#     g = g_tick()

#     while True:
#         await job()

# loop = asyncio.get_event_loop()
# task = loop.create_task(periodic(1))
# loop.call_later(5, task.cancel)

# try:
#     loop.run_until_complete(task)
# except asyncio.CancelledError:
#     pass


def main():
    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(job())
    finally:
        loop.close()


main()
