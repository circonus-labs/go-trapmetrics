# v0.0.13

* fix(lint): unused args
* chore(lint): struct alignment
* build(deps): bump github.com/circonus-labs/go-trapcheck from 0.0.12 to 0.0.13

# v0.0.12

* build(deps): bump github.com/circonus-labs/go-trapcheck from 0.0.10 to 0.0.11
* build(deps): bump github.com/circonus-labs/go-apiclient from 0.7.19 to 0.7.23
* build(deps): bump golangci/golangci-lint-action from 3.3.0 to 3.4.0

# v0.0.11

* feat: add BytesSentGzip stat
* build(deps): update go-trapcheck from v0.0.9 to v0.0.10

# v0.0.10

* fix: send all 64bit numbers as bignum_as_string in json

# v0.0.9

* feat: add `UpdateCheckTags` method
* feat: add `QueueCheckTag` method
* feat(deps): bump go-apiclient from 0.7.15 to 0.7.18
* feat(deps): bump go-trapcheck from 0.0.8 to 0.0.9
* chore: update to go1.17
* fix(lint): ioutil deprecation

# v0.0.8

* upd: go-trapcheck v0.0.8
* add: additional text tests (leading/trailing spaces, non-printable char)
* upd: generic text cleaner method
* add: trim leading/trailing spaces from text metric values
* add: replace non-printable chars in text metric values
* add: replace 'smart' quotes (with regular quotes) in text metric values
* add: escape any embedded quotes in text metric values
* add: non-printable char replacement to config (default '_')
* upd: check type assertions
* add: more tests for each metric type

# v0.0.7

* upd: go-trapcheck v0.0.7
* upd: use bytes.Buffer
* add: FlushRawJSON to send pre-formatted metrics
* add: FlushWithBuffer to pass in a buffer from a pool
* upd: use bytes.Buffer for metrics

# v0.0.6

* upd: dep go-trapcheck

# v0.0.5

* upd: dep go-trapcheck
* fix: handle zero metrics to send correctly

# v0.0.4

* upd: dep go-trapcheck

# v0.0.3

* build(deps): bump github.com/circonus-labs/go-trapcheck

# v0.0.2

* build(deps): bump github.com/circonus-labs/go-trapcheck
* add: dependabot config
* fix: lint issues
* add: lint config/action

# v0.0.1

* initial
