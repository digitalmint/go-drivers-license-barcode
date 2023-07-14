# Go Drivers License Barcode

Golang package to help with extracting barcode data from a Driver's License.
Currently Supported:

* Date of Birth
* Date of Expiration
* Serial

This package does not attempt to look at version headers or even properly follow the specification since not all barcodes properly follow it.
You can try [this GO package](https://github.com/tka-tech/DLID) which follows the spec. However, for our needs, we had issues with this package not working and the code in this repo worked better for extracting the above fields.