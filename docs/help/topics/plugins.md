# Aspect CLI Plugins

Plugins allow you to customize Bazel's behavior, and they're easy to write!
A plugin can subscribe to the Build Event Protocol (BEP), to react in real-time during the build.
Plugins can contribute custom commands like `lint` so developers can live in a single tool.

## High-level design

A plugin is any program with a gRPC server that implements our plugin protocol.

We provide convenient support for writing plugins in Go, but this is not required.
You can write a plugin in any language.
Plugins are hosted and versioned independently from the aspect CLI.

The aspect CLI process spawns the plugin as a subprocess, then connects as a
gRPC client to it. The client and server run a negotiation protocol to determine
version compatibility and what capabilities the plugin provides.

The plugin system is based on the excellent system developed by HashiCorp for the `terraform` CLI.
You can read more about this archecture here:
<https://github.com/hashicorp/go-plugin/blob/master/README.md>

## Quickstart

Use the https://github.com/aspect-build/aspect-cli-plugin-template repo to create a starter repo.

Follow instructions on the README to customize the plugin for your org.

## Plugin configuration

In a `.aspect/cli/plugins.yaml` file at the repository root, list the plugins you'd like to install.

This is a YAML file. A typical example is as follows:

```yaml
- name: hello-world
  from: github.com/aspect-build/aspect-cli-plugin-template
  version: v0.2.0
```

The `from` line points to the plugin binary and can take one of these forms:

1. A string with no slash characters, which is interpreted as a program on your system `PATH`.
2. A filesystem path, either relative to the workspace root or absolute.
3. A string of the form `github.com/some-org/some-repo`.

    In this case, a `version` property is required as well.
    This form follows the convention in https://github.com/aspect-build/aspect-cli-plugin-template
    where a GitHub release at a tag contains the plugin binaries as assets.

    To get a binary for the right platform, we append one of these platform suffixes before fetching:
    `-darwin_amd64`, `-darwin_arm64`, `-linux_amd64`, `-linux_arm64`, `-windows_amd64.exe`

    In the yaml example above, on an x86_64 architecture Linux machine, we would download from
    `https://github.com/aspect-build/aspect-cli-plugin-template/releases/download/v0.2.0/hello-world-linux_amd64`

4. An http/https URL from which the plugin can be downloaded.

    As in the previous case, a platform suffix is appended to the URL before fetching.

## Roadmap

In the future, we plan to allow semantic versioning ranges to constrain the versions which can be used.
When aspect runs, would then prompt you to re-lock the dependencies to exact versions if they
have changed, and can verify the integrity of the plugin contents against what was first installed.

> The locking semantics follow the [Trust on first use] approach.

Another future enhancement is for From to accept a string starting with `//`, which is interpreted as a [Bazel Label] in the current workspace.

    When the `from` line is a label, it will be a `*_binary` rule which builds a plugin binary.
    When the CLI loads this plugin, it first builds it from source.
    This is useful as a local development round-trip while authoring a plugin. However, it is not a
    great way to deploy a plugin to users, as it causes them to perform an extra build every time
    they run `aspect`, whether they intend to use the plugin or not.

[trust on first use]: https://en.wikipedia.org/wiki/Trust_on_first_use
[bazel label]: https://bazel.build/concepts/labels
