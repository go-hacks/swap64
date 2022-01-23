# swap64
## Multi-threaded, scalable, x64 CBC 4096-bit block cipher that performs 7 rounds per block w/ an NLFSR key schedule, key dependant movements, and a 512-bit OTP key seed masked w/ passphrase hash.

Currently, this only works on Linux as it uses /dev/random
for seeding but I will probably update this to crypto-rand
later so that it works on all OSs.

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
-r/--rounds but be aware that this value is not stored
in the file anywhere so you MUST specify it at
encryption AND decryption time.

Note: There is no hash checking of the passphrase so you
can decrypt with any passphrase you like and will process
the file but you will end up with a binary blob of nothing. XD
