## Содержание

#### Основная часть

- [Ansible-роль для установки сервиса CRI-O](#ansible-роль-для-установки-сервиса-cri-o)
- [Ansible-роль для установки утилит kubectl, kubeadm и kubelet](#ansible-роль-для-установки-утилит-kubectl-kubeadm-и-kubelet)

#### Дополнительно

- [Проверка идемпотентности](#проверка-идемпотентности)

### Ansible-роль для установки сервиса CRI-O

В первую очередь после прочтения инструкции к установке было решено вынести необходимые переменные в `defaults`.

```yaml
crio_kubernetes_version: "1.30"
crio_project_path: "prerelease:/main"
```

Далее была заполнена `meta` информация.

```yaml
galaxy_info:
  namespace: solovev
  role_name: crio
  author: Solovev Pavel <solovev_pavel21@mail.ru>
  description: This role is designed to install and configure cri-o.
  company: yadro

  license: MIT

  min_ansible_version: "2.1"
```

И, в конце концов, написана сама роль для установки CRI-O.

```yaml
- name: Ensure necessary packages are installed
  ansible.builtin.apt:
    name: "{{ item }}"
    state: present
    update_cache: true
  loop:
    - ca-certificates
    - curl
    - gpg

- name: Ensure keyrings dir exists
  ansible.builtin.file:
    path: /etc/apt/keyrings
    state: directory
    mode: "0o755"
    owner: root
    group: root

- name: Add official GPG key
  ansible.builtin.apt_key:
    url: "https://pkgs.k8s.io/addons:/cri-o:/{{ crio_project_path }}/deb/Release.key"
    state: present
    keyring: /etc/apt/keyrings/cri-o-apt-keyring.gpg

- name: Add CRI-O apt repository
  ansible.builtin.apt_repository:
    repo: "deb [signed-by=/etc/apt/keyrings/cri-o-apt-keyring.gpg] https://pkgs.k8s.io/addons:/cri-o:/{{ crio_project_path }}/deb/ /"
    state: present
    filename: cri-o

- name: Ensure CRI-O package is installed
  ansible.builtin.apt:
    name: cri-o={{ crio_kubernetes_version }}*
    state: present
    update_cache: true
    cache_valid_time: 0
    clean: true

- name: Ensure CRI-O service is started
  ansible.builtin.service:
    name: crio
    state: started
    enabled: true

- name: Check if swap is enabled
  ansible.builtin.command:
    cmd: swapon -s
  register: swap_check
  changed_when: false

- name: Ensure swap is turned off
  when: swap_check.stdout | length > 0
  block:
    - name: Ensure swap is disabled
      ansible.builtin.command:
        cmd: swapoff -a
      changed_when: true

    - name: Ensure swap is removed from fstab
      ansible.builtin.lineinfile:
        path: '/etc/fstab'
        regexp: '\sswap\s'
        state: 'absent'

- name: Ensure br_netfilter module is loaded
  community.general.modprobe:
    name: br_netfilter
    state: present

- name: Ensure IP forwarding is enabled
  ansible.posix.sysctl:
    name: net.ipv4.ip_forward
    value: 1
    state: present
    reload: true
```

### Ansible-роль для установки утилит kubectl, kubeadm и kubelet

Процесс написания роли для установки kubectl, kubeadm и kubelet оказался несколько менее тривиальным.

В переменные в `defaults` были вынесены следующие значения:

```yaml
k8s_tools_version: "1.30"
k8s_tools_kubelet_args:
  container-runtime-endpoint: unix:///var/run/crio/crio.sock
  runtime-request-timeout: 10m
```

Для удобной работы с `k8s_tools_kubelet_args` флаги и передаваемые с их помощью значения были описаны как объект, а затем, перед записью в `var/lib/kubelet/kubeadm-flags.env`, переведены в строку с использованием `Jinja2`.

```
KUBELET_ARGS="{{ k8s_tools_kubelet_args.items() | map('join', '=') | map('regex_replace', '^', '--') | join(' ') }}"
```

Без явного указания CRI-O рантайма сервис Kubelet упорно отказывается работать.

Так как необходимо указать эндпоинт CRI-O рантайма, необходимо также установить CRI-O сам по себе. Конечно, можно сделать ровно так, как указано в инструкции, и одной ролью установить и CRI-O, и изначально требуемые утилиты, однако Ansible роли были созданы как раз чтобы решать проблему с дублированием кода.

Поэтому, чтобы роль могла в дальнейшем существовать и быть переиспользована, необходимо указать, что она зависит от роли, созданной для установки CRI-O.

Для этого укажем зависимости в `meta/main.yml`:

```yaml
dependencies:
  - name: solovev.crio
```

И в `requirements.yml`:

```yaml
- name: solovev.crio
  scm: git
  src: https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev
  version: e2d7e109c3223fa82b12327cee16b8755a89492b
```

Теперь, перед запуском задач непосредственно по установке kubectl, kubeadm и kubelet, будет также проверяться наличие и, в случае необходимости, проводится установка CRI-O.

<b>Update</b>

Чтобы роль была single-responsibility было принято решение отказаться от внедрения зависимости роли k8s_utils от crio. Вместо CRI-O может быть использован любой другой рантайм, поэтому лучше отказаться от прямой зависимости.

В README к роли должно быть указано, что требуется для ее использования, и должен быть пример с импортом предлагаемой роли.

```yaml
- name: Ensure necessary packages are installed
  ansible.builtin.apt:
    name: "{{ item }}"
    state: present
    update_cache: true
  loop:
    - apt-transport-https
    - ca-certificates
    - curl
    - gpg

- name: Ensure keyrings dir exists
  ansible.builtin.file:
    path: /etc/apt/keyrings
    state: directory
    mode: "0o755"
    owner: root
    group: root

- name: Ensure official GPG key exists
  ansible.builtin.apt_key:
    url: "https://pkgs.k8s.io/core:/stable:/v{{ k8s_tools_version }}/deb/Release.key"
    state: present
    keyring: /etc/apt/keyrings/kubernetes-apt-keyring.gpg

- name: Ensure Kubernetes apt repository exists
  ansible.builtin.apt_repository:
    repo: "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v{{ k8s_tools_version }}/deb/ /"
    state: present
    update_cache: true
    filename: kubernetes

- name: Ensure kubelet, kubeadm, kubectl packages are installed
  ansible.builtin.apt:
    name: "{{ item }}"
    state: present
    update_cache: true
    cache_valid_time: 0
  loop:
    - kubelet
    - kubeadm
    - kubectl

- name: Ensure packages versions are held
  ansible.builtin.dpkg_selections:
    name: "{{ item }}"
    selection: hold
  loop:
    - kubelet
    - kubeadm
    - kubectl

- name: Ensure necessary kubelet arguments are set
  ansible.builtin.template:
    src: var/lib/kubelet/kubeadm-flags.env.j2
    dest: /var/lib/kubelet/kubeadm-flags.env
    mode: "0o644"

- name: Ensure kubelet service is enabled
  ansible.builtin.service:
    name: kubelet
    state: started
    enabled: true
```

### Проверка идемпотентности

Идемпотентность в Ansible подразумевает под собой то, что роль может выполняться неограниченное количество раз, но при этом изменения в конфигурации сервера будут вноситься только в том случае, если он не находится в желаемом состоянии.

Чтобы убедиться в идемпотентности роли можно воспользоваться тестированием (например, `molecule`), либо запустить плейбук несколько раз.

В результатах первого запуска плейбука можно увидеть 11 задач с изменениями.

![image](/docs/summary/hw5_pictures/first_playbook_run.png)

Однако, при последующих запусках вывод в терминал всегда будет одинаков: изменений 0.

![image](/docs/summary/hw5_pictures/second_playbook_run.png)

Просмотреть логи можно в [файле](/ansible/playbook_log.txt).
