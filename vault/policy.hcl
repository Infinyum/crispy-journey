path "pki_int*"                        { capabilities = ["read", "list"] }
path "pki_int/roles/crispy-dot-com"    { capabilities = ["create", "update"] }
path "pki_int/sign/crispy-dot-com"     { capabilities = ["create", "update"] }
path "pki_int/issue/crispy-dot-com"    { capabilities = ["create"] }