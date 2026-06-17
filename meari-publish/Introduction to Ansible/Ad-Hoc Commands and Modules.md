---
created: "2026-06-16"
id: ad-hoc-commands-and-modules
source: meari-course
study:
  answer: 'Ad-hoc commands let you run a single Ansible task against hosts without writing a playbook, using the form ''ansible <pattern> -m <module> -a "<args>"''. The pattern selects hosts from your inventory (a group name, a hostname, or ''all''); ''-m'' names the module that does the work; ''-a'' passes that module its arguments. For example, ''ansible web -m ping'' checks connectivity, ''ansible all -m service -a "name=nginx state=restarted" --become'' restarts a service, and ''ansible db -m copy -a "src=./my.cnf dest=/etc/my.cnf"'' pushes a file. Common modules include ping, command, shell, copy, file, apt/yum, service, and user. Ad-hoc is ideal for quick, throwaway operations: checking uptime, gathering a fact, restarting one service, or copying a file once. When the work is repeatable, needs ordering, variables, or must be version-controlled and shared, a playbook is the right tool instead. Ad-hoc commands are essentially playbooks of one task, run from the command line.'
  kind: essay
  prompt: When would you reach for an ad-hoc Ansible command instead of writing a playbook, and how is an ad-hoc command structured? Give two concrete examples.
subject: Introduction to Ansible
title: Ad-Hoc Commands and Modules
---

An **ad-hoc command** is a way to run a single Ansible task directly from the command line, without writing a playbook file. It is perfect for quick, one-off jobs: checking if hosts are reachable, restarting a service, or copying a file to many machines at once. The general shape is always the same:

```bash
ansible <pattern> -m <module> -a "<arguments>"
```

The **pattern** decides *which* hosts to target, drawn from your [[Inventory - Defining Your Hosts]]. It can be a group name like `web`, a single hostname, or the special word `all` for every host. The `-m` flag names the **module**, the unit of code that actually does the work, and `-a` passes that module its arguments as `key=value` pairs.

A few of the modules beginners use most often:

- **ping** — not an ICMP ping, but a connectivity and Python check: `ansible all -m ping`
- **command** — runs a program safely without a shell (no pipes or variables): `ansible web -m command -a "uptime"`
- **shell** — like command but through a shell, so `|`, `>`, and `$VAR` work: `ansible web -m shell -a "ps aux | grep nginx"`
- **copy** / **file** — push a file, or manage permissions and directories
- **apt** / **yum** — install packages on Debian/Ubuntu or RHEL/CentOS
- **service** — start, stop, restart, or enable a service
- **user** — create or remove user accounts

Many tasks change the system and need root, so you add `--become` (privilege escalation, like `sudo`):

```bash
# Install nginx on all web hosts
ansible web -m apt -a "name=nginx state=present" --become

# Make sure it is running
ansible web -m service -a "name=nginx state=started" --become
```

Notice that `ping` and `command` only *read* the system, but `apt` and `service` can *change* it. Ansible modules are written to be **idempotent** where possible: running the install command twice will report `ok` the second time instead of `changed`, because nginx is already present. This is the same safety property that makes playbooks reliable.

So when should you use ad-hoc versus a playbook? Reach for **ad-hoc** when the action is quick, exploratory, and you do not need to remember or repeat it: "is everything up?", "restart this one service", "what kernel is each box running?". Reach for a **playbook** when the work is repeatable, has multiple ordered steps, uses variables, or needs to live in version control so your team can review and rerun it. An ad-hoc command is really just a playbook with a single task, typed out by hand.
