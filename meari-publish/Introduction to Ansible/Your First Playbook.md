---
created: "2026-06-16"
id: your-first-playbook
source: meari-course
study:
  answer: 'A playbook is a YAML file containing a list of plays, where each play maps a set of hosts to an ordered list of tasks. The key directives are ''hosts'' (which inventory pattern the play targets), ''become'' (whether to escalate privileges, like sudo), and ''tasks'' (the list of work to do). Each task has a human-readable ''name'' and exactly one module call with its arguments. Tasks run top to bottom, in order, on every targeted host. A complete play to install a web server might target the ''web'' group, set ''become: true'', then use the apt module to install nginx and the service module to start and enable it. You run a playbook with ''ansible-playbook site.yml'', and Ansible prints a PLAY and TASK header for each step plus a final PLAY RECAP summarizing ok, changed, and failed counts per host. Because the underlying modules are idempotent, rerunning the same playbook reports ok rather than changed when the system already matches the desired state, which makes playbooks safe to run repeatedly.'
  kind: essay
  prompt: Describe the anatomy of an Ansible playbook (plays, hosts, become, tasks) and walk through what a small playbook that installs and starts a web server does when you run it twice.
subject: Introduction to Ansible
title: Your First Playbook
---

A **playbook** is a YAML file that describes the desired state of your systems. Where an ad-hoc command runs one task by hand, a playbook captures a whole sequence of tasks in a file you can save, review, and rerun. Because it is written in YAML, everything you learned in [[YAML Basics for Ansible]] applies directly here.

At the top level, a playbook is a **list of plays**. Each **play** connects a group of hosts to an ordered list of tasks. The most important directives in a play are:

- **hosts** — the inventory pattern this play runs against (a group, a host, or `all`)
- **become** — whether to escalate privileges, the equivalent of `sudo`
- **tasks** — the list of work to perform, in order

Each **task** has a human-friendly `name` and exactly one **module** call with its arguments. The name is optional but strongly recommended, because it is what Ansible prints as each task runs.

Here is a complete, runnable playbook that installs and starts the nginx web server. Save it as `site.yml`:

```yaml
---
- name: Set up web servers
  hosts: web
  become: true
  tasks:
    - name: Install nginx
      apt:
        name: nginx
        state: present
        update_cache: true

    - name: Start and enable nginx
      service:
        name: nginx
        state: started
        enabled: true
```

This single play targets the `web` group, escalates to root with `become: true`, then runs two tasks: install the nginx package, and ensure the service is running and set to start on boot. You run it from your control node with:

```bash
ansible-playbook site.yml
```

Ansible prints a `PLAY` header, then a `TASK` line for each step, ending with a `PLAY RECAP` that summarizes the results per host:

```
PLAY RECAP ******************************************************
web1 : ok=3  changed=2  unreachable=0  failed=0
```

The magic appears when you run the same playbook a **second time**. Now nginx is already installed and already running, so both tasks report `ok` instead of `changed`, and the recap shows `changed=0`. This is **idempotency**: the playbook describes the desired end state, and Ansible only acts when reality differs from it. That property is what makes playbooks safe to run again and again, and it is explored further in [[Tasks Modules and Idempotency]].
