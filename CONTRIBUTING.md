# Contributing to goscript

This is mostly a personal experiment, but contributions welcome for:

## Welcome
- 📚 Documentation improvements
- 🐛 Bug fixes in existing features
- 📝 Example additions
- 🎓 Educational materials
- 🚀 Roadmap features

## Discouraged
- 🔧 Performance and security work (not the goal)

## How to Contribute

1. Open issue first to discuss
2. Fork and create branch
3. Keep PRs focused and small
4. Add examples if relevant
5. Update docs

## Code Style

- Follow standard Go formatting (`gofmt`)
- Comment non-obvious logic
- Keep functions small and readable
- Educational clarity > performance

## Documenting Partial Features

If you are implementing a feature partially (e.g., a function stub or incomplete logic), please:
1. Add a test case in `transpiler/transpiler_test.go` and set it to `skip: true`.
2. Add comments in the code explaining what is missing or why it is a stub.
3. Update `ROADMAP.md` to reflect the partial status if applicable.
