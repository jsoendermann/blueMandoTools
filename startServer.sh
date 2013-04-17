#! /bin/bash

rm -f vocab.html
haml vocab.haml > vocab.html

rm -f assets/vocab.js 
coffee -c assets/vocab.coffee 

rm -f assets/sentences.js
coffee -c assets/sentences.coffee

go install 
blueServer 
