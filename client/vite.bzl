"vite_project rule"

load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo", "ExternalNpmPackageInfo", "run_node")
load("@build_bazel_rules_nodejs//internal/linker:link_node_modules.bzl", "module_mappings_aspect")
load("@bazel_skylib//lib:paths.bzl", "paths")

_DEFAULT_VITE = "@npm//vite/bin:vite"

def _vite_project_impl(ctx):
    arguments = ctx.actions.args()
    execution_requirements = {}
    progress_prefix = "Compiling Vite project"

    arguments.add_all(["build", paths.dirname(ctx.file.index_html.short_path)])

    deps_depsets = []
    inputs = ctx.files.srcs[:] + ctx.files.index_html[:]
    for dep in ctx.attr.deps:
        if ExternalNpmPackageInfo in dep:
            deps_depsets.append(dep[ExternalNpmPackageInfo].sources)
        if DeclarationInfo in dep:
            deps_depsets.append(dep[DeclarationInfo].transitive_declarations)
        if DefaultInfo in dep:
            deps_depsets.append(dep[DefaultInfo].files)
    inputs.extend(depset(transitive = deps_depsets).to_list())

    outputs = []
    dist_dir = ctx.actions.declare_directory("dist")
    outputs.extend([dist_dir])
    arguments.add_all(["--outDir", "../" + dist_dir.path])

    run_node(
        ctx,
        inputs = inputs,
        arguments = [arguments],
        outputs = outputs,
        mnemonic = "ViteProject",
        executable = "vite",
        execution_requirements = execution_requirements,
        progress_message = "%s" % (
            progress_prefix,
        ),
    )

    providers = [
        DefaultInfo(
            files = depset(outputs),
        ),
    ]

    return providers

vite_project = rule(
    implementation = _vite_project_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
        ),
        "index_html": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "deps": attr.label_list(
            providers = [
                [DeclarationInfo],
            ],
            aspects = [module_mappings_aspect],
        ),
        "vite": attr.label(
            default = Label(_DEFAULT_VITE),
            executable = True,
            cfg = "host",
        ),
    },
)

# Avoid using non-normalized paths (workspace/../other_workspace/path)
def _to_manifest_path(ctx, file):
    if file.short_path.startswith("../"):
        return file.short_path[3:]
    else:
        return ctx.workspace_name + "/" + file.short_path

def _vite_prodev_impl(ctx):
    out = ctx.actions.declare_file(ctx.attr.name + ".sh")
    ctx.actions.expand_template(
        template = ctx.file.launcher_template,
        output = out,
        substitutions = {
            "TEMPLATED_main": _to_manifest_path(ctx, ctx.executable.prodevserver),
        },
        is_executable = True,
    )

    deps_depsets = []
    inputs = ctx.files.srcs[:] + ctx.files.index_html[:]
    for dep in ctx.attr.deps:
        if ExternalNpmPackageInfo in dep:
            deps_depsets.append(dep[ExternalNpmPackageInfo].sources)
        if DeclarationInfo in dep:
            deps_depsets.append(dep[DeclarationInfo].transitive_declarations)
        if DefaultInfo in dep:
            deps_depsets.append(dep[DefaultInfo].files)
    inputs.extend(depset(transitive = deps_depsets).to_list())

    files = inputs

    transitive = [
        ctx.attr.prodevserver[DefaultInfo].default_runfiles.files,
    ]

    runfiles = ctx.runfiles(
        files = files,
        transitive_files = depset([], transitive = transitive),
    )

    return [
        DefaultInfo(
            executable = out,
            runfiles = runfiles,
        ),
    ]

vite_prodev = rule(
    implementation = _vite_prodev_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
        ),
        "index_html": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "deps": attr.label_list(
            providers = [
                [DeclarationInfo],
            ],
            aspects = [module_mappings_aspect],
        ),
        "vite": attr.label(
            default = Label(_DEFAULT_VITE),
            executable = True,
            cfg = "host",
        ),
        "prodevserver": attr.label(
            default = "//client/vite:prodevserver",
            executable = True,
            cfg = "host",
        ),
        "launcher_template": attr.label(
            allow_single_file = True,
            default = "//client/vite:launcher_template.sh",
        ),
    },
    executable = True,
)
