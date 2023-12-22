# Astronomer Provider 

[GK Consulting](https://gkconsulting.dev/) is a DevOps / Infrastructure shop. We think Astronomer is the right choice for a company who doesn't want to deal with the organizational challenges of managing their own Airflow cluster. However, we didn't like the fact we couldn't provision it like all our other infrastructure - via Terraform. We built this provider so that it would benefit our clients and allow us to continue to adhere to best practices. We are committed to maintaining this repository - PRs are welcome. 

Astronomer's API isn't out of beta yet but when it comes out of beta, we will bump our provider to v1.0 (along with any necessary changes).

Please drop us a line if you'd like support, or have other IaC / DevOps needs. We would love to connect and see how we can partner together. 

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

You will need an Astronomer API token to use this provider. You'll either need to pass this into the provider via the `token` parameter or you'll need to set the `ASTRONOMER_API_TOKEN` correctly.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
