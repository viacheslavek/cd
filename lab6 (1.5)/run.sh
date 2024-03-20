#!/bin/sh
flex main.cpp
g++ lex.yy.c
./a.out