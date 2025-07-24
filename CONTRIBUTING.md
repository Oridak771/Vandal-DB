# Contributing to Vandal

First off, thank you for considering contributing to Vandal! It's people like you that make Vandal such a great tool.

## Where do I go from here?

If you've noticed a bug or have a question, [search the issue tracker](https://github.com/vandal/vandal/issues) to see if someone else has already created a ticket. If not, feel free to [create a new issue](https://github.com/vandal/vandal/issues/new).

## Fork & create a branch

If you want to contribute with code, please fork the repository and create a new branch from `main`.

```bash
git checkout -b my-new-feature
```

## Development

To get started with the development, you will need to have Go and Docker installed. You can run the controller manager locally against a Kubernetes cluster by running:

```bash
make run
```

## Testing

To run the tests, you can use the following command:

```bash
make test
```

## Submitting a pull request

When you're ready to submit a pull request, please make sure that your code is well-tested and that you have updated the documentation if necessary.

Thank you for your contribution!
