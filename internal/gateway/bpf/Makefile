# SPDX-License-Identifier: Apache-2.0
CXX ?= g++
CXXFLAGS ?= -O3 -Wall -Wextra -Werror -std=c++17 -mavx512f -fPIC
ASAN_FLAGS ?= -fsanitize=address -g -fno-omit-frame-pointer

TARGET = vector_jit.o
TARGET_ASAN = vector_jit_asan.o

.PHONY: all asan clean

all: $(TARGET)

$(TARGET): vector_jit.cpp vector_jit.h
	$(CXX) $(CXXFLAGS) -c vector_jit.cpp -o $(TARGET)

asan: vector_jit.cpp vector_jit.h
	$(CXX) $(CXXFLAGS) $(ASAN_FLAGS) -c vector_jit.cpp -o $(TARGET_ASAN)

clean:
	rm -f $(TARGET) $(TARGET_ASAN)
