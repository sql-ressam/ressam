# Ressam [WIP]

Ressam (Turkish meaning artist) - is lightweight CLI database diagram explorer 
with relationships prediction (e.g. if you don't have them for performance reasons).
Ressam can show database diagrams in the self-hosted Web app, YAML or JSON.
No telemetry and external calls, fully open.

Example usage:

```sh
export RESSAM_DRIVER=pg # aliases: postgres, postgresql
export RESSAM_DSN="postgresql://user:password@127.0.0.1:5432/ressam?sslmode=disable"
ressam draw
```

TODO:

* MySQL support
* MSSQL support
* Mongo support?
* Export to JSON, YAML
* Relationship prediction
