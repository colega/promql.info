# promql.info

Source code for [promql.info](https://promql.info).

Development:

```shell
npm install
npm run install-dev
npm run watch-css # in one terminal
npm run watch-js # in another terminal
npm run watch-go # in another terminal
```

App will be served on https://localhost:8080, but [air](github.com/air-verse/air) will proxy it on https://localhost:8081, which is the URL you should visit as that will automatically reload.