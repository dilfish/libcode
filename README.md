# libcode
[![Build Status](https://travis-ci.org/dilfish/libcode.svg?branch=master)](https://travis-ci.org/dilfish/libcode)
[![codecov](https://codecov.io/gh/dilfish/libcode/branch/master/graph/badge.svg)](https://codecov.io/gh/dilfish/libcode)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)](github.com/dilfish/libcode)

### a general encoder/decoder for Chinese obfuscation
### inspired by https://sym233.github.io/core-values-encoder/

#### encode:
- 2500 common used Chinese character: map to 1-2500, map to 9 based number list with prefix 11
- unicode Chinese character: map to [1, unicode Chinese code block size], map to 9 based number list with prefix 10
- other words: map to [1, unicode code point], map to 9 based number list with prefix 9
- 9 based number list with prefix map to 12 core value word.

#### decode: 
- core value word, map to 9 based number list, check prefix to transform number list with:
- prefix 11: map to [1, 2500], map to one of the 2500 common used Chinese characters
- prefix 10: map to [1, unicode Chinese code block size], map to unicode Chinese character
- prefix 9: map to unicode code point, map to unicode character
