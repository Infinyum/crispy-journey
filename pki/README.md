# Installation

1. Add the HashiCorp Helm repository.  
`helm repo add hashicorp https://helm.releases.hashicorp.com`

2. Install the latest version of the Consul Helm chart with parameters from `consul-helm-values.yaml`  
`helm install consul hashicorp/consul --values consul-helm-values.yaml`

3. Install the latest version of the Vault Helm chart with parameters from `vault-helm-values.yaml`  
`helm install vault hashicorp/vault --values vault-helm-values.yaml`

# Initialize and unseal the Vault

1. Initialize Vault with one key share and one key threshold.  
`kubectl exec vault-0 -- vault operator init -key-shares=1 -key-threshold=1 -format=json > cluster-keys.json`

2. Create a variable named VAULT_UNSEAL_KEY to capture the Vault unseal key.  
`VAULT_UNSEAL_KEY=$(cat cluster-keys.json | jq -r ".unseal_keys_b64[]")`

3. Unseal Vault running on the `vault-0` pod.
`kubectl exec vault-0 -- vault operator unseal $VAULT_UNSEAL_KEY`

4. Run the same command on pods `vault-1` and `vault-2`

5. At this point, all three pods should be in Running and Ready state.

# Configure Vault as a certificate authority

- First, start an interactive shell session on the `vault-0` pod.  
`kubectl exec --stdin=true --tty=true vault-0 -- /bin/sh`

- Log in with the root token from `cluster-keys.json`  
`vault login`

## Set up PKI engine

1. Enable the PKI secrets engine  
`vault secrets enable pki`

2. Increase the TTL by tuning the secrets engine. The default value of 30 days may be too short, so increase it to 1 year:  
`vault secrets tune -max-lease-ttl=8760h pki`

3. Configure a CA certificate and private key.  
`vault write pki/root/generate/internal \
    common_name=crispy.com \
    ttl=8760h`

4. Update the certificate revokation list (CRL) location and issuing certificates.  
`vault write pki/config/urls \
    issuing_certificates="http://127.0.0.1:8200/v1/pki/ca" \
    crl_distribution_points="http://127.0.0.1:8200/v1/pki/crl"`

5. Configure a role that maps a name in Vault to a procedure for generating a certificate. When users or machines generate credentials, they are generated against this role:  
`
vault write pki/roles/crispy-dot-com \
    allowed_domains=crispy.com \
    allow_bare_domains=true \
    allow_subdomains=true \
    max_ttl=72h
`

## Authentication

1. Enable the Kubernetes authentication method.  
`vault auth enable kubernetes`

2. Configure the Kubernetes authentication method to use the service account token, the location of the Kubernetes host, and its certificate. The referenced files and variables were automatically created with the containers.  
<code>vault write auth/kubernetes/config \
    token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
    kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443" \
    kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
</code>

3. Write out the policy named `crispy` that enables fetching certificates at `pki_int/sign/crispy-dot-com`  
`vault policy write crispy - <<EOF
path "pki/sign/crispy-dot-com" {
    capabilities = ["read", "update"]
}
EOF`

4. Create a Kubernetes authentication role, named `crispy`, that connects the Kubernetes service account name and `crispy` policy.  
`
vault write auth/kubernetes/role/crispy \
        bound_service_account_names=vault \
        bound_service_account_namespaces=default \
        policies=crispy \
        ttl=24h
`