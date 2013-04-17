#! /bin/bash

rm -rf ./build
mkdir ./build
mkdir ./build/assets

# Copy files from asstes to build, compile to .js if it's a .coffee file
for f in ./assets/*; do
  if echo "$f" | grep -qE ".coffee$"; then
    echo "Compiling $f"
    coffee -o ./build/assets/ -c $f
  else
    echo "Copying $f"
    cp $f ./build/assets/
  fi
done

# Compile haml files
for f in ./haml/*; do
  echo "Compiling $f"
  b=`basename $f`
  # FIXME this does not work if the filename contains "haml"
  haml $f > ./build/${b/haml/html}
done

go install
