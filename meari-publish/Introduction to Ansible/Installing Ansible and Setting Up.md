---
created: "2026-06-16"
id: installing-ansible-and-setting-up
source: meari-course
study:
  answer: 'Installing Ansible only happens on the control node, since the design is agentless. The two common approaches are using your system package manager (for example apt install ansible on Debian or Ubuntu, dnf install ansible on Fedora) or using Pythons package manager with pip install ansible, which often gives a newer version inside a virtual environment. After installing, verify success by running ansible --version, which prints the version and the location of the config file Ansible will use. Next, set up SSH key-based authentication so Ansible can log into managed nodes without prompting for passwords: generate a key with ssh-keygen and copy it to each host with ssh-copy-id. You can optionally create an ansible.cfg file to set defaults such as the inventory path and remote user. Finally, confirm connectivity with the ping module by running ansible all -m ping; a SUCCESS with pong response means Ansible reached the host, authenticated over SSH, and ran Python there. At that point your control node is ready to manage your fleet.'
  kind: essay
  prompt: Describe the steps to get a working Ansible control node, from installation through verifying that it can successfully reach a managed host, and explain what a successful ping module result actually proves.
subject: Introduction to Ansible
title: Installing Ansible and Setting Up
---

Because Ansible is **agentless**, you install it in exactly one place: the **control node** (see [[How Ansible Works (Control Node, Managed Nodes, SSH)]]). Nothing gets installed on the servers you plan to manage. This keeps setup refreshingly simple.

There are two common ways to install. The first is your operating system's **package manager**, which is convenient and integrates with system updates:

```bash
# Debian / Ubuntu
sudo apt update && sudo apt install ansible

# Fedora
sudo dnf install ansible

# macOS (Homebrew)
brew install ansible
```

The second is **pip**, Python's package manager, which often gives you a more recent release and works well inside a Python virtual environment:

```bash
python3 -m pip install --user ansible
```

Once installed, **verify** it worked by checking the version. This is also your first troubleshooting command whenever something seems off:

```bash
ansible --version
```

The output shows the Ansible version and, helpfully, the path to the **config file** it will use. Next comes connectivity. Ansible logs into managed nodes over SSH, so the smoothest experience uses **SSH key-based authentication** rather than typing passwords. Generate a key pair (if you do not already have one) and copy the public key to each managed host:

```bash
ssh-keygen -t ed25519
ssh-copy-id user@your-server
```

You can shape Ansible's defaults with an **`ansible.cfg`** file. Ansible looks for it in the current directory first, which makes a project-local config very handy. A minimal one might be:

```ini
[defaults]
inventory = ./inventory
remote_user = ubuntu
host_key_checking = False
```

Here, `inventory` points at your list of hosts (covered in [[Inventory - Defining Your Hosts]]), `remote_user` sets the login name, and disabling `host_key_checking` avoids interactive prompts in lab setups (leave it on in production).

Finally, prove the whole chain works with the **`ping` module**. This is not a network ICMP ping; it is an Ansible module that connects over SSH and runs Python on the target:

```bash
ansible all -m ping
```

A reply of `"ping": "pong"` with **SUCCESS** confirms three things at once: Ansible reached the host, authenticated over SSH, and successfully executed Python there. If you see that green result, your control node is fully set up and ready to start running real automation.
