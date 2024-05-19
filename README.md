[![pipeline status](https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev/badges/main/pipeline.svg)](https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev/-/commits/main)
[![Latest Release](https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev/-/badges/release.svg)](https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev/-/releases)
[![GitHub](https://img.shields.io/badge/GitHub-malinkamedok-blue?logo=github)](https://github.com/malinkamedok)

## Учебный проект Соловьева Павла на курсе DevOps от Yadro

### Introduction

В рамках прохождения курса DevOps ожидается получение реального опыта работы в DevOps процессах, знакомство с соответствующими инструментами, развитие soft skills и углубление понимания SDLС (Software Development Life Cycle).

### Сourse program

- [x] Вводное занятие [What have I done?](https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev/-/merge_requests/1)
- [x] Разработка базового клиент-серверного приложения на Go [What have I done?](docs/summary/hw1.md)
- [x] Упаковка разработанного приложения в Docker [What have I done?](docs/summary/hw2.md)
- [x] Установка Jenkins. Jenkins Freestyle project [What have I done?](docs/summary/hw3.md)
- [x] Написание Jenkins pipeline на Groovy [What have I done?](docs/summary/hw4.md)
- [x] Настройка инфраструктуры под k8s с использованием Ansible [What have I done?](docs/summary/hw5.md)
- [x] Развертывание k8s [What have I done?](docs/summary/hw6.md)
- [ ] CI/CD для запуска приложения в Kubernetes
- [ ] Финальное демо
- [ ] Получение оффера :P

### Application description

#### Приложение для парсинга курсов валют ЦБ РФ

Приложение, отдающее курс валюты по ЦБ РФ за определенную дату. Для получения курсов валют используется официальное API ЦБ РФ.

#### Структура приложения

```bash
.
├── cmd
│   └── main                # Точка входа в приложение
├── internal
│   ├── app                 # Сборка, настройка и запуск всех компонентов
│   ├── config              # Конфигурации
│   ├── controller
│   │   └── http
│   │       └── v1          # Обработка HTTP запросов
│   ├── entity              # Сущности
│   └── usecase             # Бизнес-логика
│         └── cbrf          # Запрос и обработка данных от ЦБ РФ
└── pkg
    ├── httpserver          # Конфигурации HTTP сервера
    └── web                 # Конфигурации для обработки JSON-ответов
```

#### Документация

#### [OpenApi](/docs/openapi.yaml) документация внутри GitLab

#### OpenApi документация, доступная по GET запросу

<summary><code>GET</code> <code><b>/</b></code> <code>docs</code> <code><b>/</b></code></summary>

#### Получение информации о приложении

<details>
 <summary><code>GET</code> <code><b>/</b></code> <code>info</code></summary>

##### Example output

```json
{
  "version": "0.1.0",
  "service": "currency",
  "author": "p.solovev"
}
```

</details>

#### Получение курса валюты за определенную дату

<details>
 <summary><code>GET</code> <code><b>/</b></code> <code>info</code> <code><b>/</b></code> <code>currency</code></summary>

##### Parameters

> | name     | type     | data type | example    | description                 |
> |----------|----------|-----------|------------|-----------------------------|
> | currency | optional | string    | USD        | Валюта в стандарте ISO 4217 |
> | date     | optional | string    | 2016-01-06 | Дата в формате YYYY-MM-DD   |

##### Example output

```json
{
    "data": {
      "USD": "33,4013"
    },
    "service": "currency"
}
```

</details>

#### Запуск приложения

##### Bash

```bash
go mod tidy
go build -o app cmd/main/main.go
./app
```

##### Docker

```bash
docker build -t solovev_currency_app . && docker run solovev_currency_app
```

##### Docker compose
```bash
docker compose up
```

##### Приложение также доступно к загрузке из Docker Hub

```bash
docker pull malinkamedok/currency_app:latest
```
