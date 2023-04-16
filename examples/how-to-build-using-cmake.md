# How to build using

With environment:
```Dockerfile
FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y wget cmake gcc g++

WORKDIR /home
RUN wget https://github.com/ttroy50/cmake-examples/archive/refs/heads/master.tar.gz
```

# Download CMake source
Download release package:
```bash
echo "TODO"
```

# Unpack
Extract source archive:
```bash
tar xf master.tar.gz
rm master.tar.gz
cd cmake-examples-master/01-basic/A-hello-cmake
```

# Configure
This is the first step for building CMake-managed project. It will generate build system files in specified directory.

`-S` specifies where sources are

`-B` specifies the build directory

```bash
mkdir build
cmake -S . -B build/
```
```plaintext
<output placeholder>
```

# Build
When build system files are ready, we can ask CMake to perform the build using these files.

`--build` specifies the build directory

`-jN` configures parallel build

```bash
cmake --build ./build -j4
```
```plaintext
<output placeholder>
```

# Run the build application

First, let's see what it is supposed to do:
```bash
cat main.cpp
```
```plaintext
<output placeholder>
```

And run it:
```bash
./build/hello_cmake
```
```plaintext
<output placeholder>
```