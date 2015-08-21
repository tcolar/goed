Contents of ~/.goed :
  - config.toml : User (customized) configuration file
  - themes/ : User (customized) themes, create from scratch or copied from standard/themes/
  - actions/ : User (customized) actions

  - default/ : Original goed files, do not edit directly as will be replaced upon upgrades.
  - default/config.toml : Standard config. Do not edit.
  - default/themes/ : standard themes. Do not edit.
  - default/actions/ : standard actions. Do not edit.

  - buffers/ content raw buffers data, DO NOT EDIT !
  - instances/ contains the current goed instance(s) API sockets.
 
Standard files are used UNLESS a customized version exists,
so for example if both `actions/gofmt.toml` and `default/actions/gofmt.toml` exists, `action/gofmt.toml` will be used.
