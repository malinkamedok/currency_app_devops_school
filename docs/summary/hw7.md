## Содержание

- [Установка ArgoCD](#установка-argocd)
- [Конфигурирование Ingress](#конфигурирование-ingress)
- [Манифесты](#манифесты)
- [Тестирование приложения в Jenkins](#тестирование-приложения-в-jenkins)
- [Дополнительные инструменты](#дополнительные-инструменты)

### Установка ArgoCD

Установка ArgoCD \
`kubectl create namespace argocd`
`kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml`

Чтобы получить доступ к ArgoCD UI необходимо сначала прокинуть порт с самого ArgoCD на ВМ, а затем уже с ВМ на локальную машину. \
`kubectl port-forward svc/argocd-server -n argocd 8080:443` \
`ssh root@192.168.100.163 -p 22 -L 8080:127.0.0.1:8080`

Команда для получения пароля администратора \
`kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 --decode && echo`

![image](/docs/summary/hw7_pictures/argocd_installed.png)

### Конфигурирование Ingress

В предыдущую попытку Ingress был установлен и сконфигурирован неверно. В конце концов удалось решить все проблемы, связанные со слепотой при чтении документации.

![image](/docs/summary/hw7_pictures/ingress_installed.png)

Теперь сеть внутри кластера работает корректно.

### Манифесты

Были написаны deployment, service и ingress манифесты для приложения. Также был написан манифест и для ArgoCD.

В результате работы можно увидеть структуру кластера в Web UI ArgoCD. Также на скриншоте видно, что на данный момент запущена уже вторая версия приложения, автоматически задеплоенная в следствие изменения deployment.

![image](/docs/summary/hw7_pictures/argocd_all_running.png)

Результат выполнения команды `kubectl describe svc` и `kubectl describe ingress`.

![image](/docs/summary/hw7_pictures/describe_services_and_ingress.png)

Чтобы корректно подгружались приложения необходимо деплоить их по определенным тегам, а не по stable/latest/master. В таком случае, при обновлении манифеста, ArgoCD будет автоматически деплоить нужную версию.

### Выключение одной из ВМ

На представленном GIF продемонстрирован процесс выключения одной из виртуальных машин с целью убедиться в том, что кластер выполняет возложенные на него задачи. Запросы продолжают выполняться, ответы поступают.

![image](/docs/summary/hw7_pictures/vm_shutdown.gif)

### Тестирование приложения в Jenkins

Для тестирования приложения, написанного на Go, необходим сам язык, установленный на систему. Для этого воспользуемся плагином `Go`. Внутри самого пайплайна укажем привычные переменные окружения, используемые Go приложениями: `withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin", "GOBIN=${root}/go/bin"])`.

Для генерации `report.xml` была использована утилита `go-junit-report`, принимающая в свой STDIN STDOUT из go test.

![image](/docs/summary/hw7_pictures/test_results.png)

### Дополнительные инструменты

Для удобства просмотра информации о кластере была установлена утилита `k9s`.

![image](/docs/summary/hw7_pictures/k9s.png)

Для отображения хэдеров и красивого JSON была установлена утилита `httpie`.

![image](/docs/summary/hw7_pictures/pretty_curl_with_hostnames.png)