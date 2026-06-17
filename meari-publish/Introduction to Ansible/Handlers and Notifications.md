---
created: "2026-06-16"
id: handlers-and-notifications
source: meari-course
study:
  answer: 'A handler is a special task that does not run every time the playbook runs. Instead it sits idle until another task notifies it, and it runs only if that notifying task actually reported a change. You define handlers under a play''s handlers: section, and a regular task triggers one with the notify: keyword referencing the handler''s name. The classic use case is restarting a service after its configuration changes. If you deploy a new nginx config with the template module and the file is updated, that task notifies a "restart nginx" handler. If the file is already correct, the template task reports "ok" rather than "changed", so the handler is skipped and the service is not needlessly bounced. This makes plays both efficient and idempotent. Another key behavior is that notified handlers run once, at the very end of the play, after all normal tasks finish. So even if five tasks all notify the same handler, the service restarts a single time. Handlers keep disruptive actions tied to real changes rather than running them on every execution.'
  kind: essay
  prompt: What problem do Ansible handlers solve, how does the notify keyword connect a task to a handler, and why is it significant that handlers run only when notified and only once at the end of a play?
subject: Introduction to Ansible
title: Handlers and Notifications
---

Some actions should only happen when something actually changes. Restarting a web server, reloading a firewall, or rebooting a machine are disruptive, so you do not want them to run on every playbook execution. Ansible handles this with **handlers**: tasks that stay dormant until they are **notified**.

A handler looks exactly like a normal task, but it lives under the play's `handlers:` section instead of `tasks:`. A regular task triggers it using the **`notify:`** keyword, which references the handler by its `name`. The crucial detail is that a handler only fires if the notifying task reports a **changed** state. If the task makes no change, the notification is silently dropped.

Here is the canonical example: deploy a config file with the `template` module, and restart the service only if that file was actually modified.

```yaml
- hosts: web
  tasks:
    - name: Deploy nginx configuration
      ansible.builtin.template:
        src: nginx.conf.j2
        dest: /etc/nginx/nginx.conf
      notify: Restart nginx

  handlers:
    - name: Restart nginx
      ansible.builtin.service:
        name: nginx
        state: restarted
```

On the first run the template task writes the file and reports `changed`, so the **`Restart nginx`** handler runs. On a second run, if the rendered config is identical, the template task reports `ok`, no notification fires, and nginx is left running undisturbed. This is what keeps the play **idempotent** (see [[Tasks Modules and Idempotency]]).

Two behaviors are important to remember. First, notified handlers do not run immediately, they run **once at the very end of the play**, after every normal task has finished. Second, because of this batching, even if many tasks notify the same handler, it still runs only a single time. So if you change three separate nginx settings, each notifying `Restart nginx`, the service is restarted exactly once rather than three times.

Handlers pair naturally with [[Templates with Jinja2]], since config files generated from [[Variables and Facts]] are exactly the kind of thing whose changes should trigger a controlled restart.
