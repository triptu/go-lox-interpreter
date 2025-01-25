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



### Todo

- Once Netlify has latest bun, we can use it for build and deploy.
- Change the interpreter to return error throughout instead of panicking. As TinyGo doesn't support panic, the
WASM experience in case of runtime error is not good.
- Write lezer grammar for lox, to use it for syntax highlighting in code mirror. Currently, we've done two things on top of JS' language pack -
    - added lox based code snippets for autocomplete which overrides the js ones
    - added a view plugin to syntax highlight lox keywords not present in the js grammar
