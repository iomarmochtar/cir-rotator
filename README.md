# cir-rotator

<p align="center">
    <p align="center"><strong>Container Image Registry Rotator</strong></p>
</p>

Container image registry can be a collection of **trash** in the case of active app CI/CD deployment pipeline and for the old images are not cleaned up, the storage bills will grow gradually and you charged for something that you never use anymore

So this tools can help you create a rotation mechanism for it, by using the powerful include and exclude filters thanks to [expr-lang](https://github.com/expr-lang/expr), also it becomes the key difference compared to other existing tools. This is the sample of the deletion execution by following criterias:

- more than 6 months since it's uploaded with size more than 100 MiB.
- ignore any repo by name containing with `base-image` for latest tag.
- ignore any repo by name containing `internal-tool` for any tag.

```
./cir-rotator delete -ho asia-southeast2-docker.pkg.dev/gcp-proj/parent-repo \
                     --if "Now() - UploadedAt >= Duration('6M')  and ImageSize >= SizeStr('100 MiB')" \
                     --ef "Repository matches '.*base-image$' and 'latest' in Tags" \
                     --ef "Repository matches '.*internal-tools.*'"
```

## Features

- Supporting various image registry, since even it's complies to [registry spec](https://docs.docker.com/registry/spec/api/) in fact some of them provide more attribute(s) in providing various information (eg: size, child repo, etc). At the moment it's supported GCR (Google Container Registry).
- Various auth methods. service account file or the basic auth one (username & password).
- Include filters, set the criteria of the image that will be involved. It can be complex by using any combination such as regex, duration comparison, etc. See the `Filter` section above for more details.
- Exclude filters, same as `Include filter` but it's used as reversed so you can ignore some image to be excluded.
- Output json, dump the result as a json file for any further inspection.
- Output table, show the result in human readable format in cli's stdout.
- Skip some images, this can be useful if you want to ignore the image that is still being in K8S cluster by dumping it then passing the list file to the argument.


## Install

### Container

The image has been uploaded in docker hub, this is the sample how to run it.
```
docker run --rm -it iomarmochtar/cir-rotator  list -u _token -p $(gcloud auth print-access-token) -ho asia.gcr.io/parent-repo --output-table
```

### Binary File

The static binary file under [release page](https://github.com/iomarmochtar/cir-rotator/releases)

## Provider/Type

These are the supported Container Registry for now (PR are welcome). by default it will determining the type based on it's provided `host` argument.

### Google Container Image

Both GCP's **Container Image** and **Artifact Registry** are supported.

#### Required Roles

To make it minimal in roles that will be granted to the executor (eg: service account), the app will recursively fetch repository list through `child` from the top parent.

So, minimum role is only `storage.admin` see the details in it's [documentation page](https://cloud.google.com/container-registry/docs/access-control#permissions_and_roles)

#### Authentication

If no authentication method provided then it will fallback to [Application Default Credential](https://cloud.google.com/docs/authentication/provide-credentials-adc), which means it also supports [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) etc.

## Filters

Filter can be set more than one to make it more specific, it divided into 2 kinds: `include` (`--include-filter` or `--if`) and `exclude` (`--exclude-filter` or `--ef`) followed by the filter string pattern. Please note that:

-  `include` filters will be executed first.
- If it's provided more than one filter then it will be grouped with `OR` operation.
- See `expr`'s [language definition](https://expr-lang.org/docs/language-definition) for available syntax.

There are also some custom function available:
- `SizeStr(string): float64`, Convert the IEC size unit so it can be operated to `ImageSize` field eg: `SizeStr('10 MiB')`.
- `Date(string): time.Time`, convert the given date string by format `yyyy-mm-dd` to `Time` object eg: `Date("2022-06-13")`.
- `Duration(string): time.Duration`, convert string to golang's duration. see [this page](https://pkg.go.dev/time#ParseDuration) for the supported pattern. but i added some custom one: `d` for day, `M` for month (30 days) and `Y` for year (365 days) eg: `Duration('1Y3M20m')`.

## How To Use

### List Repositories

listing repositories, can be used to examine the target of the repository that will be deleted. It must specified one of the output stdout stable (`--output-table`) and/or dump the result to json file (`--output-json`)

<details>
    <summary>available arguments</summary>

```
NAME:
   cir-rotator list - 

USAGE:
   cir-rotator list [command options] [arguments...]

OPTIONS:
   --output-table                      show output as table to stdout (default: false)
   --output-json value                 dump result as json file
   --allow-insecure                    allow insecure ssl verify (default: false) [$ALLOW_INSECURE_SSL]
   --basic-auth-user value, -u value   basic authentication user [$BASIC_AUTH_USER]
   --basic-auth-pwd value, -p value    basic authentication password [$BASIC_AUTH_PWD]
   --host value, --ho value            registry host [$REGISTRY_HOST]
   --type value, -t value              registry type [$REGISTRY_TYPE]
   --service-account value, -f value   service account file path, it cannot be combined if basic auth args are provided [$SA_FILE]
   --exclude-filter value, --ef value  excluding result                    (accepts multiple inputs)
   --include-filter value, --if value  only process the results of filter  (accepts multiple inputs)
   --worker-count value                http client worker count (default: 1)
   --help, -h                          show help (default: false)
```
</details>

### Delete Repositories

Deleting the repository. if not specified the filters it will **deleting all repositories inside registry**, so better for you to examine first using `list` command or use option `--dry-run`.

<details>
    <summary>available arguments</summary>

```
NAME:
   cir-rotator delete - 

USAGE:
   cir-rotator delete [command options] [arguments...]

OPTIONS:
   --output-table                      show output as table to stdout (default: false)
   --output-json value                 dump result as json file
   --allow-insecure                    allow insecure ssl verify (default: false) [$ALLOW_INSECURE_SSL]
   --basic-auth-user value, -u value   basic authentication user [$BASIC_AUTH_USER]
   --basic-auth-pwd value, -p value    basic authentication password [$BASIC_AUTH_PWD]
   --host value, --ho value            registry host [$REGISTRY_HOST]
   --type value, -t value              registry type [$REGISTRY_TYPE]
   --service-account value, -f value   service account file path, it cannot be combined if basic auth args are provided [$SA_FILE]
   --exclude-filter value, --ef value  excluding result                    (accepts multiple inputs)
   --include-filter value, --if value  only process the results of filter  (accepts multiple inputs)
   --worker-count value                http client worker count (default: 1)
   --dry-run                           just log the action, will not deleting (default: false)
   --skip-list value                   path of file that contains skipping list, will be ignored if matched
   --repo-list value                   path of file containing repositories that will be deleted, this can be generated from list action
   --skip-error                        if any error happen while deleting just ignore it (default: false)
   --help, -h                          show help (default: false)
```
</details>
