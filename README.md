# swap64
## Multi-threaded, scalable, x64 CBC 4096-bit block cipher that performs 7 rounds per block w/ an NLFSR key schedule, key dependant movements, and a 512-bit OTP key seed masked w/ passphrase hash.

Currently, this only works on Linux(Maybe Mac) as it uses
/dev/random for seeding but I will probably update this
to crypto-rand later so that it works on all OSs.

Build with
``` sh
chmod +x build && ./build
```

Run with
``` sh
./swap64 fileName
```

Alternatively you can run without any arguments or
with -h/--help to show the usage information.

Encrypted files have a .sp extension and swap64 will default
to fwd and rev operating modes (-o/--opmode) based on the
extension unless otherwise specified in case you want to
encrypt something multiple times.

You can specify threads with -t/--threads otherwise
it will default to 0 which uses the number of CPUs.

The number of rounds per block can also be specified via
-r/--rounds (defaults to 7) but be aware that this value
is not stored in the file anywhere so you MUST specify
it at encryption AND decryption time. Also, due to the
design of the cipher, increasing the rounds doesn't
necessarily make it any more secure as the shuffle still
only happens one time and the rounds are a shift of that
shuffled block so since it shifts 512 bits per round,
the 8th round would circle back to the original placement.
It may add "some" additional security but just be aware
that the blocks are not actually being cycled through
the entire cipher each round.

Note: There is no hash checking of the passphrase so you
can decrypt with any passphrase you like and it will process
the file but you will end up with a binary blob of nothing. XD
