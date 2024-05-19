## Содержание

- [Проврека конфигураций ВМ](#проврека-конфигураций-вм)
- [Конфигурирование CRI-O](#конфигурирование-cri-o)
- [Запуск Control Plane](#запуск-control-plane)
- [Проверка работоспособности](#проверка-работоспособности)
- [Установка Ingress Controller](#установка-ingress-controller)

### Проврека конфигураций ВМ

Для начала убедимся в наличии требуемых пакетов и фиксировании их версий.

![image](/docs/summary/hw6_pictures/showhold.png)

Проверим, что выключен netfilter и overlay, настроена сеть.

![image](/docs/summary/hw6_pictures/check_configs.png)

### Конфигурирование CRI-O

Основываясь на документацию, стандартная конфигурация CRI-O находится в `/etc/crio/crio.conf`, однако в моем случае она была обнаружена в `/etc/crio/crio.conf.d/10-crio.conf`. Также ее содержание было несколько отличным от представленного в лекции и документации. Было принято решение скомбинировать их и записать в стандартное место с использованием Ansible роли и template.

```yaml
- name: Ensure CRI-O config is set
  ansible.builtin.template:
    src: etc/crio/crio.conf.j2
    dest: /etc/crio/crio.conf
    mode: "0o644"
```

### Запуск Control Plane

Скопируем конфигурационный файл для kubectl, необходимый для подключения к кластеру, и добавим алиасы. После мы можем убедиться, что нода control-plane поднята успешно.

![image](/docs/summary/hw6_pictures/cp_node.png)

Установим Tigera operator, необходимый для управления жизненным циклом Calico и убедимся в наличии соответствующих namespace, pod, deployment и replicaset.

![image](/docs/summary/hw6_pictures/tigera_operator.png)

Во избежание тяжелых правок конфига после запуска Canico, загрузим конфигурационный файл и исправим конфигурацию calicoNetwork с учетом изначально заданного CIDR.

![image](/docs/summary/hw6_pictures/calico_config.png)

После этого запустим установку Calico с помощью команды `k create -f custom-resources.yaml`, и на этом базовое конфигурирование можно считать оконченным.

### Проверка работоспособности

Проверим доступность worker нод.

![image](/docs/summary/hw6_pictures/nodes.png)

Запустим BusyBox и убедимся в работоспособности пода.

![image](/docs/summary/hw6_pictures/busybox.png)

Мастер нода доступна к просмотру по адресу `178.170.195.95:22002`.

### Установка Ingress Controller

Установим Nginx Ingress Controller с использованием пакетного менеджера Helm.

Для начала необходимо установить сам Helm.

![image](/docs/summary/hw6_pictures/helm.png)

Затем добавляем репозиторий с помощью команды `helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx`, после чего непосредственно устанавливаем сам Nginx ingress controller через команду `helm install ingress-nginx/ingress-nginx --generate-name`.

Убедиться в том, что ingress controller установлен, возможно посмотрев в поды.

![image](/docs/summary/hw6_pictures/pods_with_ingress_controller.png)