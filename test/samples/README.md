# BSPL Samples

## Credit

* [`example_1.bspl`](https://github.com/mikelsr/bspl/test/samples/example_1.bspl):
from [An Evaluation of Communication Protocol Languages forEngineering Multiagent Systems]
(https://arxiv.org/pdf/1901.08441.pdf#subsection.2.5)

* `circular.bspl`: same as `example_1.bspl` but a circular dependency has been created by
adding `in price` to `Request` which is outputted by `Offer` which requires `item` generated
by `Request`
