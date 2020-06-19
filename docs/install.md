# installation the "gitops way"

## required prerequisites
- [ ] define your fqdn
   ```bash
   export ZITADELFQDN="zitadel.dev"
   ```
- [ ] define api and rest endpoints for your organization
  - e.g. :
  ```bash
     export ZITADELACCOUNTS=accounts.${ZITADELFQDN}
     export ZITADELAPI=api.${ZITADELFQDN}
     export ZITADELISSUER=issuer.${ZITADELFQDN}
     export ZITADELCONSOLE=console.${ZITADELFQDN}
  ```
- [ ] create an ops repository to desribe your cluster
- [ ] create a folderstructure for your environments
  - e.g. :
  ```bash
  k8s/workload
  ├── cockroach-db
  │   └── secure
  │       ├── base
  │       └── overlay
  │           ├── cockroachsecure
  │           └── generic
  ├── storage-class 
  └── zitadel
      └── overlay
          └── dev
              ├── console
              ├── hosts
              ├── mappings
              └── migrations
  ```

## recommended prerequisites for convenience
- [ ] choose and install reconciler for kubernetes ( e.g. [argocd](https://argoproj.github.io/argo-cd/))
- [ ] create a secrets repository with [gopass](https://github.com/gopasspw/gopass)

## install with kubernetes
- [ ] create environments overlay
- [ ] configure namespace
- [ ] configure hosts
  - [ ] prepare configure SSL certificates when not using let´s encrypt
- [ ] configure mapping
- [ ] define secrets/structure
- [ ] install product

## install locally with container
