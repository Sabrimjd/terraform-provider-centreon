# Terraform Provider for Centreon

This Terraform Provider allows you to interact with Centreon through its API. It provides the ability to manage and query Centreon resources through Terraform.

> ⚠️ **Warning**: This provider is in early stages of development and is not ready for production use. Features may be incomplete, and breaking changes can occur without notice. Use it for testing and evaluation purposes only.

You can see the official documentation on the Terraform Provider dedicated webpage here : [https://registry.terraform.io/providers/Sabrimjd/centreon/latest/docs](https://registry.terraform.io/providers/Sabrimjd/centreon/latest/docs)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
```sh
git clone git@github.com:your-username/terraform-provider-centreon.git
```

2. Enter the repository directory
```sh
cd terraform-provider-centreon
```

3. Build the provider
```sh
make install
```

## Documentation

The provider documentation can be found in the [docs](docs/) directory.

## Development Note

Please note that the development of this provider is primarily driven by my work requirements. New features and updates are added based on my professional needs. While contributions are welcome, the scope of development is mainly focused on functionalities needed in my work environment.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

```sh
make build
```

To generate or update documentation, use the following commands:

```sh
make generate
```

In order to run the full suite of tests, run `make test`.

```sh
make test
```

## License

[MPL-2.0](LICENSE)
