# Create VirtualMachine from OKD Template Task

#### Task is deprecated and will be removed in future versions.

This task creates a VirtualMachine from OKD Template.
Virtual machines can be described and parametrized in a generic form with these templates.
A bundle of predefined templates to use can be found in [Common Templates](https://github.com/kubevirt/common-templates) project.

### Parameters

{% for item in task_yaml.spec.params %}
{% if item.name == "templateParams" or 'Volume' in item.name %}
- **{{ item.name }}**: {{ item.description | replace('[', '`[')   | replace(']', ']`')}}
{% else %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endif %}
{% endfor %}

### Results

{% for item in task_yaml.spec.results %}
- **{{ item.name }}**: {{ item.description | replace('"', '`') }}
{% endfor %}

### Usage

Please see [examples](examples) on how to create VMs from a template.

### Usage in different namespaces

You can use task to do actions in different namespace. To do that, tasks requires special permissions. Apply these RBAC objects and permissions and update accordingly task run object with correct serviceAccount:

```
{% for item in rbac_yaml %}
{{ item | to_nice_yaml }}---
{% endfor %}
```
