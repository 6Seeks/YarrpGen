
# YarrpGen: A Tool for Enabling Yarrp Support for "Entire Mode" in IPv6

## What is YarrpGen?

YarrpGen is a tool implemented in __Go__. It's designed to generate random IPv6 targets across multiple prefixes at a specified granularity. Importantly, __it guarantees an even distribution of candidate addresses over all input prefixes__.

This tool is motivate by the papers titled:

- [Yarrp'ing the Internet: Randomized High-speed Active Topology Discovery](https://dl.acm.org/doi/pdf/10.1145/2987443.2987479)
- [Yarrpbox: Detecting Middleboxes at Internet-scale](https://dl.acm.org/doi/pdf/10.1145/3595290)

## Is YarrpGen Necessary?

Imagine you have three prefixes: 2001::/16, 2004:1234:5678::/48, and 2006::1234:5678::/48. How do you choose any /64 with the same probability within them? Simply selecting a prefix at random and then choosing a /64 within it will result in uneven probabilities.

The probability of selecting a /64 in 2001::/16 is $2^{32}$ times greater than that of the other two.

## The Principle

Please excuse the brevity; I am currently busy with my PhD. I will elaborate on the principles of YarrpGen in due time.

## How Do You Use It with `yarrp`?

Remember, `yarrp` does not support "entire mode" in IPv6 by default, but it can read input addresses from stdin. `YarrpGen` provides these evenly-distributed random targets for scanning.

1. Basic Usage


`cat IPv6Prefixes.txt | ./YarrpGen`

YarrpGen requires three parameters:

- `-l` the granularity (unit prefix/network) you wish to scan, e.g., /64 or /60.
- `-c` the number of target addresses you want to generate, e.g., 1e7.
- `-i` the type of Interface Identifiers you wish to use, options: lowbyte1/fixiid/random.

2. Using with Yarrp

`cat IPv6Prefixes.txt | ./YarrpGen -l 56 -c 10000000 -i random | yarrp -I interface_name -a SrcIPv6 -M SrcMAC -G DstMAC -r 10000 -t ICMP6 -o output.yrp -i -`

Sample 10 million /56 prefixes uniformly across all assigned IPv6 Internet prefixes, then perform a traceroute with Yarrp.
