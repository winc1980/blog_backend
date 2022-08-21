FROM python:3.10-slim-bullseye

ENV APP_PATH /opt/apps
ENV HOME /root
ENV PATH ${PATH}:${HOME}/.poetry/bin
COPY . ${APP_PATH}
WORKDIR ${APP_PATH}

RUN apt update && apt upgrade -y && apt install -y curl
RUN curl -sSL https://raw.githubusercontent.com/python-poetry/poetry/master/get-poetry.py | python -

RUN poetry install --no-dev

ARG version

ENTRYPOINT [ "poetry", "run" ]

CMD ["uvicorn", "api.main:app", "--host", "0.0.0.0", "--reload"]