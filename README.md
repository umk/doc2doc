# doc2doc

**doc2doc** is a simple command-line tool that takes an input file, then creates or updates an output file. It compares the new input with the previous one to keep unchanged parts of the output and formatting intact, without needing extra instructions for the AI.

Usage: `doc2doc [arguments...]`

### Arguments

| Argument           | Description                                                                 |
|---------------------|-----------------------------------------------------------------------------|
| `-d string`         | Metadata file path                                                          |
| `-force`            | Force generation                                                            |
| `-gen.p float`      | Generation top P                                                           |
| `-gen.seed int`     | Generation seed                                                            |
| `-gen.t float`      | Generation temperature                                                     |
| `-i value`          | Input file path (required)                                                 |
| `-meta`             | Only save metadata given input and output                                  |
| `-o string`         | Output file path (required)                                                |
| `-svc.base string`  | Service base URL                                                           |
| `-svc.key string`   | Service key                                                                 |
| `-svc.model string` | Service model name                                                          |
| `-y`                | Confirm automatically                                                       |

The `-i` value can also be set to `-`, which tells the program to read the value from `stdin`.

### Environment variables

| Environment Variable | Description                 |
|----------------------|----------------------------|
| `D2D_BASE_URL`       | Base URL of the OpenAI-compatible LLM service          |
| `D2D_KEY`            | A key for LLM service                               |
| `D2D_MODEL`          | A model to use when generating the output |

If base URL is omitted, OpenAI service is used. The key must be specified either way to make the OpenAI client library happy. If model is not specified, `gpt-4o` is used.

### Example

Refresh the application arguments list in README.md based on the output of `doc2doc --help`

```sh
./doc2doc --help 2>&1 | ./d2d-ollama.sh "qwen2.5:14b" -d .doc2doc/arguments.d2d -i - -o README.md
```
