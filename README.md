# cuvva #

This is a Golang package that implements a simple web crawler.

The package is developed by [Andrey Porubov][aporubov-github]
([mailto:aporubov@gmail.com][aporubov-email]).

This is the solution to the coding challenge and it should not
be used for any purposes other than educational.

Please note that I don't have vast production experience with
Golang, it's not my primary language at the moment and I don't
use Golang on every day basis.

## Pre-requisites ##

This package relies on [x/net/html][x/net/html] package
for parsing HTML which isn't included in the standard Golang
distribution, thus it should be fetched first:

```shell
go get golang.org/x/net/html
```

## Usage ##

```go
import (
    "net/url"
    "github.com/aporubov/cuvva"
)
```

Construct the crawler providing the entry URL, then invoke the
`Crawl` method to perform the crawling and `Print` method to
pretty-print the generated sitemap, for example:

```go
URL, _ := url.Parse("https://www.cuvva.com")
crawler := cuvva.NewCrawler()
crawler.Crawl(URL)
crawler.Print()
```

The example above should produce the output similar to

```shell
...
(webpage) https://www.cuvva.com/support/cuvva-privacy-policy
    Links
        (webpage) https://www.cuvva.com/how-insurance-works
        (webpage) https://www.cuvva.com/support
        (webpage) https://www.cuvva.com/support/cuvva-privacy-policy
        (webpage) https://www.cuvva.com/about
        (webpage) https://www.cuvva.com/
        (webpage) https://www.cuvva.com/car-insurance/subscription
        (webpage) https://www.cuvva.com/get-an-estimate
        (webpage) https://www.cuvva.com/car-insurance/temporary
        (webpage) https://www.cuvva.com/car-insurance/temporary-van-insurance
        (webpage) https://www.cuvva.com/car-insurance/learner-driver
        (webpage) https://www.cuvva.com/free-car-checker
        (webpage) https://www.cuvva.com/insurance-groups
        (webpage) https://www.cuvva.com/news
        (webpage) https://www.cuvva.com/careers
        (webpage) https://www.cuvva.com/support/cuvva-cookie-policy
        (webpage) https://www.cuvva.com/support/cuvvas-terms-conditions
        (webpage) https://www.cuvva.com/support/contacting-support
        (webpage) https://www.cuvva.com/support/how-is-cuvva-regulated
        (webpage) https://www.cuvva.com/category/data-security
    Assets
        (image/vnd.microsoft.icon) https://www.cuvva.com/favicon.98576fa3.ico
        (application/javascript; charset=utf-8) https://www.cuvva.com/website.11813109.js
(webpage) https://www.cuvva.com/cuvva/how-do-our-product-teams-work-with-cops
    Links
        (webpage) https://www.cuvva.com/
        (webpage) https://www.cuvva.com/support
        (webpage) https://www.cuvva.com/how-insurance-works
        (webpage) https://www.cuvva.com/about
        (webpage) https://www.cuvva.com/get-an-estimate
        (webpage) https://www.cuvva.com/car-insurance/temporary
        (webpage) https://www.cuvva.com/car-insurance/subscription
        (webpage) https://www.cuvva.com/cuvva/1-minute-response-time-24-7-365-days-a-year
        (webpage) https://www.cuvva.com/cuvva/how-we-test-and-roll-out-new-product-features
        (webpage) https://www.cuvva.com/careers
        (webpage) https://www.cuvva.com
        (webpage) https://www.cuvva.com/car-insurance/learner-driver
        (webpage) https://www.cuvva.com/car-insurance/temporary-van-insurance
        (webpage) https://www.cuvva.com/free-car-checker
        (webpage) https://www.cuvva.com/insurance-groups
        (webpage) https://www.cuvva.com/news
        (webpage) https://www.cuvva.com/support/contacting-support
        (webpage) https://www.cuvva.com/support/cuvva-cookie-policy
        (webpage) https://www.cuvva.com/support/cuvva-privacy-policy
        (webpage) https://www.cuvva.com/support/cuvvas-terms-conditions
    Assets
        (application/javascript; charset=utf-8) https://www.cuvva.com/website.11813109.js
        (image/vnd.microsoft.icon) https://www.cuvva.com/favicon.98576fa3.ico
...
```

Please note that the crawler is limited to one domain from the entry point URL, so when crawling
`cuvva.com`, it would crawl all pages within the `cuvva.com` domain, but not follow the links
to other websites (e.g. Twitter/Facebook accounts).

## Tests ##

There is only one basic test at the moment that constructs and runs the crawler for
the `cuvva.com` domain, exactly as it's described in [Usage](#usage) section.

```shell
go test
```

TODO: Add more tests.

[aporubov-github]: https://github.com/aporubov
[aporubov-email]: mailto:aporubov@gmail.com
[x/net/html]: https://pkg.go.dev/golang.org/x/net/html