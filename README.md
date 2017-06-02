# aglais
[![Build Status](https://travis-ci.org/Mitchell-Riley/aglais.svg?branch=master)](https://travis-ci.org/Mitchell-Riley/aglais)

## About
aglais is an implementation of the [io language](https://github.com/stevedekorte/io) in Go. I started this because as I was fooling around with the original language implementation, I wanted to make my own changes. Here's the problem: it's written in C. I hate C. Go is much nicer. 

However, this is not a direct translation. I am also using this opportunity to learn Go, as well as compiler design.

## Installation
`go get github.com/mitchell-riley/aglais`

## Resources
* Jonathan Boyett's inspiring [europa](https://github.com/saysjonathan/europa)
* Rob Pike's [concurrent lexer](https://www.youtube.com/watch?v=HxaD_trXwRE) design
* Christoph Zenger's [course notes on grammar manipulation](http://lampwww.epfl.ch/teaching/archive/compilation-ssc/2000/part4/parsing/node3.html#SECTION00030000000000000000)
* Context free grammar checkers:
    1. http://smlweb.cpsc.ucalgary.ca/start.html
    2. http://mdaines.github.io/grammophone/
* Recursive descent parsers:
    1. http://cogitolearning.co.uk/?p=573
    2. https://en.wikipedia.org/wiki/Recursive_descent_parser#C_implementation
    3. http://www.cs.cornell.edu/courses/cs412/2008sp/lectures/lec07.pdf
