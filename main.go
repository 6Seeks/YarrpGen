package main

import (
  "encoding/binary"
  "flag"
  "fmt"
  "math"
  "math/rand"

  "net"
  "time"
)

var prefixes = make([]uint64, 0)
var masks = make([]uint64, 0)
var records = make([]uint64,0)
var acceptance = make([]float64, 0)
var alternative = make([]int, 0)
var s = rand.Float64()
var num float64 = 0

func FlushAreaDivision() {
  for i := 0; i < len(prefixes); i++ {
    acceptance[i] *= float64(len(prefixes))
  }
  // two pipes
  small := make([]int, 0)
  large := make([]int, 0)

  for i := 0; i < len(prefixes); i++ {
    if acceptance[i] < 1.0 {
      small = append(small, i)
    } else {
      large = append(large, i)
    }
  }

  var s, l int
  for len(small) > 0 && len(large) > 0 {
    s = small[0]
    small = small[1:]
    l = large[0]
    large = large[1:]
    alternative[s] = l
    acceptance[l] = acceptance[l] + acceptance[l] - 1.0
    if acceptance[l] < 1.0 {
      small = append(small, l)
    } else {
      large = append(large, l)
    }
  }
  for len(small) > 0 {
    s = small[0]
    small = small[1:]
    alternative[s] = s
  }
  for len(large) > 0 {
    l = large[0]
    large = large[1:]
    alternative[l] = l
  }
}

func Generate() int {
  s = math.Mod(s+0.6180339887498949, 1.0)
  column := int(math.Floor(s * num))
  if acceptance[column] < s*num-float64(column) {
    return alternative[column]
  }
  return column
}

func fnv1(value uint64) uint64 {
  var hash uint64 = 14695981039346656037
  for i := 0; i < 8; i++ {
    hash ^= value & 0xff
    hash *= 1099511628211
    value >>= 8
  }
  return hash
}

func main() {
  var prefixLen int
  var count int
  var iid string
  flag.IntVar(&prefixLen, "l", 64, "")

  flag.IntVar(&count, "c", 1e7, "")
  flag.StringVar(&iid, "i", "lowbyte1", "lowbyte1/fixed/random")
  flag.Parse()
  rand.Seed(time.Now().UnixNano())

  if prefixLen > 64 {
    panic("illegal prefix length")
  }

  var entire int = 0
  for {
    var line string
    if _, err := fmt.Scanln(&line); err != nil {
      break
    }
    if ip6, ip6net, err := net.ParseCIDR(line); err != nil {
      panic(line)
    } else {
      n, _ := ip6net.Mask.Size()
      if n > prefixLen {
        panic("illegal prefix length")
      }
      entire += (1 << (prefixLen - n))
      masks = append(masks, uint64((1<<(64-n))-(1<<(64-prefixLen))))
      records = append(records, rand.Uint64())
      prefixes = append(prefixes, binary.BigEndian.Uint64(ip6[:8]))
      acceptance = append(acceptance, float64(int64(1<<(64-n))))
      alternative = append(alternative, 0)
      num += 1.0
    }
  }

  FlushAreaDivision()

  ip := net.IPv6zero
  if count > entire {
    count = entire
  }

  for i := 0; i < count; i++ {
    index := Generate()
    base := prefixes[index]

    offset := fnv1(records[index]) & masks[index]
    records[index]++
    binary.BigEndian.PutUint64(ip[:8], base+offset)
    switch iid {
    case "lowbyte1":
      binary.BigEndian.PutUint64(ip[8:16], uint64(1))
    case "fixed":
      binary.BigEndian.PutUint64(ip[8:16], 0x1234567812345678)
    case "random":
      binary.BigEndian.PutUint64(ip[8:16], rand.Uint64())
    }
    fmt.Println(ip)
  }
}
