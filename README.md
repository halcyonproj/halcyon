# Halcyon

> The open control plane for the software-defined datacenter.

Halcyon is an open-source control plane for modeling, enforcing, and operating physical datacenter infrastructure as software.
It allows infrastructure teams to declaratively describe racks, power, cooling, networking, servers, and firmware, continuously enforcing the constraints that keep real-world hardware safe, predictable, and supportable.

Halcyon is built for teams running on-prem or private-cloud infrastructure who want cloud-level discipline (declarative state, enforcement, safe rollouts) without giving up control of their own hardware.

## Getting Started

> [!WARNING]
> Halcyon is under active development. APIs and models are evolving.
>
> Please see the [docs](./docs) for architecture, concepts, and current usage.

A step-by-step getting started guide is coming soon.

## What Halcyon Is

Halcyon treats the physical datacenter as a governed system, not a collection of loosely related vendor devices.

At its core, Halcyon provides:
- A datacenter model: sites, rooms, rows, racks, servers, PDUs, feeds, networks, and their relationships
- A policy and constraint engine that enforces physical limits (power, cooling, redundancy, placement)
- Drift detection between declared intent and observed reality
- Firmware and BIOS rollout primitives with safety checks and controlled concurrency
- Telemetry-aware enforcement, tying real measurements back to declared assumptions
- Extensible integrations (Redfish, SNMP, vendor APIs, network devices, etc)
- A clean API, CLI, and UI for operators and automation

Halcyon is designed to answer questions like:
- Is this rack operating within its validated power envelope?
- Can I safely place another GPU node here?
- What assumptions does this deployment rely on, and are they still true?
- What changed before things broke?

Halcyon is a good fit if you:
- Operate on-prem or private-cloud infrastructure
- Run GPU-dense, power-constrained, or failure-sensitive environments
- Care about knowing when systems are operating safely
- Want infrastructure decisions encoded and enforceable, not implicit

## What Halcyon Is Not

Halcyon is intentionally opinionated about what it does not try to be.
- Not a hardware vendor. Halcyon integrates with hardware vendors; it does not replace them.
- Not traditional DCIM. It is not a static inventory system or a visualization-only tool. Halcyon _enforces_ behavior, not just documents it.
- Not a workload scheduler. Halcyon governs infrastructure safety and constraints. It does not schedule applications, jobs, or containers.
- Not a monitoring or SIEM replacement. Halcyon consumes telemetry to enforce constraints, but it is not a general-purpose observability platform.
- Not a magic abstraction layer. Physical reality still matters. Halcyon does not hide constraints, it makes them intentional and explicit.
- Not a cloud. Halcyon is designed to help teams operate hardware and software safely, not pretend it isn’t there.

Halcyon may not be a good fit if you:
- Want a simple inventory tracker
- Don’t control your hardware
- Expect physical constraints to be someone else’s problem

## Support
For questions, discussions, or issues, please [open a GitHub issue](https://github.com/halcyonproj/halcyon/issues).
