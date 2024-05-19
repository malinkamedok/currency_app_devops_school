Установка ArgoCD
`kubectl create namespace argocd`
`kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml`

Прокидывание портов с ВМ на локальную машину
`ssh pavel@178.170.195.95 -p 22002 -L 8081:127.0.0.1:8081`

Команда для получения пароля администратора
`kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 --decode && echo`