K8s Tools
=========

This role is used to install kubelet, kubeadm and kubectl.

Role Variables
--------------

- k8s_tools_version: "1.30"
- k8s_tools_kubelet_args:
    - container-runtime-endpoint: unix:///var/run/crio/crio.sock
    - runtime-request-timeout: 10m

`k8s_tools_version` variable is used to determine the correct version of packages to be installed.

`k8s_tools_kubelet_args` variable contains the necessary parameters that need to be specified for the correct operation of the Kubelet service when using the CRI-O runtime.

Dependencies
------------

To ensure that the role works correctly, you will need to install CRI-O, containerd, or another container runtime. It is recommended to use the imported role to install CRI-O.

```yaml
- name: solovev.crio
  scm: git
  src: https://gitlab-pub.yadro.com/devops-school-2024/student/p.solovev
  version: e2d7e109c3223fa82b12327cee16b8755a89492b
```

To import the proposed role, please use the following command.

```bash
ansible-galaxy role install -r requirements.yml
```

Example Playbook
----------------

```yaml
- name: Ensure kubelet, kubeadm and kubectl are installed and configured
  hosts: workers
  gather_facts: true
  become: true
  tasks:
    - name: Import role
      ansible.builtin.import_role:
        name: crio
```
