# Best Practices for Setting Up a Rootstock Node

This guide provides recommendations and references for setting up and configuring an Rootstock node. For complete instructions and the latest configuration details, refer to the [Rootstock DevPortal](https://dev.rootstock.io/).

## Node Setup Recommendations

When setting up an Rootstock node, it’s important to follow the official setup guidelines to ensure stability, performance, and security.

**Reference:**
- [Rootstock Node Setup Guide](https://dev.rootstock.io/node-operators/setup/)

### Key Recommendations
- Follow the installation instructions for your operating system (Linux, Windows, or Docker).
- Keep your node software updated to the latest release to maintain compatibility with the network.
- Configure your node to run as a service for better reliability and automatic restarts.
- Ensure proper disk space and memory allocation according to the [system requirements](https://dev.rootstock.io/node-operators/setup/requirements/).

#### System configurations for Ubuntu
It’s highly recommended to follow the [system configurations suggested for Ubuntu](https://github.com/rsksmart/artifacts/tree/master/rskj-ubuntu-installer).
The most straightforward way to achieve this is to [install the RSKj node using the Ubuntu package](https://dev.rootstock.io/node-operators/setup/installation/ubuntu/).

## Module Configuration Recommendations

Depending on your use case (full node, archive node, or light node), some modules can be enabled or disabled to optimize resource usage. The minimal modules set-up to run a node which will be consumed by a Flyover Liquidity Provider Server can be found in [this file](../docker-compose/rskj/rsk.conf). Please notice that the provided configuration is just a regtest example so it shouldn't be used in production, however, the enabled modules part can be used as a reference.

**Reference:**
- [Node Configuration and Modules](https://dev.rootstock.io/node-operators/setup/configuration/)


## Additional Resources

- [Rootstock Node Runners documentation](https://dev.rootstock.io/node-operators/setup/node-runner/)
- [Rootstock Network Monitoring Tools](https://dev.rootstock.io/dev-tools/)
- [Troubleshooting Guide](https://dev.rootstock.io/node-operators/troubleshooting/)

_This document is meant as a high-level best practices summary. Always refer to the official Rootstock DevPortal for up-to-date configuration instructions._
