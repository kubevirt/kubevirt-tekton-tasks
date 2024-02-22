# Modify Windows ISO file

This tasks is modifying windows iso file. It replaces prompt bootloader with non prompt one. This helps with automation of 
windows which requires EFI - the prompt bootloader will not continue with installation until some key is pressed. The non prompt 
bootloader will not require any key pres.

### Parameters

{% for item in task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}


### Usage

Task run using resolver:
```
{{ task_run_resolver_yaml | to_nice_yaml }}```

### Platforms

The Task can be run on linux/amd64 platform.
