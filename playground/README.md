## Golox Playground

Repo for the Golox Playground.


## Commands

[Bun](https://bun.sh/) is used as the package manager and bundler.


### Build
```sh
bun run build
```

The above command, creates build files in `playground/dist` ready to be served.

### Dev Server

```sh
bun run dev
```

The above command will watch for changes and rebuild the playground. You can open `playground/dist/index.html` directly in your browser and reload on any change(Bun doesn't support hot reloading at the moment).

