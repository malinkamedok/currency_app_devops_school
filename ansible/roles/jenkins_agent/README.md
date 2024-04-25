Setup Jenkins agent
=========

Роль для быстрого поднятия агента для Jenkins.

Описание
------------

#### Тестирование

Для прохождения тестирования необходимо поднять локальный Jenkins. Сделать это можно с помощью [docker compose](../../test_requirements/docker-compose.yml).

После установки Jenkins необходимо создать пользователя для агента c username и password `agent`. 

Также требуется установить [`Swarm`](https://plugins.jenkins.io/swarm/) плагин для Jenkins.

Запуск осуществляется с помощью команды `molecule test`.

#### Запуск

Ansible роль запускается с помощью команды 
```bash
ansible-playbook setup-agent.yml -e jenkins_agent_password=PUT_YOUR_TOKEN_HERE
```
