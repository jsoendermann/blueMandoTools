#! /bin/bash

rm -f vocab.html
haml vocab.haml > vocab.html
rm -f assets/vocab.js 
coffee -c assets/vocab.coffee 
go install 
blueServer 
