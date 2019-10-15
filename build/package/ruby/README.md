# Building the gem

Before building the gem it needs to have the darwin and linux `.so` files in it's `lib/native` directory. The `build_all` target in the Makefile will create those resources as well as build the gem.

```sh
make build_all
```

After the gem is built, you can install it with:

```sh
gem install ./vault-gem/vault-0.0.0.gem
```
