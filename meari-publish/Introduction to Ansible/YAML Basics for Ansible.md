---
created: "2026-06-16"
id: yaml-basics-for-ansible
source: meari-course
study:
  answer: 'YAML is the format Ansible uses for playbooks, and a few rules prevent most beginner errors. Indentation defines structure and must use spaces, never tabs; nested items line up by consistent space counts (commonly two). A mapping is written as ''key: value'' with a required space after the colon. A list uses a leading ''- '' (dash plus space) for each item, and a list of dictionaries combines both. Booleans are bare words: ''true''/''false'' (or yes/no), not quoted. Strings usually need no quotes, but quote them when the value contains a colon-space, a leading special character like ''{'', ''#'', ''*'', or ''@'', or when you want to force it to stay a string (for example a version like "1.10" or "yes"). The ''---'' line marks the start of a document and conventionally opens a playbook. Because indentation carries meaning, a stray tab or a missing space after a colon is the most common cause of "could not find expected" parse errors, so configuring your editor to show whitespace pays off immediately.'
  kind: essay
  prompt: Explain the core YAML rules an Ansible beginner must get right (indentation, mappings, lists, quoting, booleans) and describe two mistakes that commonly break a playbook.
subject: Introduction to Ansible
title: YAML Basics for Ansible
---

Ansible playbooks are written in **YAML**, a data format designed to be human-readable. You do not need to learn all of YAML, just the handful of rules that show up in every playbook. Getting these right will save you from the cryptic parse errors that trip up nearly every beginner.

The single most important rule is **indentation defines structure, and you must use spaces, not tabs**. Items that belong to the same level must be indented the same amount; nesting deeper means indenting further. The convention in Ansible is **two spaces per level**. A literal tab character will break the file, so set your editor to insert spaces and to show whitespace.

A **mapping** (also called a dictionary) is a set of `key: value` pairs. The space after the colon is required:

```yaml
name: install nginx
state: present
```

A **list** is a series of items, each starting with a dash and a space (`- `):

```yaml
packages:
  - nginx
  - git
  - curl
```

Very often you combine the two: a **list of dictionaries**, which is exactly how a list of tasks looks. Each `- ` begins a new item, and the keys inside it align together:

```yaml
tasks:
  - name: install nginx
    apt:
      name: nginx
      state: present
  - name: start nginx
    service:
      name: nginx
      state: started
```

**Strings, booleans, and quoting** cause the most subtle bugs. Most strings need no quotes, but you should quote a value when it contains a colon followed by a space (otherwise YAML thinks it is a new key), or when it starts with a special character like `{`, `#`, `*`, `@`, or `%`. Booleans are written as bare words `true`/`false` (or `yes`/`no`) with no quotes. Watch out: an unquoted `yes` or a version number like `1.10` may be read as a boolean or a number rather than text, so quote them when you mean a string:

```yaml
greeting: "Build status: passing"   # quoted: contains a colon-space
enabled: true                        # boolean, no quotes
version: "1.10"                      # quoted to stay a string
```

Finally, the **`---`** line marks the start of a YAML document and conventionally opens every playbook. Together these rules are all you need to read and write the playbooks in [[Your First Playbook]]. The two mistakes to watch for above all: a stray **tab** in the indentation, and a **missing space after a colon**.
