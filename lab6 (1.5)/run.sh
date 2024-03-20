#!/bin/sh
flex main.cpp
g++ -std=c++17 lex.yy.c -o main
./main