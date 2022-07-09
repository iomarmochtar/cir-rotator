# cir-rotator

<p align="center">
    <p align="center"><strong>Container Image Registry Rotator</strong></p>
    <p align="center">
        <a href="https://goreportcard.com/report/github.com/iomarmochtar/cir-rotator"><img src="https://goreportcard.com/badge/github.com/iomarmochtar/cir-rotator" alt="Go Report Card"></a>
        <a href="https://codecov.io/gh/iomarmochtar/cir-rotator" > 
            <img src="https://codecov.io/gh/iomarmochtar/cir-rotator/branch/main/graph/badge.svg?token=MM0M02CDL1"/> 
        </a>
</p>

Container image registry can be a collection of trash since it is mostly used for service deployment and the total size can be increased gradually so then it bills us for something we don’t use anymore. 
 
So this tools can help you create a rotation mechanism for it, by using the powerful include and exclude filters thanks to [antonmed’s expr](github.com/antonmedv/expr).

## Install

[TODO]

## How To Use

[TODO]

### Filters

[TODO]

## TODO
- [ ] Add registry type generic ([image registry spec](https://github.com/opencontainers/distribution-spec/blob/main/spec.md))
- [ ] Add registry type ACR (Azure Container Registry)
- [ ] Add retry mechanism in http client
- [ ] Github release
- [ ] Docker release