# BlueServer

BlueServer is a web app that turns lists of Chinese words and sentences into beautiful Anki cards. It is written in plain Go, HAML and coffee script and uses cc-cedict as its Chinese-English dictionary and the online dictionary of the Taiwanese Ministry of Education as its Chinese-Chinese dictionary. You can see BlueServer live [here](http://thebluemandarinlab.com:8080/)

## Installation

BlueServer requires a sqlite database of the CC-CEDICT as generated by [cedictTxt2Db](https://github.com/yangchuanzhang/cedictTxt2Db). If you do not need the newest version of the dictionary, you can download a recent version [here](https://www.dropbox.com/s/277fmaofyaf0dvn/cedict.sqlite3). When you execute blueServer, $CEDICT_DB has to contain the path of this file.

BlueServer also requires haml, coffee script and go to be installed on your system.

Once you have installed all the requirements, you can install go by changing into the $GOPATH/src/github.com/yangchuanzhang/blueServer directory and executing build.sh. 

## Running BlueServer

To run blueServer, change into the build/ subdirectory of your installation and type blueServer. Building and executing the server can be combined by executing
        ./build.sh && (cd build && CEDICT_DB="/path/to/cedict.sqlite3" blueServer )
in a terminal.

## Implementation Details

BlueServer is written in Go using "net/http". blueServer.go contains the main() function that sets up the webserver. Once it is running, it responds to the following routes:
        /vocab/
        /vocab/lookup/<word>
        /sentences/
        /sentences/lookup/<sentence>
        /assets/<file>
