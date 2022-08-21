from sqlalchemy import create_engine

from api.models.user import Base as user_base

from api.models.article import Base as article_base

DB_URL = "mysql+pymysql://root:root@db:3306/winc_blog?charset=utf8"
engine = create_engine(DB_URL, echo=True)


def reset_database():
    user_base.metadata.drop_all(bind=engine)
    user_base.metadata.create_all(bind=engine)
    article_base.metadata.drop_all(bind=engine)
    article_base.metadata.create_all(bind=engine)


if __name__ == "__main__":
    reset_database()