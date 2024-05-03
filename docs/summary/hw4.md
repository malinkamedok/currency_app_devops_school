## Содержание

#### Основная часть

- [Docker Hub credentials](#docker-hub-credentials)
- [Multibranch Pipeline](#multibranch-pipeline)

#### Дополнительно

- [Установка и конфигурирование Buildah на агентах с помощью Ansible](#установка-и-конфигурирование-buildah-на-агентах-с-помощью-ansible)
- [Подключение и использование Shared Library](#подключение-и-использование-shared-library)
- [Передача переменных с помощью Map](#передача-переменных-с-помощью-map)
- [Проверка наличия переданных переменных](#проверка-наличия-переданных-переменных)
- [Очистка рабочего пространства](#очистка-рабочего-пространства)

### Docker Hub credentials

В ходе выполнения работы по написанию `Freestyle project` была имплементирована загрузка образа на Docker Hub. Данные для входа на платформу передавались через параметры. Такой подход работает, но имеет существенный недостаток - пароль видно в логах пайплайна.

![image](/docs/summary/hw4_pictures/dockerhub_creds_freestyle.png)

Для решения проблемы воспользуемся предложенным на лекции методом с добавленим `credentials` внутри Jenkins и использованием их в пайплайне с плагином `Credentials Binding`.

![image](/docs/summary/hw4_pictures/jenkins_add_creds.png)

### Multibranch Pipeline

Для выполнения задачи был создан Multibranch Pipeline. Затем был автоматически запущен скрипт сканирования репозитория на наличие `Jenkinsfile`.

![image](/docs/summary/hw4_pictures/multibranch_pipeline.png)

Jenkinsfile был найден, и скрипт успешно выполнен.

![image](/docs/summary/hw4_pictures/job_succeeded.png)

### Установка и конфигурирование Buildah на агентах с помощью Ansible

В продолжение темы изучения Ansible и поднятия агентов Jenkins с его помощью была дополнительно написана роль, устанавливающая и конфигурирующая Buildah. Затем данная роль была запущена на ВМ с агентами, и теперь ими можно полноценно пользоваться.

```yaml
- name: Install buildah
  ansible.builtin.apt:
    update_cache: true
    cache_valid_time: 0
    clean: true
    name: buildah
```

```yaml
---
# tasks file for buildah

- name: Include appropriate tasks for OS family
  ansible.builtin.include_tasks: "{{ os_specific_tasks }}"
  with_first_found:
    - "os/{{ ansible_distribution | lower | replace(' ', '_') }}-{{ ansible_distribution_version }}.yml"
    - "os/{{ ansible_distribution | lower | replace(' ', '_') }}-{{ ansible_distribution_major_version }}.yml"
    - "os/{{ ansible_distribution | lower | replace(' ', '_') }}.yml"
    - "os/{{ ansible_os_family | lower }}.yml"
    - "os/default.yml"
  loop_control:
    loop_var: os_specific_tasks

- name: Configure buildah registries
  ansible.builtin.template:
    dest: /etc/containers/registries.conf
    mode: "0644"
    src: registries.conf.j2
```

![image](/docs/summary/hw4_pictures/ansible_install_buildah.png)

### Подключение и использование Shared Library

Первым делом для использования Shared Library необходимо создать пустой репозиторий, либо пустую ветку, внутри которых уже могут находиться функции, переиспользуемые в разных пайплайнах. Для этого была создана ветка [p.solovev/shared_library](https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev/-/tree/p.solovev/shared_library?ref_type=heads).

Первая функция представляет собой сборку docker образа с передачей в нее параметров через Config map.

```groovy
def call(Map config = [:]) {
    sh("buildah bud --tag ${config.imageName}:${config.tag} -f ${config.pathToDockerfile} .")
}
```

В Jenkinsfile вызов функции из Shared Library выглядит следующим образом.

```groovy
@Library("shared_library") _
stage('Build docker image') {
    dir("repo") {
        buildDockerImage(imageName:"malinkamedok/currency_app", tag:"latest", pathToDockerfile:"./Dockerfile")
    }
}
```

Чтобы включить использование Shared Library внутри пайплайна необходимо в конфигурации пайплана во вкладке `Properties/Pipeline Libraries` указать название библиотеки и ее источник.

![image](/docs/summary/hw4_pictures/pipeline_properties.png)

В логах пайплайна можно убедиться в корректности выполнения данной функции.

![image](/docs/summary/hw4_pictures/pipeline_log_with_shared_library.png)

По такому же принципу были написаны функции для чекаута и загрузки docker образа на любой container registry.

```groovy
def call(Map config = [:]) {
    sh """
    buildah login -u ${config.hub_username} -p ${config.hub_password} ${config.registry_url}
    buildah push ${config.imageName}:${config.tag}
    buildah logout ${config.registry_url}
    """
}
```

Убедиться в том, что функции из Shared Library действительно используются в пайплайне, можно в окне `Replay`.

![image](/docs/summary/hw4_pictures/jenkins_replay.png)

### Передача переменных с помощью Map

Для повышения читаемости кода и переиспользования уже объявленных переменных было принято решение вынести их объявление в отдельную map. Таким образом не придется искать все места, где эти переменные требуются, и достаточно будет лишь раз их определить.

```groovy
def configMap = [
        imageName:"malinkamedok/currency_app",
        tag:"latest",
        pathToDockerfile:"./Dockerfile",
        registryUrl: "docker.io",
        branchName: BRANCH_NAME,
        credentialsId: 'gitlab',
        repoUrl: 'https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev'
        ]
```

Внутри библиотечных функций были объявлены стандартные значения для переменных, которые зачастую могут быть одинаковы от проекта к проекту. Таким образом, если внутрь функции не передать значения этих переменных, будут использваны стандартные.

```groovy
def call(Map config = [:]) {
    withCredentials([usernamePassword(credentialsId: 'DockerHub', passwordVariable: 'HUB_PASSWORD', usernameVariable: 'HUB_USERNAME')]) {
        def defaultConfig = [
            'registry_url': 'docker.io',
            'tag': 'latest'
        ] << config

        sh """
        buildah login -u ${HUB_USERNAME} -p ${HUB_PASSWORD} ${config.registry_url}
        buildah push ${config.imageName}:${config.tag}
        buildah logout ${config.registry_url}
        """
    }
}
```

После проведенного рефакторинга код основного пайплайна стал намного читаемее.

```groovy
def checkoutConfigMap = [
    branchName: BRANCH_NAME,
    credentialsId: 'gitlab',
    repoUrl: 'https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev'
]

def dockerConfigMap = [
    imageName:"malinkamedok/currency_app",
    tag:"latest",
    pathToDockerfile:"./Dockerfile",
    registryUrl: "docker.io",
]

@Library("shared_library") _
node {
    stage("Checkout repo") {
        checkoutRepo(checkoutConfigMap)
    }
    stage('Build docker image') {
        buildDockerImage(dockerConfigMap)
    }
    stage('Push docker image to docker hub') {
        pushDockerImage(dockerConfigMap)
    }
}
```

### Проверка наличия переданных переменных

Также в ходе рефакторинга было решено ввести проверки на наличие всех переменных, значения которых не могут быть взяты из стандартных. Таким образом, в случае отсутствия переменной, в консоли пайплайна будет видно, какой из них не хватает.

```groovy
if (!config.containsKey('credentialsId')) {
    error("Не передан обязательный параметр 'credentialsId'")
}

if (!config.containsKey('repoUrl')) {
    error("Не передан обязательный параметр 'repoUrl'")
}
```

![image](/docs/summary/hw4_pictures/pipeline_fails_on_empty_var.png)

### Очистка рабочего пространства

Рабочие пространства могут занимать значительное количество дискового пространства, особенно если в проекте накапливается большое количество веток. Чтобы поддерживать систему в чистоте, убирать все лишние артефакты и не терять в прозводительности со временем был использован плагин `Workspace Cleanup`.

Для его использования в конце пайплайна в любом случае, независимо от результата выполнения предыдущих шагов, была использована конструкция try-finally.

```groovy
@Library("shared_library") _
node ('swarm') {
    try {
        stage("Checkout repo") {
        checkoutRepo(checkoutConfigMap)
        }
        stage('Build docker image') {
            buildDockerImage(dockerConfigMap)
        }
        stage('Push docker image to docker hub') {
            pushDockerImage(dockerConfigMap)
        }
    } finally {
        cleanWs(notFailBuild: true)
    }
}
```