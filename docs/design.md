
# MAAS

Metal As A Service (MAAS) allows physical and virtual machines to be managed
using an API.

MAAS supports a number of management controllers, for example, IPMI, HMC, 
AMT, etc. For a more complete list see [BMC Power Types][power].

The lifecycle of a _MAAS machine_ consists of the following actions:

- [*Enlist*][enlist] registers MAC addresses, etc. and changes the machine 
status to `New`.
- [*Commission*][commission] builds an inventory of the RAM, CPU, SSDs, NICs
GPUs, etc. and changes the machine status to `Ready`.
- Acquire reserves a node for use and changes the machine status to `Allocated`.
- [*Deploy*][deploy] installs or boots an operating system image and changes the
node status to `Deployed`.
- *Release* changes the machine status back to `Ready`. Local storage may be 
erased.

Each MAAS machine is assigned a SystemID during *enlistment*.

For API documentation see the [gomaasapi][gomaasapi] library or [ReST][rest] API.

# Requirements

- It MUST be possible to allocate a physical or virtual machine.
- It MUST be possible to for allocated machines to optionally initialize or
join a Kubernetes cluster.
- If a implementation decides not to implement Kubernetes provisioning, the
generic implementation MUST provision Kubernetes.

## Resource management

- Given a Cluster API `Machine`, it MUST be possible to determine which MAAS
machine it represents.
- Given a MAAS machine, it MUST be possible to determine which Cluster API 
`Machine` it is represents.

# Design

## Create 

1. Acquire MAAS machine.
1. If error, return and allow create to be retried.
1. Update Cluster API `Machine` object with `ProviderID == SystemID`.
1. If error, release MAAS machine and allow create to be retried.
1. Deploy MAAS machine.
1. If error, log and wait for an administrator to correct problem.

Machines may be leaked if:

- "Update Cluster API `Machine`" failed and release also fails.
- "Deploy MAAS machine" fails.

## Upgrade

TODO

## Delete

1. Release MAAS machine (power off, optionally erase disk, machine is `Ready`).
1. If error, allow delete to be retried.
1. Delete Cluster API `Machine` resource.
1. If error, log and wait for an administrator to correct problem.

Machines may be leaked if:

- "Delete Cluster API `Machine`" fails.

[enlist]: https://docs.maas.io/2.0/en/installconfig-add-nodes#enlistment)
[commission]: https://docs.maas.io/2.0/en/installconfig-commission-nodes
[deploy]: https://docs.maas.io/2.0/en/installconfig-nodes-deploy
[gomaasapi]: https://github.com/juju/gomaasapi
[rest]: https://docs.maas.io/2.5/en/api
[power]: https://docs.maas.io/2.0/en/installconfig-power-types
