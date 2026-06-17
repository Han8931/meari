---
created: "2026-06-16"
id: templates-with-jinja2
source: meari-course
study:
  answer: 'A template is a configuration file with placeholders that Ansible fills in per host. You write a .j2 file using the Jinja2 templating language and deploy it with the template module, which renders the file on the control node using that host''s variables and facts, then copies the result to the destination. Jinja2 has two main constructs: {{ expression }} inserts a value, like {{ http_port }} or {{ ansible_facts[''default_ipv4''][''address''] }}, and {% logic %} runs control structures such as for loops and if statements that do not themselves output text. Filters transform values inside expressions using a pipe, for example {{ name | upper }} or {{ port | default(80) }}. Because rendering uses real per-host data, one template can produce a correct, customized config for every machine: each server gets its own IP, hostname, or port substituted automatically. A common pattern is to deploy the rendered config and notify a handler so the service restarts only when the file actually changed, keeping the play idempotent and avoiding needless restarts.'
  kind: essay
  prompt: Explain how the Ansible template module and Jinja2 work together to generate per-host configuration files, covering the difference between expression and logic syntax, the role of filters, and why facts make a single template reusable across many machines.
subject: Introduction to Ansible
title: Templates with Jinja2
---

Often you need a config file that is *almost* the same on every server but differs in a few details: this host's IP, that environment's port, the number of worker processes. Instead of maintaining one file per machine, you write a single **template** with placeholders and let Ansible fill them in. Templates use the **Jinja2** templating language and are deployed with the **`template`** module.

A template file conventionally ends in **`.j2`**. When the `template` module runs, it renders the file **on the control node** using the target host's [[Variables and Facts]], then copies the finished result to the destination path on the managed node.

```yaml
- name: Deploy the application config
  ansible.builtin.template:
    src: app.conf.j2
    dest: /etc/app/app.conf
  notify: Restart app
```

Jinja2 has two core syntaxes. **`{{ }}`** is an **expression** that inserts a value into the output. **`{% %}`** is **logic**: control structures like `for` loops and `if` statements that direct rendering but do not themselves print anything. There are also **filters**, applied with a pipe `|`, that transform a value, such as `{{ name | upper }}` or `{{ port | default(80) }}` (which supplies a fallback when the variable is unset).

Here is a small `app.conf.j2` that pulls from both variables and gathered facts:

```jinja2
# Managed by Ansible - do not edit by hand
server_name = {{ inventory_hostname }}
bind_address = {{ ansible_facts['default_ipv4']['address'] }}
port = {{ http_port | default(8080) }}

{% if enable_logging %}
log_file = /var/log/app.log
{% endif %}

allowed_users:
{% for user in app_users %}
  - {{ user }}
{% endfor %}
```

When this renders against a host where `ansible_facts['default_ipv4']['address']` is `10.0.0.5` and `app_users` is `[alice, bob]`, the output is a concrete file with those exact values substituted. Run the same play against a different host and it gets *its* IP and *its* user list. That is the payoff: **one template, correct output for every machine**, because rendering uses each host's real data.

A best-practice pattern ties this together: deploy the rendered config and `notify:` a handler so the service restarts only when the file actually changed (see [[Handlers and Notifications]]). This keeps the play **idempotent**, no restart happens when the config is already correct. As you grow, templates are usually bundled inside roles for reuse across projects.
