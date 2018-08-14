load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    prefix = "github.com/stripe/go-confidant-client",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/stripe/go-confidant-client",
    visibility = ["//visibility:private"],
    deps = [
        "//confidant:go_default_library",
        "//kmsauth:go_default_library",
    ],
)

go_binary(
    name = "go-confidant-client",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
