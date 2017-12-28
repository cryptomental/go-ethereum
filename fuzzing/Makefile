GO=go
GOFLAGS=build -buildmode=c-archive
CC=clang
CXX=clang++
CFLAGS=-O2

all: fuzzer-rlp fuzzer-bmt fuzzer-rle fuzzer-ecies fuzzer-trie

fuzzer-rlp: rlp.a
	$(CC) $(CFLAGS) fuzzer_rlp.c -c -o fuzzer_rlp.o
	$(CXX) $(CFLAGS) fuzzer_rlp.o libFuzzer.a rlp.a -o fuzzer-rlp -lpthread
rlp.a: rlp/*.go
	$(GO) $(GOFLAGS)  ./rlp/
fuzzer-bmt: bmt.a
	$(CC) $(CFLAGS) fuzzer_bmt.c -c -o fuzzer_bmt.o
	$(CXX) $(CFLAGS) fuzzer_bmt.o libFuzzer.a bmt.a -o fuzzer-bmt -lpthread
bmt.a: bmt/*.go
	$(GO) $(GOFLAGS)  ./bmt/
fuzzer-rle: rle.a
	$(CC) $(CFLAGS) fuzzer_rle.c -c -o fuzzer_rle.o
	$(CXX) $(CFLAGS) fuzzer_rle.o libFuzzer.a rle.a -o fuzzer-rle -lpthread
rle.a: rle/*.go
	$(GO) $(GOFLAGS)  ./rle/
fuzzer-ecies: ecies.a
	$(CC) $(CFLAGS) fuzzer_ecies.c -c -o fuzzer_ecies.o
	$(CXX) $(CFLAGS) fuzzer_ecies.o libFuzzer.a ecies.a -o fuzzer-ecies -lpthread
ecies.a: ecies/*.go
	$(GO) $(GOFLAGS)  ./ecies/
fuzzer-trie: trie.a
	$(CC) $(CFLAGS) fuzzer_trie.c -c -o fuzzer_trie.o
	$(CXX) $(CFLAGS) fuzzer_trie.o libFuzzer.a trie.a -o fuzzer-trie -lpthread
trie.a: trie/*.go
	$(GO) $(GOFLAGS)  ./trie/

clean:
	rm -rf rlp.a bmt.a rle.a ecies.a trie.a
	rm -rf rlp.h bmt.h rle.h ecies.h trie.h
	rm -rf fuzzer_rlp.o fuzzer_bmt.o fuzzer_rle.o fuzzer_ecies.o fuzzer_trie.o
	rm -rf fuzzer-rlp fuzzer-bmt fuzzer-rle fuzzer-ecies fuzzer-trie
