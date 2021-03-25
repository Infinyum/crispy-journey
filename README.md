# crispy-journey

Simple web server in Golang used as an excuse to play with different tools

# `.github/workflows/`

- GitHub Actions as a CI platform
- Will be compared with a Jenkins instance running locally

# `ec2/`

- Terraform (.tf file) to create a simple AWS EC2 instance which accepts SSH connections and exposes port 80 where the server is running
- Ansible playbook for setting up the EC2 instance and deploying the web server as a Docker image
- HashiCorp Vault configuration file for a Vault standalone instance running locally storing the Dockerhub credentials

# `k8s/`

- yaml configuration files to deploy the web server on a Kubernetes cluster (specifically a Minikube instance running locally)
- Ingress controller to enable HTTPS for the web server, which requires a TLS certificate
- cert-manager resources to do just that

# `pki/`

- Contains some resources to enable Vault as a certificate authority within Kubernetes

# `src/`

- Golang web server in its own `crispy` package, along with a basic test file to be added as a step in CI pipeline