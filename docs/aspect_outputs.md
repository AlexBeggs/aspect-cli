## aspect outputs

Print paths to declared output files

### Synopsis

Queries for the outputs declared by actions generated by the given target.

Prints each output file on a line, with the mnemonic of the action that produces it,
followed by a path to the file, relative to the workspace root.

You can optionally provide an extra argument, which is a filter on the mnemonic.

```
aspect outputs <target> [mnemonic] [flags]
```

### Examples

```
# Show all outputs of the //cli/pro target, which is a go_binary:

% aspect outputs //cli/pro
 
GoCompilePkg bazel-out/k8-fastbuild/bin/cli/pro/pro.a
GoCompilePkg bazel-out/k8-fastbuild/bin/cli/pro/pro.x
GoLink bazel-out/k8-fastbuild/bin/cli/pro/pro_/pro
SourceSymlinkManifest bazel-out/k8-fastbuild/bin/cli/pro/pro_/pro.runfiles_manifest
SymlinkTree bazel-out/k8-fastbuild/bin/cli/pro/pro_/pro.runfiles/MANIFEST
Middleman bazel-out/k8-fastbuild/internal/_middlemen/cli_Spro_Spro_U_Spro-runfiles

# Show just the output of the GoLink action, which is the executable produced by a go_binary:

% aspect outputs //cli/pro GoLink
bazel-out/k8-fastbuild/bin/cli/pro/pro_/pro
```

### Options

```
  -h, --help   help for outputs
```

### Options inherited from parent commands

```
      --aspect:config string   config file (default is $HOME/.aspect/cli/config.yaml)
      --aspect:interactive     Interactive mode (e.g. prompts for user input)
```

### SEE ALSO

* [aspect](aspect.md)	 - Aspect CLI
