# mth-core

This repository hosts common code for our services, which help with logging, authentication, API middlewares etc.
We version every minor change for `glide update` to follow.

This project produces no runnable binary. However, you can check `test` target in Makefile if you're wondering how to run all unit tests.

## Working with other repositories

Most of the time, suggested changes are not immediately merged and require testing where it's used. It's usually
a good idea to push changes to a remote branch and reference commits from consuming service:

1) Pick the latest commit hash like `fd8332512979ef10fd23c6074e5992fbe3d0341b`.
2) Navigate to project directory e.g. `mth-gateway-api`.
3) Find `mth-core` in `glide.lock` and replace revision number with `fd8332512979ef10fd23c6074e5992fbe3d0341b`.
4) Do `glide install`.

When you MR is merged, please push a new version tag to `mth-core`. Afterwards, please update consumer service's `glide.yaml` with latest version tag and do `glide update`.
