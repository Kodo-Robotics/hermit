# Hermit 🧙‍♂️

Hermit is a blazing-fast, minimal, and extensible CLI tool to manage virtual development environments.  
Built in Go as a modern alternative to Vagrant, Hermit focuses on performance, simplicity, and cross-platform support.

---

## ✨ Features

- 🪄 Easy-to-use CLI for creating and managing VMs
- ⚡ Fast startup and lightweight execution
- 🖥️ Native support for VirtualBox (Docker and others coming soon)
- 🧩 Modular backend design
- 📦 YAML-based configuration
- 🔐 Isolated dev environments — like a hermit for your code

---

## 🛠️ Installation

Hermit is under active development. To try it locally:

```bash
git clone https://github.com/<your-username>/hermit.git
cd hermit
go build -o hermit .
./hermit help
```

---

## 🚀 Quick Start

```bash
hermit init               # Generates a default config file
hermit up                 # Spins up a VirtualBox VM
hermit destroy            # Removes the VM
```

---

## 🧱 Configuration

Create a hermit.yaml file like this:

```yaml
name: hermit-dev
memory: 2048
cpus: 2
image: ubuntu-22.04
network:
  type: nat
  port_mappings:
    - host: 8080
      guest: 80
```

---

## 💬 Community

Hermit is built by developers for developers. Contributions, ideas, and discussions are welcome!

## 📄 License

Apache 2.0 License
