# GP

Do you remember that old `Build, Ship, Run` thing? well, `gp` is a tool that abstract all those things.

If you have a file in your repo describing the component, gp can build it, create the deployment scripts and run it.

Yes, this is correct, as a developer, you can focus on writing your code and forget about how to create, configure and deploy `YAML` files.

The system allows a hands-free experience while providing all the resources in a manner that, if you want to, you can tinker with the different configurations.

Configuration hidration follows a hierarchical process allowing you to contribute and collaborate to the level you feel comfortable.


## What's the idea behind it?

Building Distributed Systems is a very complex task. These are all the tasks necessary to build

- Write a set of microservices
- Write the CI Pipeline to build microservices
- Build those microservices
- Store them in a registry
- Write the Deployment Configuration (Environment dependant) for those microservices
- Deploy them
- Monitor them
- Update them

There's nothing really new in this list, however, there's a lot of pressure for Devs to learn as much Infrastructure as possible, but there's a catch, they not only have to learn infra, they have to learn how infra is built and managed in their companies (yes, everybody does it differently).

There is a clear push-back from Devs to learn Kubernetes, Helm, Kustomize... and someone has to wonder why do they.

There are 2 issues:

- Developers getting frustrated because they need to work on things they do not enjoy: wiring infra yamls
- Operations people getting frustrated because developers `do not understand` how things work.

Both sides are correct. Why does it happen?

Developers do not have the knowledge to run Operations. Operations do not have enough knowledge to understand the applications.

Can we automate it? yes, with a caveat: we need more info, from both parties. This is where `gp` comes along. `gp` is just one piece of the solution, the other 2 are `Defaults` and a `Metadata System`.

Once you put these 3 components together you can automate a lot of the current friction between parties.

- gp automates the configuration translation.
- Defaults sets the minimum sets of constraints and configurations needed for a target system (or environment, or context).
- The Metadata System analyses the system and enhances both the application/component defintions and the live system (or environment, or context).

One more thing... The system allows strict control of the build process. This means that developers have the freedom to write and build the way they want but the organisation can enforce a set of controls at any stage.

This approach guarantees hands-off standarisation and enforced regulation across all the applications without impossing restrictions to any team (this is done by giving ownership of the resources to the right teams).

An application team might be responsible for defining a set of microservices but they're not responsible for the specific configuration of those microservices into a high available multi-cloud environment.

Network concerns, storage concerns, security concerns are given to the right team. This is achieved by redefining the `Interface` of the systems. Teams do not interact with Kubernetes, or Terraform, or Cloudformation, teams interact with gp exposing application definitions and constraints.

### What does it mean in real life?

Application teams will focus on writing applications, microservices, components.
Ops teams will focus on defining the rules to deploy applications on their clusters.
DevOps teams will focus on writing automation components that facilitate running systems at scale.

What the model does is to separate responsibilities. It's not about building a wall between the 3 teams, it's about establishing proper communication channels with a lingua franca for all.

At its core, this model is meta automation: building a system that enables automation by automating a set of concerns.



### Quick Start

First you need to configure gp:

```
$gp configure --user ivan --registry harbor.alacasa.uk --verbose true 

```

This command will create a config file and a set of helper templates under `$HOME/.config/gp`. The previous command will create the following config file:

```
registry: harbor.alacasa.uk/ivan/
registry_user: ivan
targets:
- name: local
  type: docker
  platform: darwin/arm64
  allow_latest: false
- name: k8s
  type: helm
  platform: linux/amd64
  allow_latest: false
docker:
  push: false
  overwrite: false
  build_info: true
```

Helper Templates are fetched from a `Defaults` repo. That repo can be specified with the `configure` command.

Let's look at a component descriptor file:

```
name: ghreleases
desc: Simple API to list the last releases of a set of Github repos.
lang: go
port: 8080
src: 
  git: https://gitea.alacasa.uk/ivan/ghreleases.git
cmd: "ghreleases -f ./conf/catalog"

config:
  files:
  - name: catalog
    path: ./conf/catalog.json
    include: true


connects:
  - name: github
    type: http
    url: https://api.github.com
    constraint:
      hardcoded: true
```

This file provides all the information to build and run this component in any defined target, in this case, `docker` and `kubernetes` using helm charts.





## Simple case, one component

- Dockerfile: defines the composition of the component's runtime.
- docker-compose: defines how to run the component.
- helm chart: defines how to deploy the component in a k8s cluster.
- helmrelease: defines the instance of the component in a cluster.

## One Component with Dependencies and Configuration


## Two components


```
name: ghreleases-airtable
# meta:
lang: go
port: 8080
src: 
  git: ssh://ivan@synology.alacasa.uk:/volume1/git/ghreleases-airtable
cmd: "ghreleases-airtable -k $AirtableAPIToken -i $AirtableDBID -t $AirtableTableID"

config:
  files:
  - name: catalog
    path: /tmp/catalog.json
    include: true
    constraints:
      hardcoded: true
      storage:
        type: readonly

secrets:
  third-party:
    - AirtableAPIToken
    - AirtableDBID
    - AirtableTableID

connects:
  - name: github
    type: http
    url: "https://api.github.com"
  - name: airtable
    type: http
    url: "https://airtable.com"
    
```