# strategy-gen
This is a tool to generate strategy for attacker service.

The tool can generate a `strategy.json` with a config file, you can get a default config with `./strategy-gen generate export`, it will export default config to `default-config.yaml`.

After you change the action config, then you can generate a strategy with `./strategy-gen generate`, it will generate a `strategy.json` file.


# get all point and action
```shell
./strategy-gen display
```

# generate strategy
```shell
./strategy-gen generate
```

or with a special config.
```shell
./strategy-gen generate -config config.yaml
```

# export default config
```shell
./strategy-gen generate export
```