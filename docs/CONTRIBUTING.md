# Contributors Guide

**starknet.go** is an open-source Golang library for Cairo developed by NethermindEth. We welcome all contributions to this repository to help enrich the Starknet community.

## How to Contribute

We operate and maintain this project using an issue and pull request model. Please track the GitHub issues section [Issues](https://github.com/NethermindEth/starknet.go/issues) of this repository, and contributors can submit [Pull Requests](https://github.com/NethermindEth/starknet.go/pulls) for review by the maintainers.

### General Work-Flow

We recommend the following work-flow for contributors:

1. **Find an issue** to work on and use comments to communicate your intentions and ask questions.
2. **Work in a feature branch** of your personal fork (github.com/YOUR_NAME/starknet.go) of the main repository (github.com/NethermindEth/starknet.go).
3. After you have implemented or resolved the issue, **create a pull request** to merge your changes into the main repository.
4. Wait for the repository maintainers to **review your changes** to ensure the issue is addressed.
5. If the issue is resolved, the repository maintainers will **merge your pull request**.

### Linter Checks

**starknet.go** now requires linter checks to pass. Please follow these steps to install and run the linter:

1. **Install `golangci-lint`:**

   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.0
   ```

2. **Run the linter:**

   ```bash
   golangci-lint run
   ```

   Ensure that there are no linter errors before submitting your pull request.

3. **Run the linter on a specific file in a specific directory:**

   To run the linter on a specific file, use:

   ```bash
   golangci-lint run path/to/your/file.go
   ```

   To run the linter on all files in a specific directory, use:

   ```bash
   golangci-lint run path/to/your/directory
   ```

   Replace `path/to/your/file.go` and `path/to/your/directory` with the actual file and directory paths.

### Helpful Article on Contribution

For a detailed guide on contributing to a GitHub project, check out this [guide](https://akrabat.com/the-beginners-guide-to-contributing-to-a-github-project/).
