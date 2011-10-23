#!/bin/sh

./reloader-test.app && echo RELOAD && exec ./loop.sh

case $? in
    2)
        echo ERROR
        ;;
    1)
        echo QUIT
        ;;
esac
