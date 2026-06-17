---
created: "2026-06-16"
id: ansible-vault-and-secrets
source: meari-course
study:
  answer: 'Storing passwords, API keys, or TLS certificates as plaintext in playbooks or variable files is dangerous because those files are usually committed to version control, where anyone with repo access, or anyone who breaches it, can read the secrets. Ansible Vault solves this by encrypting sensitive content with a password using strong symmetric (AES256) encryption, so the file is safe to commit while still being usable by Ansible. You manage encrypted files with ansible-vault create to make a new one, edit to change it in place, encrypt to protect an existing file, and decrypt to reverse it. When you run a playbook that references vaulted data you must supply the password, either interactively with --ask-vault-pass or non-interactively with --vault-password-file pointing to a protected file, which is handy for CI. You do not have to encrypt whole files: ansible-vault encrypt_string produces a single encrypted variable you can paste inline, so most of a vars file stays readable while only the secret value is protected. The golden rule is that the vault password itself must never be committed alongside the encrypted data.'
  kind: essay
  prompt: Why is it dangerous to keep secrets in plaintext in your playbooks, and how does Ansible Vault let you store and use sensitive data safely? Mention the key commands and how you supply the vault password at runtime.
subject: Introduction to Ansible
title: Ansible Vault and Secrets
---

Real automation almost always needs **secrets**: database passwords, API tokens, SSH private keys, and TLS certificates. The tempting shortcut is to drop them straight into a variable file. This is a serious mistake. Playbooks and their variables are usually stored in **version control**, so a plaintext password becomes permanently visible to everyone with repository access, and to anyone who ever breaches that repository. Once a secret is committed, it effectively lives in the git history forever.

**Ansible Vault** is the built-in answer. It encrypts files (or individual values) with a password using strong **AES256** encryption. The encrypted file is safe to commit, yet Ansible can transparently decrypt it at runtime when you supply the password. The encryption is symmetric, so the same password locks and unlocks the data.

You work with vaulted files through the `ansible-vault` command. The four core operations are:

```bash
ansible-vault create secrets.yml     # make a new encrypted file
ansible-vault edit secrets.yml       # decrypt in memory, open editor, re-encrypt
ansible-vault encrypt vars/prod.yml  # encrypt an existing plaintext file
ansible-vault decrypt vars/prod.yml  # turn it back into plaintext (use sparingly!)
```

`create` and `edit` are the everyday commands: they prompt for the password and open the content in your editor, never writing plaintext to disk. Use `decrypt` only when you genuinely need the file readable again, since it leaves the secret exposed.

When you run a playbook that references vaulted data, you must provide the password so Ansible can decrypt it. There are two common ways:

```bash
# Prompt for the password interactively
ansible-playbook site.yml --ask-vault-pass

# Read the password from a protected file (great for CI/automation)
ansible-playbook site.yml --vault-password-file ~/.vault_pass.txt
```

The password file is convenient for pipelines, but it must itself be kept off version control and locked down with strict permissions (for example `chmod 600`). **Never commit the vault password next to the encrypted data**, that would defeat the entire purpose.

Finally, you rarely need to encrypt a whole file. Often only one value is sensitive. Ansible Vault can encrypt a **single variable** so the rest of your vars file stays readable in diffs and reviews:

```bash
ansible-vault encrypt_string 'S3cr3tP@ss' --name 'db_password'
```

This prints an encrypted block you paste directly into a normal YAML vars file. Ansible decrypts just that value at runtime. The security-minded habit is simple: keep secrets encrypted, keep the vault password separate and protected, and rotate any secret the moment you suspect it leaked.
