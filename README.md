# mangacrawl
This is a quite funcional and performatic pet project that I made two years ago. It's a crawler written in Go that focus in performance. The original goal was to download all the images(page scanlations) from a manga reader website. Our campus network back then was too unstable and often would be offline during weekends. Thus, we(the students that maintained it), tried to mirror as many web services as possible. This was one of the first, we intended to download as fast as possible, to make good use of the bandwidth at late light.

The goal was acomplished, managing to download data at saturation speeds(that for our link back then was 100mbps). Since the data volume was beyond expected, and for organization, the program tarred and gzipped the chapters as download finished. This was my first program in go and I was quite impressed by the power of goroutines and overall lightweightness of the language.

Build and run instructions are ususal go program ones. Also, since the program does not target the website, but a CDN, it shouldn't cause any major damage to the website.
