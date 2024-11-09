# tilt-linkerd-demo

A demo of using Tilt to run an application with Linkerd on your local machine

## Quick Setup

1. Setup local Kubernetes cluster

   Recommend using [Orbstack](https://orbstack.dev/) (macOS only) for local Kubernetes

   Other options:

   - [Docker Desktop](https://www.docker.com/products/docker-desktop/) (macOS / Windows)
   - [Rancher Desktop](https://rancherdesktop.io/) (macOS / Windows)

   Once Kubernetes is running, verify it's working and you can connect to it with `kubectl`

   ```sh
   kubectl get nodes
   ```

2. Install Tilt

   ```sh
   brew install tilt-dev/tap/tilt
   ```

   _Full instructions are available [here](https://docs.tilt.dev/install.html)_

3. Setup `/etc/hosts` (optional, choose one of the following)

   (option A) Add the following to `/etc/hosts`:

   ```sh
   127.0.0.1   linkerd.localhost
   127.0.0.1   foo.localhost
   ```

   (option B) use `hostctl`

   ```sh
   brew install guumaster/tap/hostctl
   sudo hostctl add tilt-linkerd-demo < .etchosts
   ```

4. Install Linkerd cli (optional)

   ```sh
   curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install-edge | sh
   export PATH="$PATH:/Users/username/.linkerd/bin"
   ```

   _Full instructions are available [here](https://linkerd.io/2.16/getting-started/#step-1-install-the-cli)_

## Run the demo

```sh
tilt up
```
