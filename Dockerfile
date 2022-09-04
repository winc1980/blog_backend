FROM python:3.10-slim-bullseye

ENV APP_PATH /opt/apps
ENV HOME /root
COPY . ${APP_PATH}
WORKDIR ${APP_PATH}

RUN apt update && apt upgrade -y && apt install -y curl
RUN curl -sSL https://install.python-poetry.org | POETRY_HOME=/opt/poetry python3 - 
ENV PATH="/opt/poetry/bin:$PATH" 

RUN poetry install --no-dev

ARG version

ENTRYPOINT [ "poetry", "run" ]

CMD ["uvicorn", "api.main:app", "--host", "0.0.0.0", "--reload"]