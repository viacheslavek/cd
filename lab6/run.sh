#!/bin/sh
flex main.cpp
gcc lex.yy.c
./a.out