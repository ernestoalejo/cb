{{ if not .Exists }}'use strict';

{{ else }}
{{ end }}
describe('Service: {{ .Data.Name }}', function() {
  beforeEach(module('services.{{ .Data.Module }}'));
});
