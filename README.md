# guinea-completion

A tiny Go HTTP service that forwards prompts to the OpenAI completions API
using only the standard library. POST `{"prompt": "..."}` to `/complete` and
it returns the generated text.

Requires the `OPENAI_API_KEY` environment variable; `PORT` overrides the
default port 8080. Run with `go run .`.
