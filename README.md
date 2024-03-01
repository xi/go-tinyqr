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

## Links

- [ISO/IEC 18004:2006](http://www.iso.org/iso/catalogue_detail.htm?csnumber=43655)
- [Nayuki: Creating a QR Code step by step](https://www.nayuki.io/page/creating-a-qr-code-step-by-step)
- [Thonky's QR Code tutorial](https://www.thonky.com/qr-code-tutorial/)
