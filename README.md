# Devoxx OPA demo

This is the demo application for my _decoupling policy from application code
using Open Policy Agent_ talk at Devoxx 2023.

It consists of a web service (PEP - Policy Enforcement Point) that serves a
templated HTML file with a hardcoded list of "files" and their attributes. There
is also a "Policy Administration Point" (PAP) under `tools/pap` that does live
rebuilds of the policy files in `tools/pap/policy`, serving them to the OPA
PDP - Policy Decision Point - integration in the PEP.

There is no login system. The user attributes are stored under
`data/userinfo.yaml`.

OPA is integrated using the OPA SDK.

## How to run

Prerequisites:

- Golang

Steps:

- `go run ./tools/pap`
- `go run .`
- go to http://127.0.0.1:8080/
- edit the `tools/pap/policy/policy.rego` file a bit and refresh the page to see
  the results

If you want to play more with Rego as a language, you might also be interested
in the Rego playground at https://play.openpolicyagent.org/
