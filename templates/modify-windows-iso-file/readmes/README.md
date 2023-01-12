# Modify Windows ISO file

This tasks is modifying windows iso file. It replaces prompt bootloader with non prompt one. This helps with automation of win 11 installation.

### Parameters

{% for item in task_yaml.spec.params %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}


### Usage

Please see [examples](examples) on how to run iso modification task.
The task run has to specify spec.podTemplate.securityContext! See [examples](examples) for example how to specify it.
