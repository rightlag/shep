<p align="center">
  <a href="https://github.com/rightlag/shep"><img src="https://cdn.rawgit.com/rightlag/shep/master/assets/title.svg" alt="shep"></a>
</p>

<p align="center">
  <a href="https://travis-ci.org/rightlag/shep"><img src="https://img.shields.io/travis/rightlag/shep.svg?style=flat-square" alt="Build Status"></a>
  <a href="https://coveralls.io/github/rightlag/shep"><img src="https://img.shields.io/coveralls/rightlag/shep.svg?style=flat-square" alt="Coverage Status"></a>
  <a href="https://goreportcard.com/report/github.com/rightlag/shep"><img src="https://goreportcard.com/badge/github.com/rightlag/shep?style=flat-square" alt="Go Report Card"></a>
</p>

shep is a Go library that parses JSON Schema vocabulary documents and validates client-submitted data. It supports validation for the seven primitive types defined below:

- `null`
- `boolean`
- `object`
- `array`
- `number`
- `integer`
- `string`

# Testing

    $ go test ./... -v

# License

```
MIT License

Copyright (c) 2017 Jason Walsh

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
