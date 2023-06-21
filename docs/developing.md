## Developing

### Dependency

1. 1 coordinator node and 1 executor node to collect from via ssh, need to have public key auth setup
2. kubectl access to a kubernetes cluster that can be used for testing
3. copy default-test.json to root directory and name it test.json

### Scripts

On Linux, Mac, and WSL there are some shell scripts modeled off the [GitHub ones](https://github.com/github/scripts-to-rule-them-all)

to get started run

```sh
./script/bootstrap
```

after a pull it is a good idea to run

```sh
./script/update
```

tests

```sh
./script/test
```

before checkin run

```sh
./script/cibuild
```

to cut a release do the following

```sh
#dont forget to update changelog.md with the release notes
git tag v0.1.1
git push origin v0.1.1
./script/release v0.1.1
```


