# doc2doc

**doc2doc** is a simple command-line tool that takes an input file and an optional prompt, then creates or updates an output file. It compares the new input with the previous one to keep unchanged parts of the output and formatting intact, without needing extra instructions for the AI.

Usage: `doc2doc [arguments...] <prompt>`

### Arguments

| Argument           | Description                                                                 |
|---------------------|-----------------------------------------------------------------------------|
| `-d string`         | Metadata file path                                                          |
| `-force`            | Force generation                                                            |
| `-gen.p float`      | Generation top P                                                           |
| `-gen.seed int`     | Generation seed                                                            |
| `-gen.t float`      | Generation temperature                                                     |
| `-i value`          | Input file path (required)                                                 |
| `-o string`         | Output file path (required)                                                |
| `-svc.base string`  | Service base URL                                                           |
| `-svc.key string`   | Service key                                                                 |
| `-svc.model string` | Service model name                                                          |

You can include an optional `prompt` argument to guide how the inputs should be transformed into the output file. Once provided, the prompt is saved in a metadata file, so you don't need to enter it again when updating the output, unless you want to change it. The `prompt` or `-i` values (but not both) can also be set to `-`, which tells the program to read the value from `stdin`.

### Environment variables

| Environment Variable | Description                 |
|----------------------|----------------------------|
| `D2D_BASE_URL`       | Base URL of the OpenAI-compatible LLM service          |
| `D2D_KEY`            | A key for LLM service                               |
| `D2D_MODEL`          | A model to use when generating the output |

If base URL is omitted, OpenAI service is used. The key must be specified either way to make the OpenAI client library happy. If model is not specified, `gpt-4o` is used.

### Examples

Refresh the application arguments list in README.md based on the output of `doc2doc --help`

```sh
./doc2doc --help 2>&1 | ./d2d-ollama.sh "qwen2.5:14b" -d .doc2doc/arguments.d2d -i - -o README.md "Command line arguments list in a form of a table: one column for the argument in monospace font, and another for description. The description must start from title case."
```

Or the prompt can be omitted if using existing metadata file:

```sh
./doc2doc --help 2>&1 | ./d2d-ollama.sh "qwen2.5:14b" -d .doc2doc/arguments.d2d -i - -o README.md
```

In this case the prompt will be reused from the metadata. The prompt can be edited directly in the metadata file later.