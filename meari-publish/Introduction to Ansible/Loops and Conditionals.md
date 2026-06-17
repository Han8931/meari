---
created: "2026-06-16"
id: loops-and-conditionals
source: meari-course
study:
  answer: 'Loops let a single task repeat over a list instead of writing the task many times. You add a loop: key with a list of items, and inside the task you reference the current element with the special variable item. A common use is installing several packages or creating several users in one block, which keeps playbooks short and easy to maintain. Conditionals let a task run only when a test is true, controlled by the when: keyword. The expression after when: is plain Jinja-style and usually checks a fact or variable, for example when: ansible_facts[''os_family''] == "Debian", so the task is skipped on non-Debian hosts. You can combine the two: when a task has both loop: and when:, the condition is evaluated separately for each item, so some items run and others are skipped. Together, loops reduce repetition and conditionals add decision-making, letting one portable playbook adapt its behavior to each host''s operating system, environment, or variables rather than needing a different playbook per case.'
  kind: essay
  prompt: Describe how loop and when work in Ansible tasks, including what the item variable refers to, and explain what happens when a single task uses both loop and when together.
subject: Introduction to Ansible
title: Loops and Conditionals
---

Writing the same task five times with only the name changed is tedious and error-prone. **Loops** fix this by letting one task iterate over a list. You add a **`loop:`** key containing the list, and inside the task you refer to the current element with the special variable **`item`**.

A typical example is installing several packages at once:

```yaml
- name: Install required packages
  ansible.builtin.apt:
    name: "{{ item }}"
    state: present
  loop:
    - nginx
    - git
    - curl
```

Loops also work well with structured data. To create several users, each item can be a small dictionary, and you access its fields with `item.field`:

```yaml
- name: Create application users
  ansible.builtin.user:
    name: "{{ item.name }}"
    groups: "{{ item.group }}"
  loop:
    - { name: alice, group: admins }
    - { name: bob, group: developers }
```

**Conditionals** let a task decide whether to run at all. The **`when:`** keyword takes an expression, and the task executes only if it evaluates to true. The expression is plain Jinja-style (note: you do **not** wrap it in `{{ }}` here), and it usually tests a fact or variable. This is how you make a playbook behave differently across operating systems using gathered [[Variables and Facts]]:

```yaml
- name: Install Apache on Debian-family hosts only
  ansible.builtin.apt:
    name: apache2
    state: present
  when: ansible_facts['os_family'] == "Debian"
```

You can build richer tests with `and`, `or`, and comparisons, for example `when: http_port == 80 and environment == "prod"`.

The two features **combine** powerfully. When a single task has both `loop:` and `when:`, Ansible evaluates the condition **separately for each item**. That means some items in the list can run while others are skipped:

```yaml
- name: Install packages, skipping any named "skip-me"
  ansible.builtin.apt:
    name: "{{ item }}"
    state: present
  loop:
    - nginx
    - skip-me
    - git
  when: item != "skip-me"
```

Together, loops cut down repetition and conditionals add decision-making, so one portable playbook can adapt to each host instead of needing a separate version for every case. Next, see how the same variables drive file generation in [[Templates with Jinja2]].
