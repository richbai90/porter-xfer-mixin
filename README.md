# A Porter mixin for file transfers from build to target machine

## Configuration

```yaml
required:
  - docker:
      privileged: false
      mounts:
         # Configure a bind mount for the runtime container to use as the destination during the install step
        - Type: "bind"
          Source: "/Users/Rich/restore"
          Target: "/restore"
          ReadOnly: false
mixins:
  - xfer:
      # Use a pre-created volume as our source during build.
      # TODO: In the future use a single config option for src and infer the value to be a volume or folder path. 
      #       Similar to how the docker -v option works currently
      volume: test

install:
  - xfer:
      # Destination should match the bind mount specified above.
      # During install the runtime container has no context of the machine it is running on.
      description: File Transfer
      destination: /restore


uninstall:
  - xfer:
      # TODO: Define an actual uninstall step that removes the copied files from the destination
      description: Obligatory uninstall step
      
```