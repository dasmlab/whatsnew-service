#!/bin/bash

msg=$1
echo "Adding all things"
git add -A . 

echo "commiting with msg: $msg"
git commit -m "$msg"
echo "pushing"
git push origin main

