## Что было сделано и изучено в ходе выполнения заданий

### Основные задачи

- [x] Написан multi-stage [Dockerfile](../../Dockerfile) для приложения
- [x] Написан [docker-compose](../../docker-compose.yml), поднимающий приложение, спулливая его с docker hub, либо выполняя билд из Dockerfile

### Дополнительные задачи

- [x] Прикручен hadolint в пайплайн сначала на свой репозиторий, затем на репозитории всех студентов
```yml
lint:dockerfile:
 image: hadolint/hadolint:fcbd01791c9251d83f2486e61ecaf41ee700a766-debian-amd64
 tags:
    - common
 script:
    - hadolint Dockerfile
 rules:
    - exists:
        - Dockerfile
```
- [x] Написана джоба на билд и пуш docker image в docker hub с импользованием kaniko
```yml
build_and_push_image:
  stage: build
  image:
    name: gcr.io/kaniko-project/executor:v1.22.0-debug
    entrypoint: [""]
  script:
    - echo "{\"auths\":{\"$CI_REGISTRY\":{\"username\":\"$CI_REGISTRY_USER\",\"password\":\"$CI_REGISTRY_PASSWORD\"}}}" > /kaniko/.docker/config.json
    - /kaniko/executor --context "${CI_PROJECT_DIR}" --dockerfile "${CI_PROJECT_DIR}/Dockerfile" --destination "${CI_REGISTRY_IMAGE}:${CI_COMMIT_TAG}"
```

```bash
CI_REGISTRY = https://index.docker.io/v1/
CI_REGISTRY_IMAGE = docker.io/malinkamedok/currency_app
CI_REGISTRY_USER = malinkamedok
CI_REGISTRY_PASSWORD = dckr_hub_token_aaaabbbbb1234
```

- [x] Для облегчения проверки домашних работ студентов билд докер образа с помощью kaniko добавлен и в общий пайплайн

```yml
build image:
  stage: build
  tags:
    - common
  image:
    name: gcr.io/kaniko-project/executor:v1.22.0-debug
    entrypoint: [""]
  script:
    - /kaniko/executor --context "${CI_PROJECT_DIR}" --dockerfile "${CI_PROJECT_DIR}/Dockerfile" --no-push
  rules:
    - exists:
        - Dockerfile
```

- [x] Изучены и применены dev контейнеры
```json
{
  "name": "currency app",
  "dockerFile": "Dockerfile",
  "postCreateCommand": "go mod tidy",
}
```

```docker
FROM docker.io/golang:1.22
WORKDIR /go/src/app
COPY . .
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
RUN go mod init app
RUN go get -d -v ./...
```

Аналогично выполненной работе в рамках 1 домашнего задания, выполненные автоматизации для проверки студенческих работ также сэкономили по меньшей мере 10 часов рабочего времени инженеров. 