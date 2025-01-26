## Golox Playground

Repo for the [Golox Playground](https://golox.tushartripathi.me/).

![Playground Web](../assets/playground-web.gif)

![Playground Mobile](../assets/playground-mobile.gif)

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

The above command will watch for changes and rebuild the playground. You can open `playground/dist/index.html` directly in your browser and reload on any change(Bun doesn't support hot reloading at the moment). Also if you change the Go code, lox examples or the worker code, you would need to stop and rereun the command.

### Updating wasm_exec.js 

This is copied from inside `tinygo` with some changes done to make it work with Bun. You can run `bun run copy-wasm-exec` to copy the new version if needed, and then add it on top of the worker file, as we run WASM in a web worker.

```sh
$(tinygo env TINYGOROOT)/targets/wasm_exec.js
```

### Todo

- Once Netlify has latest bun, we can use it for build and deploy. Currently it's a manual process, using Netlfiy's CLI.
- Change the interpreter to return error throughout instead of panicking. As TinyGo [doesn't support recover](https://github.com/tinygo-org/tinygo/pull/4380) in WASM yet, the experience in case of runtime error is not good.
- Write lezer grammar for lox, to use it for syntax highlighting in code mirror. I tried this but it was quite complex and didn't seem worth the effort. Currently, I've done two things on top of JS' language pack to make the experience nicer -
    - added lox based code snippets for autocomplete which overrides the js ones
    - added a view plugin to syntax highlight lox keywords not present in the js grammar
- Use Tailwind v4 instead of the current [Tailwind Lite](https://tailwind-lite.com/), this will signicantly reduce the size of the bundle but it's not supported in build [yet](https://bun.sh/docs/bundler/fullstack#this-is-a-work-in-progress).
