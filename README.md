## Markdown to SlateJS

Project status: needs update to current version of Slate (not working with 0.57).

Currently I do not plan to update this code.

If you need a slate state to plain text converter see [here](https://github.com/glebtv/slate_to_text)

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
