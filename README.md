# 5l4pp3r

> 📸 Your system's configuration, frozen in time.



## 🚀 High-Level Purpose

5l4pp3r is a forensic snapshot tool designed to capture a comprehensive view of your system's configuration environment. It's like a high-resolution camera for your system's state, providing IT professionals and forensic analysts with a powerful lens to examine system configurations at any given point in time.

## 🔍 What It Does

5l4pp3r meticulously collects and stores:

- 🖥️ **System Information**: Hostname and timestamp
- 🌐 **Network Details**: IP addresses, MAC addresses, interface names
- 📁 **Configuration Files**: From standard system directories and user-specific locations


All this data is compressed and stored in a structured database (SQLite or PostgreSQL), creating a space-optimized, point-in-time record of your system's state.

## 🏗️ Architectural Overview

### Key Components:

1. **Configuration Loading** (`internal/config`)

1. Reads `config.toml` for flexible customization
2. Defines database settings, compression algorithms, scan directories, and more



2. **Logging and Instrumentation**

1. Utilizes `zerolog` for structured, timestamped logs



3. **Storage Setup** (`internal/storage`)

1. Supports SQLite (local) and PostgreSQL (centralized)
2. Ensures proper schema creation and verification



4. **Data Gathering** (`internal/gatherer`)

1. Collects system info, network details, and configuration files
2. Compresses file contents for space efficiency





## 💾 Data Ingestion and Persistence Flow

1. Insert System Info (creates `system_id`)
2. Assign `system_id` to Config Files
3. Insert Network Interfaces (linked to `system_id`)
4. Insert Config Files (compressed, with metadata)
5. Commit the Transaction


## 🕵️ Forensic and IT Professional Value

- **Immutable Point-in-Time State**: Reconstruct system settings at snapshot time
- **Relational Data Model**: Powerful querying capabilities
- **Repeatable and Extensible**: Track configuration evolution over time
- **Centralization and Aggregation**: Create a global forensic data lake (with PostgreSQL)


## 🚀 Getting Started

1. Clone the repository:

```plaintext
git clone https://github.com/copyleftdev/5l4pp3r.git
```


2. Configure `config.toml` with your desired settings
3. Build and run:

```plaintext
go build
./5l4pp3r
```




## 📊 Example Output

```plaintext
11:25AM INF Starting 5l4pp3r...
11:26AM INF Snapshot completed successfully.
```

## 🛠️ Possible Enhancements

- Filtering and Exclusions
- Extended Metadata and Integrity Checks
- Integration with CI/CD and Automation Tools


## 🤝 Contributing

We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- All the amazing open-source libraries that made this project possible
- The forensic IT community for inspiration and use cases


---

Remember: With great power comes great responsibility. Use 5l4pp3r ethically and legally! 🦸‍♂️🦸‍♀️
