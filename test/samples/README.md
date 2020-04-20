# BSPL Samples

## Credit

* [`example_1.bspl`](https://github.com/mikelsr/bspl/test/samples/example_1.bspl):
from [An Evaluation of Communication Protocol Languages forEngineering Multiagent Systems]
(https://arxiv.org/pdf/1901.08441.pdf#subsection.2.5)

* `example_2.bspl`: same as `example_1.bspl` but the Seller can also initiate the process.
Added the action `Seller -> Buyer: Offer[out ID, out item, out price]` for that.

* `circular.bspl`: same as `example_1.bspl` but a circular dependency has been created by
adding `in price` to `Request` which is outputted by `Offer` which requires `item` generated
by `Request`
