# Create key & certificate using Kubernetes CA

Three variables used in the following steps:

```bash
# SERVICE is the name of the Vault service in Kubernetes.
# It does not have to match the actual running service, though it may help for consistency.
SERVICE=vault-server

# NAMESPACE where the Vault service is running.
NAMESPACE=default

# SECRET_NAME to create in the Kubernetes secrets store.
SECRET_NAME=vault-server-tls

# TMPDIR is a temporary working directory.
TMPDIR=~/tmp
```

1. Create a key for Kubernetes to sign

    ```bash
    openssl genrsa -out ${TMPDIR}/vault.key 2048
    ```

2. Create a Certificate Signing Request (CSR)
    1. Create a file `${TMPDIR}/vault-csr.conf` with the following contents:

        ```bash
        cat <<EOF >${TMPDIR}/vault-csr.conf
        [req]
        req_extensions = v3_req
        distinguished_name = req_distinguished_name
        [req_distinguished_name]
        [ v3_req ]
        basicConstraints = CA:FALSE
        keyUsage = nonRepudiation, digitalSignature, keyEncipherment
        extendedKeyUsage = serverAuth
        subjectAltName = @alt_names
        [alt_names]
        DNS.1 = ${SERVICE}
        DNS.2 = ${SERVICE}.${NAMESPACE}
        DNS.3 = ${SERVICE}.${NAMESPACE}.svc
        DNS.4 = ${SERVICE}.${NAMESPACE}.svc.cluster.local
        IP.1 = 127.0.0.1
        EOF
        ```

    2. Create the CSR

        ```bash
        openssl req -new -key ${TMPDIR}/vault.key -subj "/CN=${SERVICE}.${NAMESPACE}.svc" -out ${TMPDIR}/vault.csr -config ${TMPDIR}/vault-csr.conf
        ```

3. Sign the CSR with K8s' CA
    - From `kubectl`

        It did not work with `apiVersion: certificates.k8s.io/v1` because of incompatibility between `usages` and `signerName` which does not happen when using API version `v1beta1`

        1. Create a file `${TMPDIR}/csr.yaml` with the following contents:

            ```bash
            export CSR_NAME=vault-csr
            cat <<EOF >${TMPDIR}/csr.yaml
            apiVersion: certificates.k8s.io/v1
            kind: CertificateSigningRequest
            metadata:
              name: ${CSR_NAME}
            spec:
              groups:
              - system:authenticated
              request: $(cat ${TMPDIR}/server.csr | base64 | tr -d '\n')
              signerName: kubernetes.io/kube-apiserver-client
              usages:
              - digital signature
              - key encipherment
              - client auth
            EOF
            ```

        2. Send the CSR to Kubernetes

            ```bash
            kubectl apply -f ${TMPDIR}/csr.yaml
            ```

        3. Once created, approve the CSR in Kubernetes

            ```bash
            kubectl certificate approve ${CSR_NAME}
            ```

        4. Retrieve the certificate

            ```bash
            serverCert=$(kubectl get csr ${CSR_NAME} -o jsonpath='{.status.certificate}')
            ```

        5. Write the certificate out to a file

            ```bash
            echo "${serverCert}" | openssl base64 -d -A -out ${TMPDIR}/vault.crt
            ```

        6. Retrieve Kubernetes CA

            ```bash
            kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}' | base64 -d > ${TMPDIR}/vault.ca
            ```

    - Manually
        1. If running Minikube, retrieve the CA files

            ```bash
            cp ~/.minikube/ca.{crt,key} $TMPDIR
            ```

        2. Create a file `${TMPDIR}/vault-ca.conf` with the following contents:

            ```bash
            cat <<EOF >${TMPDIR}/vault-ca.conf
            extensions = v3_ca
            [ v3_ca ]
            basicConstraints = CA:FALSE
            keyUsage = nonRepudiation, digitalSignature, keyEncipherment
            extendedKeyUsage = serverAuth
            subjectAltName = @alt_names
            [alt_names]
            DNS.1 = ${SERVICE}
            DNS.2 = ${SERVICE}.${NAMESPACE}
            DNS.3 = ${SERVICE}.${NAMESPACE}.svc
            DNS.4 = ${SERVICE}.${NAMESPACE}.svc.cluster.local
            IP.1 = 127.0.0.1
            EOF
            ```

        3. Sign the CSR

            ```bash
            openssl x509 -req -in $TMPDIR/vault.csr -CA $TMPDIR/ca.crt -CAkey $TMPDIR/ca.key -CAcreateserial -out $TMPDIR/vault.crt -extfile $TMPDIR/vault-ca.conf -days 365
            ```

4. Store the key, cert, and Kubernetes CA into Kubernetes secrets

    ```bash
    kubectl create secret generic ${SECRET_NAME} \
            --namespace ${NAMESPACE} \
            --from-file=vault.key=${TMPDIR}/vault.key \
            --from-file=vault.crt=${TMPDIR}/vault.crt \
            --from-file=vault.ca=${TMPDIR}/ca.crt
    ```

# Same steps for Vault Injector

This time, variables are:

```bash
SERVICE=vault-injector

NAMESPACE=default

SECRET_NAME=vault-injector-tls

TMPDIR=~/tmp
```

# Helm configuration

1. Create a `helm-values.yaml` file with the following contents:

    ```bash
    global:
      tlsDisable: false

    injector:
      # True if you want to enable vault agent injection.
      enabled: true
      certs:
        # secretName is the name of the secret that has the TLS certificate and
        # private key to serve the injector webhook. If this is null, then the
        # injector will default to its automatic management mode that will assign
        # a service account to the injector to generate its own certificates.
        secretName: vault-injector-tls

        # caBundle is a base64-encoded PEM-encoded certificate bundle for the
        # CA that signed the TLS certificate that the webhook serves. This must
        # be set if secretName is non-null.
        caBundle: "LS0tLS1CRUdJTiBDRVJUSU...(see command below)"

        # certName and keyName are the names of the files within the secret for
        # the TLS cert and private key, respectively. These have reasonable
        # defaults but can be customized if necessary.
        certName: vault-injector.crt
        keyName: vault-injector.key

    server:
      extraEnvironmentVars:
        VAULT_CACERT: /vault/userconfig/vault-server-tls/vault.ca

      extraVolumes:
        - type: secret
          name: vault-server-tls # Matches the ${SECRET_NAME} from above

      standalone:
        enabled: true
        config: |
          listener "tcp" {
            address = "[::]:8200"
            cluster_address = "[::]:8201"
            tls_cert_file = "/vault/userconfig/vault-server-tls/vault.crt"
            tls_key_file  = "/vault/userconfig/vault-server-tls/vault.key"
            tls_client_ca_file = "/vault/userconfig/vault-server-tls/vault.ca"
          }

          storage "file" {
            path = "/vault/data"
          }
    ```

    To get the CA bundle value:

    ```bash
    cat ${TMPDIR}/ca.crt | base64
    ```

2. Install the Helm chart using the values file:

    ```bash
    helm repo add hashicorp https://helm.releases.hashicorp.com
    helm repo update
    helm install vault hashicorp/vault --values=helm-values.yaml
    ```

# Vault

Now that connection to the Vault server is done over TLS, we are setting up the public key infrastructure (PKI).

We will set up both a self-signed root CA **and** an intermediary CA within Vault because of the following (from the [documentation](https://learn.hashicorp.com/tutorials/vault/pki-engine)):

> If users do not import the CA chains, the browser will complain about self-signed certificates.

Useful [link](https://learn.hashicorp.com/tutorials/vault/kubernetes-cert-manager) more specific to doing this when Vault is deployed through a Helm chart

## Initializing, unsealing and logging in

1. Initialize the Vault

    ```bash
    kubectl exec vault-0 -- vault operator init -key-shares=1 -key-threshold=1 -format=json > keys.json
    ```

2. Grab the unseal key

    ```bash
    VAULT_UNSEAL_KEY=$(cat keys.json | jq -r ".unseal_keys_b64[]")
    ```

3. Unseal the Vault

    ```bash
    kubectl exec vault-0 -- vault operator unseal $VAULT_UNSEAL_KEY
    ```

4. Get the root token

    ```bash
    VAULT_ROOT_TOKEN=$(cat keys.json | jq -r ".root_token")
    ```

5. Log in to the Vault using the root token

    ```bash
    kubectl exec vault-0 -- vault login $VAULT_ROOT_TOKEN
    ```

6. Start an interactive shell session with the Vault pod

    ```bash
    kubectl exec --stdin=true --tty=true vault-0 -- /bin/sh
    ```

## Generate Root CA

In this step, we are going to generate a self-signed root certificate using PKI secrets engine.

1. Enable the `pki` secrets engine

    ```bash
    vault secrets enable pki
    ```

2. Give it a maximum TTL of 10 years

    ```bash
    vault secrets tune -max-lease-ttl=87600h pki
    ```

3. Generate root certificate

    ```bash
    vault write -field=certificate pki/root/generate/internal \
            common_name="crispy.com" \
            ttl=87600h > CA_cert.crt
    ```

    You may need to do `cd /tmp` in order for this to work

4. Configure CA and certificate revokation list (CRL) URLs

    ```bash
    vault write pki/config/urls \
            issuing_certificates="http://127.0.0.1:8200/v1/pki/ca" \
            crl_distribution_points="http://127.0.0.1:8200/v1/pki/crl"
    ```

## Intermediate CA

Now, we are going to create an intermediate CA using the root CA you regenerated in the previous step.

1. Enable a second `pki` secrets engine at the `pki_int/` path

    ```bash
    vault secrets enable -path=pki_int pki
    ```

2. Give it a maximum TTL of 5 years

    ```bash
    vault secrets tune -max-lease-ttl=43800h pki_int
    ```

3. Generate the intermediate CA

    ```bash
    vault write -field=csr pki_int/intermediate/generate/internal \
            common_name="crispy.com Intermediate Authority" > pki_intermediate.csr
    ```

4. Sign the intermediate certificate request with the root CA

    ```bash
    vault write -field=certificate pki/root/sign-intermediate csr=@pki_intermediate.csr \
            format=pem_bundle ttl="43800h" > intermediate.cert.pem
    ```

5. Import this signed certificate back in Vault

    ```bash
    vault write pki_int/intermediate/set-signed certificate=@intermediate.cert.pem
    ```

## Create a role

1. Create a role allowing subdomains and bare domain names

    ```bash
    vault write pki_int/roles/crispy-dot-com \
            allowed_domains="crispy.com" \
            allow_subdomains=true \
            require_cn=false \
            allow_bare_domains=true \
            max_ttl="72h"
    ```

2. Request a certificate to test

    ```bash
    vault write pki_int/issue/crispy-dot-com common_name="test.crispy.com" ttl="24h"
    ```

3. Create a policy which will allow the token to request the certificates associated with the `crispy-dot-com` role

    ```bash
    vault policy write pki - <<EOF
    path "pki_int*"                        { capabilities = ["read", "list"] }
    path "pki_int/roles/crispy-dot-com"    { capabilities = ["create", "update"] }
    path "pki_int/sign/crispy-dot-com"     { capabilities = ["create", "update"] }
    path "pki_int/issue/crispy-dot-com"    { capabilities = ["create"] }
    EOF
    ```

## Kubernetes authentication

This will allow `cert-manager` to authenticate with Vault in order to request certificates

1. Enable the authentication method

    ```bash
    vault auth enable kubernetes
    ```

2. Configure it to use the service account token, the Kubernetes host URL and its certificate

    ```bash
    vault write auth/kubernetes/config \
        token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
        kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443" \
        kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    ```

    The `token_reviewer_jwt` and `kubernetes_ca_cert` reference files written to the container by Kubernetes upon creation. The environment variable `KUBERNETES_PORT_443_TCP_ADDR` references the internal network address of the Kubernetes host.

3. Create a Kubernetes authentication role named `issuer` bound to the `pki` policy created above and linked to a Kubernetes service account named `issuer`, which is what was mounted to the pod and configured in the `token_reviewer_jwt` field

    ```bash
    vault write auth/kubernetes/role/issuer \
        bound_service_account_names=issuer \
        bound_service_account_namespaces=default \
        policies=pki \
        ttl=20m
    ```

    The role connects the Kubernetes service account, `issuer`, in the `default` namespace with the `pki` Vault policy. The tokens returned after authentication are valid for 20 minutes. This service account is created further down.

4. Lastly, exit the `vault-0` pod

    ```bash
    exit
    ```

# Cert Manager

## Installation

1. Using Helm, install `cert-manager` in one command:

    ```bash
    helm install \
      cert-manager jetstack/cert-manager \
      --namespace cert-manager \
      --create-namespace \
      --set installCRDs=true
    ```

2. Get all the pods within the cert-manager namespace

    ```bash
    kubectl get pods --namespace cert-manager
    ```

    Wait until the pods prefixed with `cert-manager` are running and ready (`1/1`).

## Configure an issuer

> `cert-manager` enables you to define Issuers that interface with the Vault certificate generating endpoints. These Issuers are invoked when a Certificate is created.

1. Create a service account `issuer` within the default namespace, which is what will be used when authenticating with Vault

    ```bash
    kubectl create serviceaccount issuer
    ```

2. Create a variable to hold the service account's secret name

    ```bash
    ISSUER_SECRET_REF=$(kubectl get serviceaccount issuer -o json | jq -r ".secrets[].name")
    ```

3. Create another variable to hold the certificate used to validate the Vault server (in this case, the Minikube certificate), PEM encoded and in base64

    ```bash
    CA_BUNDLE=$(cat $TMPDIR/ca.crt | base64)
    ```

4. Create an Issuer that defines Vault as the CA

    ```bash
    cat <<EOF | kubectl apply -f -
    apiVersion: cert-manager.io/v1
    kind: Issuer
    metadata:
      name: vault-issuer
      namespace: default
    spec:
      vault:
        server: https://vault.default.svc.cluster.local:8200
        path: pki_int/sign/crispy-dot-com
        caBundle: $CA_BUNDLE
        auth:
          kubernetes:
            role: issuer
            secretRef:
              name: $ISSUER_SECRET_REF
              key: token
    EOF
    ```

5. Running `kubectl get issuer -o wide` should return `Vault verified` as a `STATUS`

## Secure the Ingress resource

1. Install `ingress-nginx` using Helm

    ```bash
    helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
    helm repo update
    helm install ingress-nginx ingress-nginx/ingress-nginx
    ```

2. Create an Ingress resource with the proper *annotation*

    ```yaml
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: crispy
      namespace: default
      annotations:
        nginx.ingress.kubernetes.io/rewrite-target: /
        kubernetes.io/ingress.class: nginx
        cert-manager.io/issuer: vault-issuer
    spec:
      rules:
      - host: crispy.com
        http:
          paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: crispy
                port:
                  number: 80
      tls:
      - hosts:
        - crispy.com
        secretName: crispy-com-tls
    ```

3. Now you should be able to connect to the application using HTTPS! If using Minikube, you will need to enter a command to be able to reach the Ingress resource:

    ```bash
    minikube tunnel
    ```