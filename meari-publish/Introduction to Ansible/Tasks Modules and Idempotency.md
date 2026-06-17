---
created: "2026-06-16"
id: tasks-modules-and-idempotency
source: meari-course
study:
  answer: 'A task is a single unit of work in a playbook: it pairs a human-readable name with exactly one module call and its arguments. The module is the code that actually does the work on the managed node, such as apt installing a package or service starting a daemon. Idempotency is the crucial property that running the same task twice leaves the system in the same state: the first run makes the change, the second sees the desired state already holds and does nothing. This is why Ansible reports one of three task states per host: ''ok'' means already in the desired state and no change was needed, ''changed'' means the module had to modify the system to reach that state, and ''failed'' means the task errored. A well-written module checks current state before acting, so a playbook describes the end result rather than a sequence of commands. This differs sharply from a shell script, which blindly re-executes its commands every run and may error or duplicate work on the second pass. Idempotency makes playbooks safe to rerun, which is the foundation of reliable configuration management.'
  kind: essay
  prompt: What does it mean for an Ansible task to be idempotent, how do the ok/changed/failed states reflect this, and why is this safer than running an equivalent shell script repeatedly?
subject: Introduction to Ansible
title: Tasks Modules and Idempotency
---

A **task** is the basic building block of a playbook: one unit of work. Every task has a `name` you write for humans, and a single **module** call that does the actual work. The module receives arguments and carries them out on the managed node. For example, this task uses the `apt` module to ensure a package is present:

```yaml
- name: Install nginx
  apt:
    name: nginx
    state: present
```

The **module** is where the real intelligence lives. Ansible ships hundreds of them, each an expert at one kind of job: `apt` and `yum` manage packages, `service` manages daemons, `copy` and `template` manage files, `user` manages accounts. A good module does not just blindly run a command; it first **checks the current state** of the system, then acts only if a change is needed to reach the state you asked for.

That behavior is the heart of **idempotency**: running the same task twice leaves the system in the same final state. The first time you install nginx, the module installs it. The second time, the module sees nginx is already present and does nothing. You can run a playbook ten times in a row and reach the same result as running it once, with no harm done.

This is why Ansible reports one of three states for each task on each host:

- **ok** — the system was already in the desired state, so nothing changed
- **changed** — the module had to modify the system to reach the desired state
- **failed** — the task hit an error and could not complete

After a run you will see a recap like `ok=4 changed=1 failed=0`. Watching `changed` is how you understand what Ansible actually altered; on a converged system, a healthy rerun shows `changed=0`.

Contrast this with a plain **shell script**. A script like `apt-get install nginx && useradd deploy` simply re-executes every command each time you run it. The second run might error because the `deploy` user already exists, or might duplicate work, or might leave the system in an unexpected state. The script describes *steps to perform*; an Ansible playbook describes the *desired end state* and lets idempotent modules figure out whether any action is required.

This shift, from "do these commands" to "make the system look like this", is what makes configuration management reliable. It lets you safely rerun playbooks to fix drift, apply them to new servers, and trust that they will converge to the same result every time. Keep this in mind as you write each task in [[Your First Playbook]]: aim to describe state, and prefer real modules over the `shell` module so you keep idempotency on your side.
