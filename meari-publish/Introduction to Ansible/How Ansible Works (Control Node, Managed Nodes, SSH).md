---
created: "2026-06-16"
id: how-ansible-works-control-node-managed-nodes-ssh
source: meari-course
study:
  answer: 'Ansible operates on a simple control-node-to-managed-node model. The control node is the single machine where Ansible is installed and from which you run commands and playbooks; it can be your laptop or a dedicated server. The managed nodes are the machines you are automating. Ansible uses a push model: when you run something, the control node initiates a connection out to the managed nodes, rather than agents on those nodes phoning home. This connection happens over standard SSH, so no special networking or client software is required, only SSH access and a Python interpreter on the target. The workflow for each task is: Ansible generates a small Python program representing the module, copies it to the managed node over SSH, executes it there, captures the result as JSON, and then deletes the temporary file. Because nothing is left behind and no background daemon runs, Ansible is described as agentless. This design keeps managed nodes clean and lowers the operational burden, since there is no agent to install, upgrade, or secure.'
  kind: essay
  prompt: Walk through what happens, step by step, when Ansible runs a single task against a managed node, and explain why this design is described as agentless and push-based.
subject: Introduction to Ansible
title: How Ansible Works (Control Node, Managed Nodes, SSH)
---

To use Ansible effectively, it helps to picture the two kinds of machines involved. The **control node** is the one machine where Ansible itself is installed. This is where you run your commands and playbooks. It might be your own laptop, a build server, or a dedicated automation host. You only ever install Ansible here, not anywhere else. (See [[What Is Ansible and Why Use It]] for why that agentless property matters.)

The other side is the **managed nodes**, sometimes called hosts or targets. These are the servers you want to configure. Ansible does *not* require any Ansible software on them. They simply need to be reachable over **SSH** and to have a **Python** interpreter available, which nearly all Linux distributions ship by default.

Ansible uses a **push model**. When you trigger a run, the control node reaches *out* to the managed nodes and tells them what to do, in real time, over SSH. This is the opposite of "pull" systems, where an agent installed on each server periodically wakes up and pulls instructions from a central server. With Ansible, there is no agent and no scheduled check-in: things happen when, and only when, you run them.

Here is what actually happens for a single task. Ansible takes the **module** you asked for (say, "install a package"), generates a small self-contained Python program for it, and then:

```text
1. Connects to the managed node over SSH
2. Copies the temporary module file to the node
3. Executes it with the node's Python interpreter
4. Reads back the result as structured JSON
5. Deletes the temporary file, leaving nothing behind
```

Because the module is removed after it runs and **no daemon stays resident**, the managed node is left exactly as clean as before, plus whatever change you intended. This is the heart of what "agentless" means in practice.

This architecture has real benefits. There is no agent software to install, patch, or secure on dozens of servers, which removes a whole category of maintenance and security work. Connectivity reuses the SSH setup you already trust. And because runs are initiated by you from the control node, behavior is predictable and easy to reason about. The trade-off is that managed nodes must be reachable over SSH at run time, but for most environments that is already the case. From here, the natural next step is getting Ansible installed on your control node and confirming it can reach your hosts.
