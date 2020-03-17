# Replacement

This is an adaption of the Kustomize
[ReplacementTransformer](https://github.com/kubernetes-sigs/kustomize/tree/master/plugin/someteam.example.com/v1/replacementtransformer)
plugin to the config function interface.
It is based on the 
[kubeval](https://github.com/kubernetes-sigs/kustomize/tree/master/functions/examples/validator-kubeval) validation function.


## Function implementation

The function is implemented as an [image](image), and built using `make image`.

The function is implemented as a go program, which a configuration of
replacement sources and targets, and applies them to the other YAML
documents that are processed. The config file  ReplacementTransformer` syntax,
but is annotated with a `config.kubernetes.io/function` container image.

### Function configuration

A number of settings can be modified for `kubeval` in the function `spec`. See
the `API` struct definition in [main.go](image/main.go) for documentation.

## Function invocation

The function is invoked by authoring a [local Resource](local-resource)
with a metadata annotation of `config.kubernetes.io/function`:

    kustomize config run local-resource/

Please note that this modifies the local-resource documents in line; this can
be avoided by piping the output of the command when using the `--dry-run` flag.

## Running the Example

Run the validator with:

    kustomize config run local-resource/ --dry-run

This will return the input documents, with a single replacement applied
against the Service's `spec.environmentLabel`.

