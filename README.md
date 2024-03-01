# go-tinyqr

This is an experimental QR code generator that aims to be as small as possible.
It is based on [go-qrcode](https://github.com/skip2/go-qrcode) by Tom Harwood.

In order to reduce the size, I mostly removed options:

-   go-tinyqr always uses byte encoding. That's what you want most of the time anyway.
-   go-tinyqr always uses the medium error correction level.
-   go-tinyqr always uses mask pattern 0.

Fixing all of these options leads to slightly worse results, but massivly
simplifies the code. For example, it allows to hardcode the complete format
info.

Still, generating QR codes stays complex. I had hoped that I can reduce this
down to a small library that can just be copied to a new project. But it is
still ~1000 loc.

I did not remove support for different versions (sizes). Concentrating on a
single version would have similar benefits to the other optimizations in terms
of size, but I feel that the loss of flexibility would be much more relevant in
this case.

## Reduction in code size

| file       |  old | new | diff |
| ---------- | ---: | --: | ---: |
| gf.go      |   92 |  83 | -10% |
| ecc.go     |  177 | 102 | -42% |
| bitset.go  |  174 |  22 | -87% |
| version.go | 2927 | 416 | -86% |
| render.go  |  439 | 147 | -67% |
| qrcode.go  |  649 | 119 | -82% |
| **total**  | 4528 | 902 | -80% |

This table compares b6ab6a4 (old) and f88a46e (new). The numbers were measured
using `sloc`. In cases where I combined files I added up the line counts of the
old files.

The biggest single change is that I was able to remove ~2000 lines of version
information for error correction levels that I didn't use.

## Links

- [ISO/IEC 18004:2006](http://www.iso.org/iso/catalogue_detail.htm?csnumber=43655)
- [Nayuki: Creating a QR Code step by step](https://www.nayuki.io/page/creating-a-qr-code-step-by-step)
- [Thonky's QR Code tutorial](https://www.thonky.com/qr-code-tutorial/)
