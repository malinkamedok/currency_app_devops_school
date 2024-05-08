CRI-O
=========

This role is used to install CRI-O.

Requirements
------------

For `modprobe` module usage it is necessary to install `community.general` collection.

For `sysctl` module usage it is necessary to install `ansible.posix` collection.

```
ansible-galaxy install -r requierements.yml
```

Role Variables
--------------

- crio_kubernetes_version: "1.30"
- crio_project_path: "prerelease:/main"

`crio_kubernetes_version` variable describes the desired Kubernetes version.

`crio_project_path` variable describes the CRI-O stream that is used.

Example Playbook
----------------

```yaml
- name: Ensure CRI-O is installed
  hosts: workers
  gather_facts: true
  become: true
  tasks:
    - name: Import role
      ansible.builtin.import_role:
        name: crio
```
