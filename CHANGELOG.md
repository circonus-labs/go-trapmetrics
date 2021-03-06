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
