---
created: "2026-06-16"
id: variables-and-facts
source: meari-course
study:
  answer: 'Variables let you write one playbook that adapts to many situations instead of hard-coding values. You can define them in several places: inline under a play''s vars: key, in a separate vars file pulled in with vars_files:, in inventory next to a host or group, or passed on the command line with --extra-vars (which wins over almost everything). You reference a variable with Jinja2 syntax, {{ var }}, for example name: "{{ package }}". When values collide, precedence resolves them; at a high level command-line extras beat play vars, which beat inventory vars, with role defaults at the bottom. Facts are variables Ansible gathers automatically by running the setup module at the start of a play. They live under ansible_facts and describe each managed node: its OS family, IP addresses, CPU count, and memory. Because facts are real data about the target, you can use them in tasks and conditionals, like installing a different package on Debian versus RedHat, making playbooks portable across machines.'
  kind: essay
  prompt: Explain the different places you can define Ansible variables and what facts are, then give an example of how a gathered fact could make a single playbook work across different operating systems.
subject: Introduction to Ansible
title: Variables and Facts
---

Hard-coding values into a playbook makes it rigid. **Variables** solve this by letting you name a value once and reuse it everywhere, so the same playbook can install `nginx` on one run and `apache` on another. Ansible gives you several places to define variables, and you reference all of them the same way: with **Jinja2** double-brace syntax, `{{ var_name }}`.

The most direct place is a `vars:` block inside a play. You can also keep variables in a separate file and load it with `vars_files:`, or attach them to hosts and groups in your inventory (see [[Inventory - Defining Your Hosts]]).

```yaml
- hosts: web
  vars:
    package: nginx
    http_port: 80
  tasks:
    - name: Install the web server
      ansible.builtin.apt:
        name: "{{ package }}"
        state: present
```

When the same variable is set in more than one place, **variable precedence** decides which value wins. You do not need to memorize the full list as a beginner, but the high-level rule is: values passed on the command line with `--extra-vars` override almost everything, then play `vars:`, then inventory variables, and **role defaults** sit at the very bottom as easily-overridden fallbacks.

The other major source of variables is **facts**. At the start of a play, Ansible automatically runs the `setup` module against each managed node and collects a large dictionary of information about it, stored under **`ansible_facts`**. Facts include the operating system family, all IP addresses, total memory, CPU count, and much more. You can inspect them directly:

```bash
ansible web -m setup
```

Because facts describe the real machine, you can use them to make decisions inside tasks. For example, `ansible_facts['os_family']` lets one playbook behave correctly on both Debian and RedHat systems:

```yaml
- name: Show the host's OS family and memory
  ansible.builtin.debug:
    msg: "OS is {{ ansible_facts['os_family'] }}, RAM is {{ ansible_facts['memtotal_mb'] }} MB"
```

Variables and facts are what make playbooks reusable rather than throwaway scripts. Once you are comfortable referencing them, revisit [[Your First Playbook]] to see how they slot into a real play, and explore [[Templates with Jinja2]] to generate config files from these values.
