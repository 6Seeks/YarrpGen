package main

import (
  "encoding/binary"
  "flag"
  "fmt"
  "math/rand"

  "net"
  "time"
)

var prefixes = make([]uint64, 0)
var masks = make([]uint64, 0)
var acceptance = make([]float64, 0)
var alternative = make([]int, 0)

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
  column := rand.Intn(len(acceptance))
  if acceptance[column] < rand.Float64() {
    return alternative[column]
  }
  return column
}

func main() {
  var prefixLen int
  var count int
  var iid string
  flag.IntVar(&prefixLen, "l", 64, "")

  flag.IntVar(&count, "c", 1e8, "")
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
      prefixes = append(prefixes, binary.BigEndian.Uint64(ip6[:8]))
      acceptance = append(acceptance, float64(int64(1<<(64-n))))
      alternative = append(alternative, 0)
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

    offset := rand.Uint64()& masks[index]
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
