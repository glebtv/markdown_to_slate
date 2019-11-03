## Markdown to SlateJS

[![Build Status](https://travis-ci.org/glebtv/markdown_to_slate.svg?branch=master)](https://travis-ci.org/glebtv/markdown_to_slate)

Markdown to SlateJS editor state converter in golang.

Also supports converting slatejs state to plain text

Current state: prototype (most stuff works, has some problems)

https://slatejs.org/

Example:

go run example/main.go > view/src/data.json
cd view
yarn install
yarn start
