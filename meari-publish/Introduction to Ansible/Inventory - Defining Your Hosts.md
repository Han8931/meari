---
created: "2026-06-16"
id: inventory-defining-your-hosts
source: meari-course
study:
  answer: 'An inventory is the file (or files) that tells Ansible which machines it can manage and how they are organized. Without an inventory, Ansible has no targets. Hosts can be listed individually and arranged into named groups, so you can target, for example, all webservers or all databases with a single command instead of naming each machine. Inventories come in two formats: a simple INI style, where group names appear in square brackets followed by their hosts, and a more structured YAML style that nests hosts and variables under groups. You can attach host variables (like a custom SSH port or user) and group variables to set values that apply to specific machines or whole groups. Ansible ships with a default inventory at /etc/ansible/hosts, but the strongly recommended practice is to keep a project-local inventory file inside your project and point to it, either with the -i flag or via ansible.cfg. This keeps each projects hosts self-contained, version-controlled, and portable across machines.'
  kind: essay
  prompt: What is an Ansible inventory and why are groups and host variables useful? Compare the default /etc/ansible/hosts inventory with a project-local inventory and explain when you would prefer each.
subject: Introduction to Ansible
title: Inventory - Defining Your Hosts
---

Ansible needs to know *which* machines it is allowed to manage. That list lives in an **inventory**. Without one, commands like `ansible all -m ping` have nothing to talk to. The inventory is simply a file that names your hosts and, optionally, organizes them into groups and attaches variables to them.

The simplest inventory uses **INI format**. Group names go in square brackets, and the hosts in each group are listed underneath:

```ini
[webservers]
web1.example.com
web2.example.com

[databases]
db1.example.com
```

**Groups** are the real power here. Instead of naming machines one at a time, you can target a whole group at once. For example, `ansible webservers -m ping` talks only to the web servers. There is also an implicit `all` group that includes every host, which is what `ansible all` refers to. You can even nest groups for larger setups, but groups of related hosts are enough to start.

Ansible also supports a **YAML format** inventory, which many people prefer for larger or more structured setups because it nests cleanly and handles variables well:

```yaml
all:
  children:
    webservers:
      hosts:
        web1.example.com:
        web2.example.com:
    databases:
      hosts:
        db1.example.com:
```

You can attach **host variables** to individual machines and **group variables** to whole groups. These let you record per-host details, such as a non-standard SSH port or a specific login user, right next to the host:

```ini
[webservers]
web1.example.com ansible_user=deploy ansible_port=2222
```

Here `ansible_user` and `ansible_port` are special connection variables Ansible understands. You will learn more about using your own variables in [[Variables and Facts]].

Finally, where does the inventory live? Ansible has a **default** location at `/etc/ansible/hosts`, a system-wide file. It works, but it is global and not tied to any particular project. The recommended practice is a **project-local inventory**: a file you keep inside your project directory (often just named `inventory` or `hosts`) and point Ansible at explicitly. You can do this on the command line with `-i`:

```bash
ansible all -i ./inventory -m ping
```

Or set it once in your `ansible.cfg`, as shown in [[Installing Ansible and Setting Up]]. A project-local inventory keeps each project's hosts self-contained, easy to commit to version control, and portable between teammates, which is why it is preferred for real work. With your hosts defined, you are ready to start running commands against them.
