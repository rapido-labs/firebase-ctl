# Firebase CTL

Firebase CTL is a tool intended to automate and source control the remote configurations for firebase projects.

To make it work, we require a service account file in the machine that runs the tool, and the path to said file be provided in the environment variable,`GOOGLE_APPLICATION_CREDENTIALS`.

It offers the following facilities

### Dump existing conditions and parameters to a file
This command will query the remote-config REST API to get the configurations and dump it to the directory

```shell
firebase-ctl get remote-config --output-dir output/
```
The output will be of the following structure
```text
 output
 |__conditions
    |__conditions.json
 |__parameters
    |__parameters.json
```
### Validate a directory whether the structure is valid
Here, the user has two options
- If the `GOOGLE_APPLICATION_CREDENTIALS` environment variable is not provided, the tool just performs a validation to ensure structural integrity.
- If the `GOOGLE_APPLICATION_CREDENTIALS` environment variable is provided, the tool validates the configuration for structural correctness. Further it also performs a validation by making a call to the firebase API which makes a dry-run without applying the configuration.

The command is as follows
```shell
firebase-ctl validate remote-config --input-dir local-dir
```
The users can create multiple files under the parameters directory according to the feature set. However, uniqueness needs to be maintained across all the keys present in the files in the `parameters` directory.

### Find the diff between the source, and the current remote version
This command shows the diff for both conditions and parameters in red and green colors.
```shell
firebase-ctl diff remote-config --config-dir local-dir
```

### Apply the config
This command applies the remote-config in the input-dir
```shell
firebase-ctl apply remote-config --input-dir input-dir
```

test