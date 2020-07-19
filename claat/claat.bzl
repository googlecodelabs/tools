md_exts = [
    ".md",
]


def _impl(ctx):
    tmp = ctx.actions.declare_directory("TMP_" + ctx.label.name)

    src_paths = []
    for i in ctx.files.srcs:
        src_paths.append(i.path)

    cmd = "{claat} export -o {out} {srcs}".format(
        claat=ctx.executable._target.path,
        srcs=" ".join(src_paths),
        out=tmp.path,
    )

    ctx.actions.run_shell(
        inputs=ctx.files.srcs,
        outputs=[tmp],
        tools=[
            ctx.executable._target,
        ],
        command=cmd,
    )

    ctx.actions.run_shell(
        inputs=[tmp],
        outputs=[ctx.outputs.tar],
        command="tar -cf {output} -C {inn} .".format(
            inn=tmp.path,
            output=ctx.outputs.tar.path
        ),
    )

    return [DefaultInfo(runfiles=ctx.runfiles([ctx.outputs.tar]))]


claat = rule(
    implementation=_impl,
    attrs={
        # "src": attr.label(mandatory=True, allow_single_file=True),
        "srcs": attr.label_list(mandatory=True, allow_files=md_exts),
        "_target": attr.label(cfg="host",
                              allow_single_file=True, executable=True, default="@com_github_googlecodelabs_tools//claat")
    },
    outputs={
        "tar": "%{name}.tar",
    },
)
