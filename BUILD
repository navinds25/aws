load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "gazelle")

go_prefix("github.com/navinds25/mission-ctrl")

# bazel rule definition
gazelle(
  prefix = "github.com/navinds25/mission-ctrl",
  name = "gazelle",
  command = "fix",
)
