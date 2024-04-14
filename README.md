# capi-dagger-ci

* Dagger 101
    * A way of writing CI CD in any supported language - python, go, CUE
    * Portable CI - write once run anywhere
    * No CI vendor lockin
    * Full power of programming language
    * Simple deps - Docker + Dagger CLI = profit!
* For demo
    * Build a VPC + DOKS cluster with Terraform
    * Fetch kubeconfig with doctl
    * Install and configure cluster API
    * Dagger will be the glue and orchestration
* Deploy using 
    * `source secrets.env && make deploy-infra`
    * `source secrets.env && make install-capi`
* Other commands
    * `make destroy-infra`
    * `make deploy-plan / destroy-plan`
* While deploying
    * Dagger is made up of modules and functions - `dagger functions`
    * Modules are collections of functions which represent CI tasks
    * Functions can be chained into each other and reused in a pipeline passing context along chain - `with-docreds`
    * Secrets can be injected - dagger automatically scrubs them from logs `env:xxx` markers
